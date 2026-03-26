# lee10 💡

**lee10** is a vibe coded AI-native desktop code editor built to demonstrate the idea of interactive/collaborative coding. intelligence. 

Instead of routing your codebase through the cloud, lee10 uses a custom **Go-based intelligence backend** and **Ollama** to analyze, chat, and coach you locally on your machine.

![lee10 Screenshot](docs/screenshot.png) *(Add a screenshot here!)*

## ✨ Features
- **100% Privacy & Local LLM**: Powered natively by your machine. Runs `qwen2.5-coder` for offline, instantaneous code generation and logic analysis.
- **On-Demand Interactive Chat Coach**: The AI doesn't just auto-complete code. When you pause typing, the background observer dynamically reads the architecture surrounding your cursor and pipes conversational advice directly into your Chat Bar, acting as an AI mentor exactly when you need it.
- **Sleek Glassmorphic UI**: Powered by React, Vite, and Wails, featuring a gorgeous floating chat bar and deep Monaco Editor integration.

---

## 🚀 Getting Started

### Prerequisites
1. **[Go](https://go.dev/doc/install)** (1.18+)
2. **[Node.js](https://nodejs.org)** (18+)
3. **[Ollama](https://ollama.com)** (For Local AI)
4. **[Wails CLI](https://wails.io/docs/gettingstarted/installation)** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Setup AI Models
lee10 relies on these specific code-focused models running natively on Ollama:
```bash
ollama run qwen2.5-coder:1.5b
ollama pull nomic-embed-text
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
- **Intelligence**: Native HTTP Go-routines binding directly to Ollama REST endpoints (Qwen & Nomic)
- **Frontend Layer**: React + TypeScript + Monaco Editor + React Markdown

## 🤝 Contributing
Contributions are extremely welcome! lee10 has massive potential, from adding multi-tab support, to rebuilding a robust and highly-efficient local workspace RAG indexer!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License
Distributed under the MIT License. See `LICENSE` for more information.
