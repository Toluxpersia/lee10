package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Upvote struct {
	TotalUpvotes int64 `json:"total_upvotes"`
	IsUpvoted    bool  `json:"is_upvoted"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Children []*FileNode `json:"children,omitempty"`
}

type App struct {
	ctx        context.Context
	Workspace  []DocumentChunk
	IsIndexing bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetSuggestion(code, lineContent string, lineNumber int) string {
	prompt := fmt.Sprintf(`You are an interactive AI code editor assistant. 
Analyze the following code. The user is currently typing on line %d, which contains: %s

If you notice any syntax mistakes, bad practices, edge cases, or inefficiencies near line %d, concisely point out the issue and suggest how to improve it.
Do NOT write the code solutions yourself. Be a conversational coach. Focus only on mistakes or improvements. Keep your response extremely short (1-3 sentences max). 
If the code looks perfectly fine and there are no immediate issues or algorithm improvements, you MUST return the exact string 'OK'. Do not say anything else.

Code:
%s`, lineNumber, lineContent, lineNumber, code)


	reqBody := OllamaRequest{
		Model:  "qwen2.5-coder:1.5b", 
		Prompt: prompt,
		Stream: false,
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Ollama Error:", err)
		return ""
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		fmt.Println("Decode Error:", err)
		return ""
	}
	return ollamaResp.Response
}

func (a *App) GetChatResponse(userPrompt, code string) string {
	// RAG Integration: Retrieve context if embeddings work
	var contextString string
	queryVector, embedErr := GenerateEmbedding(userPrompt)
	if embedErr == nil && len(a.Workspace) > 0 {
		topChunks := RetrieveContext(queryVector, a.Workspace, 3)
		if len(topChunks) > 0 {
			contextString = "Here are related code snippets from other files in the user's workspace:\n\n"
			for _, chunk := range topChunks {
				contextString += fmt.Sprintf("File: %s\n```\n%s\n```\n\n", chunk.FilePath, chunk.Content)
			}
		}
	}

	prompt := fmt.Sprintf(`You are an expert AI coding assistant built into the Liten Editor.
The user asked you a direct question: "%s"

Here is their current active file code for context:
%s

%s
Please provide a concise, helpful, and direct answer.`, userPrompt, code, contextString)

	reqBody := OllamaRequest{
		Model:  "qwen2.5-coder:1.5b",
		Prompt: prompt,
		Stream: false,
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Ollama Error:", err)
		return ""
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		fmt.Println("Decode Error:", err)
		return ""
	}
	return ollamaResp.Response
}

func (a *App) OpenFolder() *FileNode {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Project Folder",
	})
	if err != nil || dir == "" {
		return nil
	}

	// Trigger Background Vector Indexer
	go func() {
		a.IsIndexing = true
		fmt.Println("Indexing workspace vectors for RAG:", dir)
		chunks, err := IndexWorkspace(dir)
		if err == nil {
			a.Workspace = chunks
			fmt.Println("RAG Indexer successful! Stored", len(chunks), "document chunk vectors.")
		} else {
			fmt.Println("RAG Indexer failed:", err)
		}
		a.IsIndexing = false
	}()

	return a.buildFileTree(dir)
}

func (a *App) ReadFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (a *App) SaveFile(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err == nil {
		go a.UpdateFileIndex(path, content)
	}
	return err
}

func (a *App) UpdateFileIndex(path, content string) {
	// Safely clear old chunks for this file and recreate
	var newWorkspace []DocumentChunk
	for _, chunk := range a.Workspace {
		if chunk.FilePath != path {
			newWorkspace = append(newWorkspace, chunk)
		}
	}

	chunks := ChunkString(content, path, 1500)
	for _, chunk := range chunks {
		vector, embedErr := GenerateEmbedding(chunk.Content)
		if embedErr == nil && vector != nil {
			chunk.Vector = vector
			newWorkspace = append(newWorkspace, chunk)
		}
	}
	a.Workspace = newWorkspace
	fmt.Println("Dynamically re-indexed saved file:", path)
}

func (a *App) CreateFile(baseDir, name string) (string, error) {
	newPath := filepath.Join(baseDir, name)
	f, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	f.Close()
	return newPath, nil
}

func (a *App) ReadFolder(dir string) *FileNode {
	return a.buildFileTree(dir)
}

func (a *App) buildFileTree(dir string) *FileNode {
	node := &FileNode{
		Name:  filepath.Base(dir),
		Path:  dir,
		IsDir: true,
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return node
	}

	for _, entry := range entries {
		if entry.Name()[0] == '.' || entry.Name() == "node_modules" {
			continue
		}
		childPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			node.Children = append(node.Children, a.buildFileTree(childPath))
		} else {
			node.Children = append(node.Children, &FileNode{
				Name:  entry.Name(),
				Path:  childPath,
				IsDir: false,
			})
		}
	}
	return node
}
