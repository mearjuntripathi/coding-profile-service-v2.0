package handler

import (
	"log"
	"net/http"
	"bytes"
	"os"
	"path/filepath"
	"github.com/yuin/goldmark"
)

// findReadme walks upward from the current directory to find README.md
func findReadme() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(dir, "README.md")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
		readmePath, err := findReadme()
		if err != nil {
			http.Error(w, "README.md not found", http.StatusNotFound)
			log.Printf("Error locating README: %v", err)
			return
		}

		data, err := os.ReadFile(readmePath)
		if err != nil {
			http.Error(w, "Failed to read README", http.StatusInternalServerError)
			log.Printf("Error reading README: %v", err)
			return
		}

		var buf bytes.Buffer
		if err := goldmark.Convert(data, &buf); err != nil {
			http.Error(w, "Failed to render README", http.StatusInternalServerError)
			log.Printf("Markdown conversion error: %v", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Project Documentation</title>
			<style>
				:root {
					color-scheme: light dark;
				}
				body {
					max-width: 900px;
					margin: 2rem auto;
					padding: 2rem;
					font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
					line-height: 1.6;
					background-color: var(--bg);
					color: var(--text);
					transition: background 0.3s, color 0.3s;
				}

				/* Light mode colors */
				@media (prefers-color-scheme: light) {
					:root {
						--bg: #ffffff;
						--text: #24292e;
						--code-bg: #f6f8fa;
						--link: #0969da;
					}
				}

				/* Dark mode colors */
				@media (prefers-color-scheme: dark) {
					:root {
						--bg: #0d1117;
						--text: #c9d1d9;
						--code-bg: #161b22;
						--link: #58a6ff;
					}
				}

				h1, h2, h3, h4 {
					border-bottom: 1px solid rgba(128,128,128,0.3);
					padding-bottom: .3em;
					margin-top: 1.4em;
				}

				a {
					color: var(--link);
					text-decoration: none;
				}
				a:hover {
					text-decoration: underline;
				}

				code {
					background: var(--code-bg);
					padding: 2px 4px;
					border-radius: 4px;
					font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
					font-size: 0.9em;
				}

				pre {
					background: var(--code-bg);
					padding: 1em;
					border-radius: 6px;
					overflow-x: auto;
					font-size: 0.9em;
				}

				img {
					max-width: 100%;
					border-radius: 4px;
				}

				blockquote {
					border-left: 4px solid var(--link);
					padding-left: 1em;
					color: #6a737d;
					margin: 1em 0;
				}

				table {
					border-collapse: collapse;
					margin: 1em 0;
					width: 100%;
				}

				th, td {
					border: 1px solid rgba(128,128,128,0.3);
					padding: 8px;
				}

				th {
					background: var(--code-bg);
				}

				hr {
					border: none;
					border-top: 1px solid rgba(128,128,128,0.3);
					margin: 2em 0;
				}
			</style>
		</head>
		<body>
		`))
		w.Write(buf.Bytes())
		w.Write([]byte("</body></html>"))
	}