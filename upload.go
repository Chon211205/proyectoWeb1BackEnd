package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type uploadResponse struct {
	ImageURL string `json:"image_url"`
}

func registerUploadRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/upload", uploadImageHandler)
}

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiError{Error: "method not allowed"})
		return
	}

	// límite total de request: 1 MB + margen pequeño
	r.Body = http.MaxBytesReader(w, r.Body, 1_100_000)

	err := r.ParseMultipartForm(1_100_000)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "image too large or invalid form"})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "image file is required"})
		return
	}
	defer file.Close()

	if header.Size > 1_000_000 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "image must be 1 MB or smaller"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowed[ext] {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "only jpg, jpeg, png and webp are allowed"})
		return
	}

	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "could not create upload directory"})
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join("uploads", filename)

	dst, err := os.Create(path)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "could not save image"})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "could not write image"})
		return
	}

	writeJSON(w, http.StatusCreated, uploadResponse{
		ImageURL: "http://localhost:8080/uploads/" + filename,
	})
}