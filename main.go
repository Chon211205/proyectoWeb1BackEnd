package main

import (
	"database/sql"
	"fmt" 
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

//Main
func main() {
	//Abre la base de datos llamada series.db
	db, err := sql.Open("sqlite", "file:series.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Manejo de errores
	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	//Router
	mux := http.NewServeMux()
	//Registro de rutas y de swagger
	registerRoutes(mux, db)
	registerSwaggerRoutes(mux)
	registerUploadRoutes(mux)

	//Mandar las imagenes.
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	//peticiones desde frontend y evitar los errores cors
	handler := corsMiddleware(mux)

	//Levantar el servidor.
	fmt.Println("API corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}