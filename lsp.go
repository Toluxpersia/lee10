package main

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow any origin since our Wails frontend runs on a custom hostname
	},
}

func StartLSPProxy(ctx context.Context, port string) {
	http.HandleFunc("/lsp", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("LSP Upgrade error:", err)
			return
		}
		defer ws.Close()

		cmd := exec.CommandContext(ctx, "gopls")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return
		}

		if err := cmd.Start(); err != nil {
			fmt.Println("LSP Start error:", err)
			return
		}

		// Read from LSP stdout, write to WS
		go func() {
			buf := make([]byte, 1024*32)
			for {
				n, err := stdout.Read(buf)
				if err != nil {
					break
				}
				if err := ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					break
				}
			}
		}()

		// Read from WS, write to LSP stdin
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				break
			}
			if _, err := stdin.Write(msg); err != nil {
				break
			}
		}
        
        cmd.Wait() // wait for process to finish
	})

	go func() {
		fmt.Println("Starting LSP Bridge on :" + port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Println("LSP Server Error:", err)
		}
	}()
}
