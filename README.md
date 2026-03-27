# lee10 💡

**lee10** is an incredibly robust, ultra-lightweight AI-native Code Editor Extension explicitly designed to integrate natively inside Visual Studio Code.

Rather than actively auto-completing templates blindly against your cursor, `lee10` operates exclusively as an **Architecture Coach**. It leverages a fully encapsulated React Webview Sidebar bridging perfectly into a Dual-Engine LLM topology, intelligently tracking your structural code dependencies natively via pure compiler AST traces!

---

## ✨ Advanced Features
- **Compiler Native LSP AST:** Unlike generic RAG models parsing vectors randomly, `lee10` explicitly extracts your chat keywords connecting securely into the `tsc` or `gopls` underlying Language Server. It fetches 100% physically exact code definitions across your workspace eliminating AI hallucination!
- **Dual Neural Routing:** Choose mathematically strict 100% Private Offline Inference hitting `Gemma 2:9B / Ollama`, or swap seamlessly to Premium Cloud Next-Gen Models (OpenAI GPT-5, Anthropic Sonnet 3.5, Google Gemini 3 Flash). 
- **SecretStorage Key Enclave:** Cloud API Keys are completely purged from generic `settings.json` limits! `lee10` leverages natively encrypted OS-level `vscode.secrets` vault arrays safely keeping tokens locked behind dynamic React placeholders (`••••••••`).
- **Passive Code Insight:** As you type, `lee10` analyzes code architectures silently dropping non-intrusive beautiful Code Insight native `hoverCards` exactly over your lines allowing 1-click execution straight into the Chat Webview!

---

## 🚀 Getting Started

### Prerequisites
1. **[Node.js](https://nodejs.org)** (v18+)
2. **VS Code** (v1.80+)
3. **[Ollama](https://ollama.com)** (Required ONLY if running the 100% Local Inference Engine)

### Installing from Source
1. Clone the repository natively.
2. Enter the VS Code extension sub-directory:
   ```bash
   cd lee10VS
   npm install
   ```
3. Boot the local Extension Development Host:
   ```bash
   # Press F5 dynamically inside the VS Code editor to spawn a compiled debugger clone natively!
   ```

### Publishing & Deployment
To physically build the immutable package representing the Production Visual Studio instance, compile via VSCE CLI:
```bash
# Globally install the Microsoft VSCE builder 
npm install -g @vscode/vsce

# Physically bundle the extension to generate the output .vsix artifact!
vsce package
```

---

## 🛠 Tech Stack
- **Extension Node Core**: Visual Studio Code API + Typescript
- **AST Mapping Engine**: Native `vscode.executeWorkspaceSymbolProvider` Context Extractions
- **React Frontend**: Vite + Typescript + React Markdown + Prism Code Rendering
- **Cryptography**: Native OS `SecretStorage` Mappings

## License
Distributed securely under the MIT License.
