# lee10 💡

**lee10** is a vibe coded, very bare-bones AI-native desktop code editor built to demonstrate interactive and collaborative coding. 

The main idea is to provide a seamless coding experience where the AI is always aware of your codebase. It observes as you type and does not auto-complete code or offer complete code suggestions. Instead, it offers insights and conversational advice based on your current code block. When the user needs help, they can ask the AI through the chat bar and receive a response based on their current code and the context of their codebase. This system is powered by a local vector database that stores the embeddings of the codebase. Instead of routing your codebase through the cloud, lee10 uses a local Gemma 2:9B model via Ollama for analysis. This setup ultimately delivers a seamless coding experience where the AI always knows your codebase and can provide insights and suggestions, keeping the developer in control and preventing knowledge atrophy over time through autocomplete and direct code suggestions.


## ✨ Features
- **100% Privacy & Local LLM**: Powered natively by your machine. Runs Gemma 2:9B for offline, instantaneous code generation and logic analysis.
- **On-Demand Interactive Chat Coach**: The AI doesn't just auto-complete code. When you pause typing, the background observer dynamically reads the architecture surrounding your cursor and pipes conversational advice, acting as an AI mentor exactly when you need it.
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
ollama run gemma2:9b
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
- **Intelligence**: Native HTTP Go-routines binding directly to Ollama REST endpoints (Gemma & Nomic)
- **Frontend Layer**: React + TypeScript + Monaco Editor + React Markdown

## License
Distributed under the MIT License. See `LICENSE` for more information.
