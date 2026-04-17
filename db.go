package main

import (
	"database/sql"
	"errors"
)

// Estructura de la tabla de series
type Series struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	CurrentEpisode int    `json:"current_episode"`
	TotalEpisodes  int    `json:"total_episodes"`
	ImageURL       string `json:"image_url"`
	Rating         int    `json:"rating"`
}

// Estructura de la tabla para agregar series con rating.
type SeriesInput struct {
	Name           string `json:"name"`
	CurrentEpisode int    `json:"current_episode"`
	TotalEpisodes  int    `json:"total_episodes"`
	ImageURL       string `json:"image_url"`
	Rating         int    `json:"rating"`
}

// Activa las claves foraneas en SQlite.
func initDB(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	//Creacion de tablas
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS series (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			current_episode INTEGER NOT NULL DEFAULT 0,
			total_episodes INTEGER NOT NULL DEFAULT 0,
			image_url TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE IF NOT EXISTS ratings (
			series_id INTEGER PRIMARY KEY,
			rating INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	//Agrega columna para imagenes
	_, _ = db.Exec("ALTER TABLE series ADD COLUMN image_url TEXT NOT NULL DEFAULT ''")
	return nil
}

// Lista todas las series con SELECT Uniendo la tabla series y rating. Con paginacion
func listSeries(db *sql.DB, page, limit int) ([]Series, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	rows, err := db.Query(`
		SELECT s.id, s.name, s.current_episode, s.total_episodes,
		       COALESCE(s.image_url, '') AS image_url,
		       COALESCE(r.rating, 0) AS rating
		FROM series s
		LEFT JOIN ratings r ON r.series_id = s.id
		ORDER BY s.id DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	series := []Series{}
	for rows.Next() {
		item := Series{}
		if err := rows.Scan(&item.ID, &item.Name, &item.CurrentEpisode, &item.TotalEpisodes, &item.ImageURL, &item.Rating); err != nil {
			return nil, err
		}
		series = append(series, item)
	}

	return series, rows.Err()
}

func countSeries(db *sql.DB) (int, error) {
	var total int
	err := db.QueryRow(`SELECT COUNT(*) FROM series`).Scan(&total)
	return total, err
}

// Busca la serie por su ID
func getSeries(db *sql.DB, id int) (Series, error) {
	item := Series{}
	err := db.QueryRow(`
		SELECT s.id, s.name, s.current_episode, s.total_episodes,
		       COALESCE(s.image_url, '') AS image_url,
		       COALESCE(r.rating, 0) AS rating
		FROM series s
		LEFT JOIN ratings r ON r.series_id = s.id
		WHERE s.id = ?
	`, id).Scan(&item.ID, &item.Name, &item.CurrentEpisode, &item.TotalEpisodes, &item.ImageURL, &item.Rating)
	return item, err
}

// inserta nueva serie en la tabla con INSERT
func createSeries(db *sql.DB, input SeriesInput) (Series, error) {
	result, err := db.Exec(
		"INSERT INTO series (name, current_episode, total_episodes, image_url) VALUES (?, ?, ?, ?)",
		input.Name,
		input.CurrentEpisode,
		input.TotalEpisodes,
		input.ImageURL,
	)
	if err != nil {
		return Series{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Series{}, err
	}

	if err := upsertRating(db, int(id), input.Rating); err != nil {
		return Series{}, err
	}

	return getSeries(db, int(id))
}

// Actualiza la serie al cambiar cualquier atributo de este con UPDATE.
func updateSeries(db *sql.DB, id int, input SeriesInput) (Series, error) {
	result, err := db.Exec(
		"UPDATE series SET name = ?, current_episode = ?, total_episodes = ?, image_url = ? WHERE id = ?",
		input.Name,
		input.CurrentEpisode,
		input.TotalEpisodes,
		input.ImageURL,
		id,
	)
	if err != nil {
		return Series{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Series{}, err
	}
	if affected == 0 {
		return Series{}, sql.ErrNoRows
	}

	if err := upsertRating(db, id, input.Rating); err != nil {
		return Series{}, err
	}

	return getSeries(db, id)
}

// Elimina la serie de la tabla con DELETE.
func removeSeries(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Actualiza el rating en la tabla.
func upsertRating(db *sql.DB, seriesID int, rating int) error {
	_, err := db.Exec(
		`INSERT INTO ratings (series_id, rating)
		 VALUES (?, ?)
		 ON CONFLICT(series_id) DO UPDATE SET rating = excluded.rating`,
		seriesID,
		rating,
	)
	return err
}

// Verifica los datos antes de guardarlos en la base de datos.
func validateSeries(input SeriesInput) error {
	if input.Name == "" {
		return errors.New("Introduce un nombre valido")
	}
	if input.CurrentEpisode < 0 {
		return errors.New("El episodio actual debe ser mayor a 0")
	}
	if input.TotalEpisodes < 0 {
		return errors.New("El total de episodios debe ser mayor a 0")
	}
	if input.TotalEpisodes > 0 && input.CurrentEpisode > input.TotalEpisodes {
		return errors.New("El episodio actual no puede ser mayor al total de episodios")
	}
	if input.Rating < 0 || input.Rating > 10 {
		return errors.New("El rating debe ser entre 0 a 10")
	}
	return nil
}
