import { useState, useRef, useEffect } from 'react';
import Editor, { Monaco } from "@monaco-editor/react";
import './App.css';
import { OpenFolder, ReadFile, SaveFile, CreateFile, ReadFolder, GetChatResponse, GetCoachInsight } from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";

let editorMounted = false;
const activeInsightsMap = new Map<number, string>();
let decoratorsRef: string[] = [];

function App() {
  const [code, setCode] = useState('// Welcome to lee10 AI Editor\n');
  const [fileTree, setFileTree] = useState<main.FileNode | null>(null);
  const [activeFile, setActiveFile] = useState<string | null>(null);
  const [aiSuggestion, setAiSuggestion] = useState<string | null>(null);
  const [isAiLoading, setIsAiLoading] = useState(false);
  const [isCreatingFile, setIsCreatingFile] = useState(false);
  const [newFileName, setNewFileName] = useState("");
  const [chatInput, setChatInput] = useState("");
  const [chatHistory, setChatHistory] = useState<{role: string, text: string}[]>([]);
  const [isChatting, setIsChatting] = useState(false);
  const editorRef = useRef<any>(null);
  const monacoRef = useRef<any>(null);

  function handleEditorMount(editor: any, monaco: Monaco) {
    editorRef.current = editor;
    monacoRef.current = monaco;
    editorMounted = true;
  }

  function handleEditorChange(value: string | undefined) {
    if (value !== undefined) {
      setCode(value);
    }
  }

  // Interactive AI Listener Hook
  useEffect(() => {
    if (!activeFile || !code || code.trim() === '') return;

    setIsAiLoading(true);
    const timer = setTimeout(async () => {
      try {
        if (editorRef.current) {
          const position = editorRef.current.getPosition();
          const model = editorRef.current.getModel();
          if (!model || !position) {
              setIsAiLoading(false);
              return;
          }
          const currentLine = model.getLineContent(position.lineNumber);
          
          // Get a localized block of code around the cursor (10 lines before and after)
          const startLine = Math.max(1, position.lineNumber - 10);
          const endLine = Math.min(model.getLineCount(), position.lineNumber + 10);
          const blockContent = model.getValueInRange({
            startLineNumber: startLine,
            startColumn: 1,
            endLineNumber: endLine,
            endColumn: model.getLineMaxColumn(endLine)
          });
          
          const insight = await GetCoachInsight(blockContent, currentLine, position.lineNumber);
          
          if (insight && insight.has_issues && insight.coach_message) {
            // Track the message internally
            activeInsightsMap.set(position.lineNumber, insight.coach_message);
            
            // Pipe the Coach's advice directly into the real-time AI Chat Panel instead of fighting Webkit Canvas DOM!
            setChatHistory(prev => [
                ...prev, 
                {role: 'ai', text: `💡 *Line ${position.lineNumber}:* ${insight.coach_message}`}
            ]);
            
            // Auto-open chat if hidden by turning it "on" logically via state
            setIsChatting(false);
            
            // Erase any old Monaco markers just in case
            if (monacoRef.current && editorRef.current) {
                const model = editorRef.current.getModel();
                if (model) {
                    monacoRef.current.editor.setModelMarkers(model, "ai-coach", []);
                }
            }
          }
        }
      } catch (err) {
        console.error(err);
      } finally {
        setIsAiLoading(false);
      }
    }, 1500);

    return () => clearTimeout(timer);
  }, [code, activeFile]);

  async function handleOpenFolder() {
    const tree = await OpenFolder();
    if (tree) {
      setFileTree(tree);
    }
  }

  async function submitCreateFile(e: React.FormEvent) {
    e.preventDefault();
    if (!fileTree) {
        alert("Please open a project folder first.");
        setIsCreatingFile(false);
        return;
    }
    if (!newFileName.trim()) {
        setIsCreatingFile(false);
        return;
    }

    try {
        const newPath = await CreateFile(fileTree.path, newFileName.trim());
        const newTree = await ReadFolder(fileTree.path);
        setFileTree(newTree);
        handleOpenFile(newPath);
    } catch (err) {
        console.error("Failed to create file:", err);
        alert("Could not create file: " + String(err));
    }
    setIsCreatingFile(false);
    setNewFileName("");
  }

  async function handleAskLLM(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Enter' && chatInput.trim()) {
      const userMessage = chatInput.trim();
      setChatInput("");
      setChatHistory(prev => [...prev, {role: 'user', text: userMessage}]);
      setIsChatting(true);
      
      try {
          const response = await GetChatResponse(userMessage, code);
          setChatHistory(prev => [...prev, {role: 'ai', text: response}]);
      } catch (err) {
          console.error(err);
      } finally {
          setIsChatting(false);
      }
    }
  }

  async function handleOpenFile(path: string) {
    try {
      const content = await ReadFile(path);
      setCode(content);
      setActiveFile(path);
      setAiSuggestion(null);
    } catch (err) {
      console.error("Failed to read file", err);
    }
  }

  useEffect(() => {
    const handleKeyDown = async (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 's') {
        e.preventDefault();
        if (activeFile) {
          try {
            await SaveFile(activeFile, code);
          } catch (err) {
            console.error("Failed to save", err);
          }
        }
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [code, activeFile]);

  const FileTreeNode = ({ node }: { node: main.FileNode }) => {
    const [expanded, setExpanded] = useState(false);
    
    if (node.isDir) {
      return (
        <div className="tree-node">
          <div className="explorer-item" onClick={() => setExpanded(!expanded)}>
            <span className="icon">{expanded ? '📂' : '📁'}</span>
            <span className="name">{node.name}</span>
          </div>
          {expanded && node.children && (
            <div className="tree-children" style={{ paddingLeft: '15px' }}>
              {node.children.map((child, i) => (
                <FileTreeNode key={i} node={child} />
              ))}
            </div>
          )}
        </div>
      );
    }

    return (
      <div 
        className={`explorer-item ${activeFile === node.path ? 'active' : ''}`} 
        onClick={() => handleOpenFile(node.path)}
      >
        <span className="icon">📄</span>
        <span className="name">{node.name}</span>
      </div>
    );
  };

  const getLanguage = (path: string | null) => {
      if (!path) return 'typescript';
      const ext = path.split('.').pop();
      if (ext === 'ts' || ext === 'tsx') return 'typescript';
      if (ext === 'js' || ext === 'jsx') return 'javascript';
      if (ext === 'go') return 'go';
      if (ext === 'md') return 'markdown';
      if (ext === 'css') return 'css';
      if (ext === 'json') return 'json';
      return 'typescript';
  };

  return (
    <div className="layout">
      <div className="sidebar">
        <div className="sidebar-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span className="title">lee10</span>
          <div style={{ display: 'flex', gap: '8px' }}>
            <button onClick={() => { setIsCreatingFile(true); setNewFileName(""); }} title="New File" style={{ background: 'none', border: 'none', color: '#007acc', cursor: 'pointer', fontSize: '12px', padding: '0' }}>➕</button>
            <button onClick={handleOpenFolder} title="Open Folder" style={{ background: 'none', border: 'none', color: '#fff', cursor: 'pointer', fontSize: '12px', padding: '0' }}>📂</button>
          </div>
        </div>
        <div className="explorer">
          {isCreatingFile && (
             <form onSubmit={submitCreateFile} style={{ padding: '8px 15px' }}>
                <input 
                   autoFocus
                   value={newFileName}
                   onChange={e => setNewFileName(e.target.value)}
                   onBlur={() => setIsCreatingFile(false)}
                   placeholder="filename.go"
                   style={{ width: '100%', background: '#3c3c3c', border: '1px solid #555', color: '#fff', padding: '4px', outline: 'none', boxSizing: 'border-box' }}
                />
             </form>
          )}
          {fileTree ? (
            <FileTreeNode node={fileTree} />
          ) : (
            <div style={{ padding: '15px', color: '#888', textAlign: 'center', fontSize: '13px' }}>
              No folder opened
            </div>
          )}
        </div>
      </div>

      <div className="editor-container" style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
        <div style={{ flex: 1, position: 'relative', overflow: 'hidden' }}>
          <Editor
            height="100%"
            theme="vs-dark"
            language={getLanguage(activeFile)}
          value={code}
          onChange={handleEditorChange}
          onMount={handleEditorMount}
          options={{
            minimap: { enabled: false },
            fontSize: 14,
            fontFamily: "'JetBrains Mono', 'Fira Code', 'Courier New', monospace",
            wordWrap: 'on',
            lineNumbersMinChars: 3,
            padding: { top: 16 },
            scrollBeyondLastLine: false,
            smoothScrolling: true,
            cursorBlinking: "smooth",
            cursorSmoothCaretAnimation: "on"
          }}
        />

        {aiSuggestion && (
          <div className="ai-suggestion-box">
            <div className="ai-header">
              <span>✨ lee10 Insight</span>
              <button className="ai-dismiss" onClick={() => setAiSuggestion(null)}>×</button>
            </div>
            <div className="ai-content">{aiSuggestion}</div>
          </div>
        )}

        {/* Chat Panel Overlay */}
        {chatHistory.length > 0 && (
           <div className="chat-panel">
              <div className="chat-header">
                <span>💬 lee10 Chat</span>
                <button onClick={() => setChatHistory([])}>Clear</button>
              </div>
              <div className="chat-messages">
                {chatHistory.map((msg, i) => (
                  <div key={i} className={`chat-message ${msg.role}`}>
                    <strong>{msg.role === 'user' ? 'You' : 'lee10'}: </strong>
                    <span>{msg.text}</span>
                  </div>
                ))}
                {isChatting && <div className="chat-message ai"><em>Thinking...</em></div>}
              </div>
           </div>
        )}
        </div>

        {/* Chat Command Bar */}
        <div className="chat-input-container">
           <input 
              className="chat-input"
              value={chatInput}
              onChange={e => setChatInput(e.target.value)}
              onKeyDown={handleAskLLM}
              placeholder="Ask lee10 about this code... (Press Enter)"
           />
        </div>
      </div>
    </div>
  );
}

export default App;
