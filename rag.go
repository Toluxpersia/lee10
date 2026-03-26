package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type DocumentChunk struct {
	FilePath string
	Text     string
	Vector   []float64
}

type EmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbedResponse struct {
	Embedding []float64 `json:"embedding"`
}

var (
	WorkspaceIndex []DocumentChunk
	CurrentFolder  string
	indexMutex     sync.Mutex
)

// EnsureIndex explicitly exposed to React to silently trigger semantic mapping exactly once
func (a *App) EnsureIndex() bool {
	indexMutex.Lock()
	defer indexMutex.Unlock()

	// If no folder is active, or we have already generated the index... return cleanly!
	if CurrentFolder == "" || len(WorkspaceIndex) > 0 {
		return true 
	}

	fmt.Println("🚀 Implicit Lazy-Indexing Triggered for:", CurrentFolder)
	WorkspaceIndex = make([]DocumentChunk, 0)

	err := filepath.Walk(CurrentFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "build" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".go", ".ts", ".tsx", ".js", ".jsx", ".md", ".json", ".css", ".html":
			content, err := os.ReadFile(path)
			if err == nil {
				text := string(content)
				// Truncate massively large singleton files to prevent Context OOM limits
				if len(text) > 5000 {
					text = text[:5000] 
				}

				vec := getEmbedding(text)
				if vec != nil {
					WorkspaceIndex = append(WorkspaceIndex, DocumentChunk{
						FilePath: path,
						Text:     text,
						Vector:   vec,
					})
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Indexing error:", err)
		return false
	}

	fmt.Printf("✅ Indexing Complete! Natively Embedded %d files into Unified Memory Vectors.\n", len(WorkspaceIndex))
	return true
}

func getEmbedding(text string) []float64 {
	reqBody := EmbedRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}
	body, _ := json.Marshal(reqBody)
	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var empResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&empResp); err != nil {
		return nil
	}
	return empResp.Embedding
}

func CosineSimilarity(a, b []float64) float64 {
	var dotProduct, normA, normB float64
	for i := 0; i < len(a) && i < len(b); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func RetrieveContext(query string) string {
	indexMutex.Lock()
	defer indexMutex.Unlock()

	if len(WorkspaceIndex) == 0 {
		return ""
	}

	queryVec := getEmbedding(query)
	if queryVec == nil {
		return ""
	}

	type ScoredChunk struct {
		Chunk DocumentChunk
		Score float64
	}
	
	var scores []ScoredChunk
	for _, chunk := range WorkspaceIndex {
		score := CosineSimilarity(queryVec, chunk.Vector)
		scores = append(scores, ScoredChunk{Chunk: chunk, Score: score})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Safely extract Top 2 Mathmatically Similar files
	var sb strings.Builder
	for i := 0; i < 2 && i < len(scores); i++ {
		if scores[i].Score > 0.35 {
			sb.WriteString(fmt.Sprintf("\n--- FILE PATH: %s ---\n%s\n", scores[i].Chunk.FilePath, scores[i].Chunk.Text))
		}
	}
	return sb.String()
}
