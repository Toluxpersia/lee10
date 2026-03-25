# Liten 💡

**Liten** is a lightning-fast, AI-native desktop code editor built from the ground up for privacy-first, ultra-low latency intelligence. 

Instead of routing your codebase through the cloud, Liten uses a custom **Go-based native Vector Indexer** and **Ollama** to analyze, chat, and coach you locally on your machine.

![Liten Screenshot](docs/screenshot.png) *(Add a screenshot here!)*

## ✨ Features
- **100% Privacy & Local LLMs**: Powered natively by your machine. Runs `qwen2.5-coder` for code generation and `nomic-embed-text` for vector math.
- **Dynamic RAG Workspace Indexing**: The moment you open a folder, a background Goroutine slices your code into vector algorithms. Save a file, and Liten intelligently updates the index on the fly.
- **Interactive Code Coach**: The Insight Box doesn't just auto-complete code; it reads the surrounding lines you're stuck on using the RAG index and acts as a conversational AI mentor popping up when you make a mistake.
- **Sleek Glassmorphic UI**: Powered by React, Vite, and Wails, featuring a gorgeous floating chat bar and deep Monaco Editor integration.

---

## 🚀 Getting Started

### Prerequisites
1. **[Go](https://go.dev/doc/install)** (1.18+)
2. **[Node.js](https://nodejs.org)** (18+)
3. **[Ollama](https://ollama.com)** (For Local AI)
4. **[Wails CLI](https://wails.io/docs/gettingstarted/installation)** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Setup AI Models
Liten heavily relies on these two specific local models running on Ollama:
```bash
ollama run qwen2.5-coder:1.5b
ollama pull nomic-embed-text
```

### Installation
1. Clone the repository: `git clone https://github.com/yourusername/liten.git`
2. Enter the directory: `cd liten`
3. Run the development server (auto-compiles frontend and backend):
```bash
wails dev
```
*(Note: If the `wails` command is not found, you may need to use the explicit Go path: `~/go/bin/wails dev`)*

To easily build a native macOS or Windows app binary:
```bash
wails build
# Or: ~/go/bin/wails build
```

---

## 🛠 Tech Stack
- **Backend Shell**: Go + Wails
- **Vector Math**: Custom `rag.go` Cosine Similarity engine mapping Ollama HTTP embeddings
- **Frontend Layer**: React + TypeScript + Monaco Editor

## 🤝 Contributing
Contributions are extremely welcome! Liten has massive potential, from adding multi-tab support, to implementing AST tree-sitter parsing for perfectly chunked embedding vectors!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License
Distributed under the MIT License. See `LICENSE` for more information.
