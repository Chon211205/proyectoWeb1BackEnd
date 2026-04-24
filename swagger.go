package main

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed openapi.json
var openAPISpec string

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/openapi.json", serveOpenAPI)
	mux.HandleFunc("/docs", serveSwaggerUI)
}

func serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(openAPISpec))
}

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>

  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui-standalone-preset.js"></script>

  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: "https://proyectoweb1backend-production.up.railway.app/openapi.json",
        dom_id: "#swagger-ui",
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}