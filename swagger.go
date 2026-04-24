package main

import (
	"net/http"
	"os"
)

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/openapi.json", serveOpenAPI)
	mux.HandleFunc("/docs", serveSwaggerUI)
}

func serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	// Verifica si el archivo existe
	if _, err := os.Stat("openapi.json"); os.IsNotExist(err) {
		http.Error(w, "openapi.json not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "./openapi.json")
}

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Swagger UI</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html, body {
      margin: 0;
      padding: 0;
      height: 100%;
      background: #f5f7f8;
    }
    #swagger-ui {
      max-width: 1200px;
      margin: 0 auto;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: window.location.origin + "/openapi.json",
        dom_id: "#swagger-ui"
      });
    };
  </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}