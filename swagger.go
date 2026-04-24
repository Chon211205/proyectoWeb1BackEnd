package main

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.json
var openAPISpec []byte

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/openapi.json", serveOpenAPI)
	mux.HandleFunc("/docs", serveSwaggerUI)
}

func serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openAPISpec)
}

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Swagger UI</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>

  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui-standalone-preset.js"></script>

  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: "#swagger-ui",
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}