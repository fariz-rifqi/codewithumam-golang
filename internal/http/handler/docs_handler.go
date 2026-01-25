package handler

import (
	"fmt"
	"net/http"
)

const scalarDocsHTML = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1"/>
    <title>API Docs</title>
    <style>
      body { margin: 0; }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </head>
  <body>
    <script>
      Scalar.createApiReference(document.body, {
        spec: { url: "/openapi.json" }
      })
    </script>
  </body>
</html>`

type DocsHandler struct {
	specPath string
}

func NewDocsHandler(specPath string) *DocsHandler {
	return &DocsHandler{specPath: specPath}
}

func (h *DocsHandler) ServeSpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, h.specPath)
}

func (h *DocsHandler) ServeDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, scalarDocsHTML)
}

func (h *DocsHandler) RedirectDocs(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/docs", http.StatusMovedPermanently)
}
