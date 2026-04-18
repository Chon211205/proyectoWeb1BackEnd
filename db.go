package main

import (
	"database/sql"
	"errors"
	"strings"
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
}

//Estructura de rating.
type RatingResponse struct {
	SeriesID int `json:"series_id"`
	Rating   int `json:"rating"`
}

//
type RatingInput struct {
	Rating int `json:"rating"`
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

// Lista todas las series con SELECT Uniendo la tabla series y rating. Con paginacion y buscar con ?q=
func listSeries(db *sql.DB, page, limit int, q, sort, order string) ([]Series, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := `
		SELECT s.id, s.name, s.current_episode, s.total_episodes,
		       COALESCE(s.image_url, '') AS image_url,
		       COALESCE(r.rating, 0) AS rating
		FROM series s
		LEFT JOIN ratings r ON r.series_id = s.id
	`

	args := []any{}

	//Filtro
	if q != "" {
		query += " WHERE LOWER(s.name) LIKE ?"
		args = append(args, "%"+strings.ToLower(q)+"%")
	}

	//Ordenar
	validSort := map[string]string{
		"name":            "s.name",
		"rating":          "r.rating",
		"current_episode": "s.current_episode",
		"total_episodes":  "s.total_episodes",
	}

	sortColumn, ok := validSort[sort]
	if !ok {
		sortColumn = "s.id"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query += " ORDER BY " + sortColumn + " " + strings.ToUpper(order)

	//Paginacion
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
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

func validateRating(input RatingInput) error {
	if input.Rating < 0 || input.Rating > 10 {
		return errors.New("rating must be between 0 and 10")
	}
	return nil
}

func getRating(db *sql.DB, seriesID int) (RatingResponse, error) {
	item := RatingResponse{}

	err := db.QueryRow(`
		SELECT series_id, rating
		FROM ratings
		WHERE series_id = ?
	`, seriesID).Scan(&item.SeriesID, &item.Rating)

	return item, err
}

func upsertRating(db *sql.DB, seriesID int, rating int) (RatingResponse, error) {
	_, err := db.Exec(`
		INSERT INTO ratings (series_id, rating)
		VALUES (?, ?)
		ON CONFLICT(series_id) DO UPDATE SET rating = excluded.rating
	`, seriesID, rating)
	if err != nil {
		return RatingResponse{}, err
	}

	return getRating(db, seriesID)
}

func removeRating(db *sql.DB, seriesID int) error {
	result, err := db.Exec(`
		DELETE FROM ratings
		WHERE series_id = ?
	`, seriesID)
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



func countSeries(db *sql.DB, q string) (int, error) {
	var total int

	query := "SELECT COUNT(*) FROM series"
	args := []any{}

	if q != "" {
		query += " WHERE LOWER(name) LIKE ?"
		args = append(args, "%"+strings.ToLower(q)+"%")
	}

	err := db.QueryRow(query, args...).Scan(&total)
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

// Verifica los datos antes de guardarlos en la base de datos.
func validateSeries(input SeriesInput) error {
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.CurrentEpisode < 0 {
		return errors.New("current_episode must be greater than or equal to 0")
	}
	if input.TotalEpisodes < 0 {
		return errors.New("total_episodes must be greater than or equal to 0")
	}
	if input.TotalEpisodes > 0 && input.CurrentEpisode > input.TotalEpisodes {
		return errors.New("current_episode cannot be greater than total_episodes")
	}
	return nil
}
