package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

//Error del api
type apiError struct {
	Error string `json:"error"`
}

// Registro de rutas
func registerRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/series", seriesCollectionHandler(db))
	mux.HandleFunc("/series/", seriesItemHandler(db))
}

// Ruta raiz
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeJSON(w, http.StatusNotFound, apiError{Error: "not found"})
		return
	}

	//Codigo HTTP 200
	writeJSON(w, http.StatusOK, map[string]string{
		"name":    "Series Tracker API",
		"version": "1.0.0",
	})
}

// Devuelve un handler que permite el acceso a la base de datos.
func seriesCollectionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		//Metodo GET (listar series)
		case http.MethodGet:
			items, err := listSeries(db)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, apiError{Error: "could not list series"})
				return
			}
			writeJSON(w, http.StatusOK, items)

		//Metodo POST (Agrega la serie)
		case http.MethodPost:
			var input SeriesInput
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				//Codigo HTTP status 400
				writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid JSON body"})
				return
			}
			if err := validateSeries(input); err != nil {
				writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
				return
			}

			item, err := createSeries(db, input)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, apiError{Error: "could not create series"})
				return
			}

			w.Header().Set("Location", "/series/"+strconv.Itoa(item.ID))
			//Codigo HTTP 201
			writeJSON(w, http.StatusCreated, item)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, apiError{Error: "method not allowed"})
		}
	}
}

// Handler para series especificas
func seriesItemHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseSeriesID(r.URL.Path)
		if !ok {
			//Codigo HTTP 404
			writeJSON(w, http.StatusNotFound, apiError{Error: "not found"})
			return
		}

		switch r.Method {

		//Metodo GET (Buscar serie)
		case http.MethodGet:
			item, err := getSeries(db, id)
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, apiError{Error: "series not found"})
				return
			}
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, apiError{Error: "could not get series"})
				return
			}
			writeJSON(w, http.StatusOK, item)

		//Metodo POST (Actualizar serie)
		case http.MethodPut:
			var input SeriesInput
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid JSON body"})
				return
			}
			if err := validateSeries(input); err != nil {
				writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
				return
			}

			item, err := updateSeries(db, id, input)
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, apiError{Error: "series not found"})
				return
			}
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, apiError{Error: "could not update series"})
				return
			}
			writeJSON(w, http.StatusOK, item)

		//Metodo DELETE (Eliminar una serie)
		case http.MethodDelete:
			err := removeSeries(db, id)
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, apiError{Error: "series not found"})
				return
			}
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, apiError{Error: "could not delete series"})
				return
			}
			//Codigo HTTP 204
			w.WriteHeader(http.StatusNoContent)
		default:
			//Codigo HTTP 405
			writeJSON(w, http.StatusMethodNotAllowed, apiError{Error: "method not allowed"})
		}
	}
}

// Verifica que el ID de la serie sea entero positivo.
func parseSeriesID(path string) (int, bool) {
	idText := strings.TrimPrefix(path, "/series/")
	if idText == "" || strings.Contains(idText, "/") {
		return 0, false
	}

	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

// Middleware cors
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Simplificar respuestas Json
func writeJSON(w http.ResponseWriter, status int, value any) {
	//Pone el header
	w.Header().Set("Content-Type", "application/json")
	//manda el codigo HTTP
	w.WriteHeader(status)
	if status != http.StatusNoContent {
		_ = json.NewEncoder(w).Encode(value)
	}
}
