package main

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

//go:embed openapi.json
var openAPISpec []byte

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/openapi.json", serveOpenAPI)
	mux.HandleFunc("/docs", serveSwaggerUI)
	mux.HandleFunc("/docs/", serveSwaggerUI)
}

func serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	spec := buildOpenAPISpec(r)
	if len(spec) == 0 {
		http.Error(w, "openapi spec unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(spec)
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
      const specUrl = new URL("../openapi.json", window.location.href).toString();

      SwaggerUIBundle({
        url: specUrl,
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

func buildOpenAPISpec(r *http.Request) []byte {
	rawSpec := openAPISpec
	if len(rawSpec) == 0 {
		fileSpec, err := os.ReadFile("openapi.json")
		if err == nil {
			rawSpec = fileSpec
		}
	}

	if len(rawSpec) == 0 {
		return nil
	}

	var spec map[string]any
	if err := json.Unmarshal(rawSpec, &spec); err != nil {
		return rawSpec
	}

	spec["servers"] = []map[string]string{
		{
			"url":         requestBaseURL(r),
			"description": "Servidor actual",
		},
	}

	encoded, err := json.Marshal(spec)
	if err != nil {
		return rawSpec
	}

	return encoded
}

func requestBaseURL(r *http.Request) string {
	scheme := "https"

	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = strings.Split(forwardedProto, ",")[0]
	} else if r.TLS == nil {
		scheme = "http"
	}

	if strings.HasPrefix(r.Host, "localhost") || strings.HasPrefix(r.Host, "127.0.0.1") {
		scheme = "http"
	}

	return scheme + "://" + r.Host
}
