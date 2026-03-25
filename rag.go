package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DocumentChunk struct {
	FilePath string
	Content  string
	Vector   []float32
}

type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GenerateEmbedding calls the local Ollama instance to vectorize text.
func GenerateEmbedding(text string) ([]float32, error) {
	reqBody := OllamaEmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ollama API error: %d", resp.StatusCode)
	}

	var embedResp OllamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, err
	}
	return embedResp.Embedding, nil
}

// CosineSimilarity computes the distance between two vectors. Higher = more similar.
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// ChunkString naively chunks a file by double-newlines (functions) up to a max limit.
func ChunkString(content, path string, maxChars int) []DocumentChunk {
	blocks := strings.Split(content, "\n\n")
	var chunks []DocumentChunk
	
	currentChunk := ""
	for _, block := range blocks {
		if len(currentChunk)+len(block) > maxChars {
			if strings.TrimSpace(currentChunk) != "" {
				chunks = append(chunks, DocumentChunk{FilePath: path, Content: currentChunk})
			}
			currentChunk = block
		} else {
			currentChunk += "\n\n" + block
		}
	}
	if strings.TrimSpace(currentChunk) != "" {
		chunks = append(chunks, DocumentChunk{FilePath: path, Content: currentChunk})
	}
	return chunks
}

// IndexWorkspace recursively traverses the project directory, chunks text files, and generates embeddings.
func IndexWorkspace(dir string) ([]DocumentChunk, error) {
	var allChunks []DocumentChunk
	allowedExtensions := map[string]bool{
		".go": true, ".ts": true, ".tsx": true, ".js": true, ".jsx": true, ".py": true, ".md": true, ".txt": true,
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil { return nil }
		if d.IsDir() {
			if d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		
		ext := filepath.Ext(path)
		if !allowedExtensions[ext] {
			return nil
		}

		bytes, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		content := string(bytes)

		// Create file chunks ~1500 chars limit
		chunks := ChunkString(content, path, 1500)
		
		// Vectorize chunks
		for i, chunk := range chunks {
			vector, embedErr := GenerateEmbedding(chunk.Content)
			if embedErr == nil && vector != nil {
				chunks[i].Vector = vector
				allChunks = append(allChunks, chunks[i])
			}
		}

		return nil
	})

	return allChunks, err
}

// RetrieveContext searches the vector DB for the most relevant code chunks given an embedded query.
func RetrieveContext(queryVector []float32, store []DocumentChunk, topK int) []DocumentChunk {
	if len(store) == 0 {
		return nil
	}

	type ScoredChunk struct {
		Chunk DocumentChunk
		Score float32
	}

	var scored []ScoredChunk
	for _, doc := range store {
		score := CosineSimilarity(queryVector, doc.Vector)
		scored = append(scored, ScoredChunk{Chunk: doc, Score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	var results []DocumentChunk
	for i := 0; i < topK && i < len(scored); i++ {
		// Minimum accuracy threshold to avoid hallucinating unrelated context
		if scored[i].Score > 0.4 {
			results = append(results, scored[i].Chunk)
		}
	}
	return results
}
