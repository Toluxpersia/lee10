package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"`
}

type SuggestionResponse struct {
	HasIssues    bool   `json:"has_issues"`
	CoachMessage string `json:"coach_message"`
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
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetCoachInsight(blockContent, currentLineContent string, lineNumber int) *SuggestionResponse {
	prompt := fmt.Sprintf(`You are an expert AI code coach analyzing a snippet.
The user is currently typing on line %d:
`+"`%s`"+`

Full context block:
`+"```"+`
%s
`+"```"+`

You MUST reply with ONLY a JSON object and nothing else. No markdown formatting.
If the code has no major architectural issues or bad practices, return:
{"has_issues": false, "coach_message": "OK"}

If there is a bad practice, edge case, or inefficiency, briefly suggest an improvement (1-2 sentences) without writing code:
{"has_issues": true, "coach_message": "Your coaching message here."}`, lineNumber, currentLineContent, blockContent)

	reqBody := OllamaRequest{
		Model:  "gemma2:9b",
		Prompt: prompt,
		Stream: false,
		Format: "json", // Enforce JSON mode natively in Ollama!
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Ollama Error:", err)
		return nil
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		fmt.Println("Decode Error:", err)
		return nil
	}

	rawResp := ollamaResp.Response
	// Clean up any markdown json blocks that instruct models hallucinate
	if len(rawResp) > 7 && rawResp[:7] == "```json" {
		rawResp = rawResp[7:]
	} else if len(rawResp) > 3 && rawResp[:3] == "```" {
		rawResp = rawResp[3:]
	}
	if len(rawResp) > 3 && rawResp[len(rawResp)-3:] == "```" {
		rawResp = rawResp[:len(rawResp)-3]
	}

	fmt.Println("--- Gemma2 Insight Check ---")
	fmt.Println("Raw JSON String:", rawResp)

	var insight SuggestionResponse
	if err := json.Unmarshal([]byte(rawResp), &insight); err != nil {
		fmt.Println("JSON Parse Error. Raw output was:", ollamaResp.Response, "| Err:", err)
		return nil
	}
	
	return &insight
}

func (a *App) GetChatResponse(userPrompt, code string) string {
	prompt := fmt.Sprintf(`You are an expert AI coding assistant built into the lee10 Editor.
The user asked you a direct question: "%s"

Here is their current active file code for context:
%s

Please provide a concise, helpful, and direct answer.`, userPrompt, code)

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
	return os.WriteFile(path, []byte(content), 0644)
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
