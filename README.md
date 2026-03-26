# lee10 💡

**lee10** is an AI-native desktop code editor built from the ground up for interactive intelligence. 

Instead of routing your codebase through the cloud, lee10 uses a custom **Go-based native Vector Indexer** and **Ollama** to analyze, chat, and coach you locally on your machine.

![lee10 Screenshot](docs/screenshot.png) *(Add a screenshot here!)*

## ✨ Features
- **100% Privacy & Local LLM**: Powered natively by your machine. Runs `qwen2.5-coder` for offline, instantaneous code generation and logic analysis.
- **On-Demand Interactive Code Coach**: The Insight Box doesn't just auto-complete code; it dynamically reads the surrounding architecture of any text you explicitly highlight with your cursor and acts as a conversational AI mentor popping up exactly when you need it.
- **Sleek Glassmorphic UI**: Powered by React, Vite, and Wails, featuring a gorgeous floating chat bar and deep Monaco Editor integration.

---

## 🚀 Getting Started

### Prerequisites
1. **[Go](https://go.dev/doc/install)** (1.18+)
2. **[Node.js](https://nodejs.org)** (18+)
3. **[Ollama](https://ollama.com)** (For Local AI)
4. **[Wails CLI](https://wails.io/docs/gettingstarted/installation)** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Setup AI Model
lee10 relies on this specific code-focused model running on Ollama:
```bash
ollama run qwen2.5-coder:1.5b
```

### Installation
1. Clone the repository: `git clone https://github.com/yourusername/lee10.git`
2. Enter the directory: `cd lee10`
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
Contributions are extremely welcome! lee10 has massive potential, from adding multi-tab support, to implementing AST tree-sitter parsing for perfectly chunked embedding vectors!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License
Distributed under the MIT License. See `LICENSE` for more information.
