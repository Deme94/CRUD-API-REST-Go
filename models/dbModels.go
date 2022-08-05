package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type DBModels struct {
	DB *sql.DB
}

// GetAllGenres returns all genres and error, if any
func (m *DBModels) GetAllGenres() (map[int]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, genre_name, created_at, updated_at 
				FROM genres
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make(map[int]string)
	for rows.Next() {
		var g Genre
		err := rows.Scan(
			&g.ID,
			&g.GenreName,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres[g.ID] = g.GenreName
	}

	return genres, nil
}

// GetAllModes returns all modes and error, if any
func (m *DBModels) GetAllModes() (map[int]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, mode_name, created_at, updated_at 
				FROM modes
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modes := make(map[int]string)
	for rows.Next() {
		var m Mode
		err := rows.Scan(
			&m.ID,
			&m.ModeName,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		modes[m.ID] = m.ModeName
	}

	return modes, nil
}

// GetAllGames returns all games and error, if any
func (m *DBModels) GetAllGames() ([]*Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, title, image_url, developers, publishers, release_date, storage, likes, created_at, updated_at 
				FROM games
				ORDER BY title
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games, err := m.getGamesFromRows(ctx, rows)
	if err != nil {
		return nil, err
	}

	return games, nil
}

// GetOneGame returns one game and error, if any
func (m *DBModels) GetOneGame(id int) (*Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, title, image_url, developers, publishers, release_date, storage, likes, created_at, updated_at 
				FROM games 
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	game, err := m.getGameFromRow(ctx, row, id)
	if err != nil {
		return nil, err
	}

	return game, nil
}

// GetGameImage returns image from game id and error, if any
func (m *DBModels) GetGameImage(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT image_url
				FROM games 
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var image string
	err := row.Scan(
		&image,
	)
	if err != nil {
		return "", err
	}

	return image, nil
}

// GetAllImages returns all images and error, if any
func (m *DBModels) GetAllImages() (map[int]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, image_url
				FROM games
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make(map[int]string)
	for rows.Next() {

		var id int
		var imageUrl string
		err := rows.Scan(
			&id,
			&imageUrl,
		)
		if err != nil {
			return nil, err
		}

		images[id] = imageUrl
	}

	return images, nil
}

// GetAllGamesByGenre returns all games of a certain genre and error, if any
func (m *DBModels) GetAllGamesByGenre(genreID int) ([]*Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	gameIDsQuery := fmt.Sprintf("SELECT gg.game_id FROM games_genres gg"+
		" LEFT JOIN genres g ON (g.id = gg.genre_id)"+
		" WHERE g.genre_id = '%d'", genreID)

	query := fmt.Sprintf("SELECT id, title, image_url, developers, publishers, release_date, storage, likes, created_at, updated_at"+
		" FROM games"+
		" WHERE id IN (%s)", gameIDsQuery)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games, err := m.getGamesFromRows(ctx, rows)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (m *DBModels) InsertGame(game *Game, genres []int, modes []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert game
	query := fmt.Sprintf(`INSERT INTO games (title, image_url, developers, publishers, release_date,`+
		` storage, likes, created_at, updated_at)`+
		` VALUES ('%s', '%s', $1, $2, $3, %d, %d, NOW(), NOW())`, game.Title, game.ImageUrl, game.Storage, game.Likes)

	_, err := m.DB.ExecContext(ctx, query, pq.Array(game.Publishers), pq.Array(game.Developers), game.ReleaseDate.UTC().Format("2006-01-02"))
	if err != nil {
		return err
	}
	fmt.Println("INSERTED")
	// Get ID of the new game inserted
	var gameID int
	query = fmt.Sprintf(`SELECT id FROM games`+
		` WHERE title = '%s'`, game.Title)

	row := m.DB.QueryRowContext(ctx, query)
	err = row.Scan(
		&gameID,
	)
	if err != nil {
		return err
	}

	// Insert game_genres
	for _, genreID := range genres {
		query = fmt.Sprintf(`INSERT INTO games_genres (game_id, genre_id, created_at, updated_at)`+
			` VALUES (%d, %d, NOW(), NOW())`, gameID, genreID)

		_, err = m.DB.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}

	// Insert game_modes
	for _, modeID := range modes {
		query = fmt.Sprintf(`INSERT INTO games_modes (game_id, mode_id, created_at, updated_at)`+
			` VALUES (%d, %d, NOW(), NOW())`, gameID, modeID)

		_, err = m.DB.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}
	return err
}

func (m *DBModels) UpdateGame(id int, game *Game, genres []int, modes []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert game
	query := fmt.Sprintf(`UPDATE games SET title = '%s', image_url = '%s', developers = $1, publishers = $2, release_date = $3,`+
		` storage = %d, likes = %d, updated_at = NOW()`, game.Title, game.ImageUrl, game.Storage, game.Likes)

	_, err := m.DB.ExecContext(ctx, query, pq.Array(game.Publishers), pq.Array(game.Developers), game.ReleaseDate.UTC().Format("2006-01-02"))
	if err != nil {
		return err
	}

	// Delete and Insert game_genres
	query = `DELETE FROM games_genres
			WHERE game_id = $1;
			`

	_, err = m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	for _, genreID := range genres {
		query = fmt.Sprintf(`INSERT INTO games_genres (game_id, genre_id, created_at, updated_at)`+
			` VALUES (%d, %d, NOW(), NOW()) ON CONFLICT DO NOTHING`, id, genreID)

		_, err = m.DB.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}

	// Delete and Insert game_modes
	query = `DELETE FROM games_modes
			WHERE game_id = $1;
			`

	_, err = m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	for _, modeID := range modes {
		query = fmt.Sprintf(`INSERT INTO games_modes (game_id, mode_id, created_at, updated_at)`+
			` VALUES (%d, %d, NOW(), NOW()) ON CONFLICT DO NOTHING`, id, modeID)

		_, err = m.DB.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}

	fmt.Println("UPDATED")

	return err
}

func (m *DBModels) DeleteGame(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Delete game
	query := `DELETE FROM games
			WHERE id = $1;
			`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// Reusable private function
func (m *DBModels) getGamesFromRows(ctx context.Context, rows *sql.Rows) ([]*Game, error) {
	var games []*Game
	for rows.Next() {

		var game Game
		err := rows.Scan(
			&game.ID,
			&game.Title,
			&game.ImageUrl,
			pq.Array(&game.Developers),
			pq.Array(&game.Publishers),
			&game.ReleaseDate,
			&game.Storage,
			&game.Likes,
			&game.CreatedAt,
			&game.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// get genres, if any
		query := `SELECT gg.genre_id, g.genre_name
					FROM games_genres gg
						LEFT JOIN genres g ON (g.id = gg.genre_id)
					WHERE gg.game_id = $1
				`

		rows2, err := m.DB.QueryContext(ctx, query, game.ID)
		if err != nil {
			return nil, err
		}

		genres := make(map[int]string)
		for rows2.Next() {
			var g Genre
			err := rows2.Scan(
				&g.ID,
				&g.GenreName,
			)
			if err != nil {
				return nil, err
			}
			genres[g.ID] = g.GenreName
		}
		game.Genres = genres

		// get modes, if any
		query = `SELECT gm.mode_id, m.mode_name
					FROM games_modes gm
						LEFT JOIN modes m ON (m.id = gm.mode_id)
					WHERE gm.game_id = $1
				`

		rows2, err = m.DB.QueryContext(ctx, query, game.ID)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()

		modes := make(map[int]string)
		for rows2.Next() {
			var m Mode
			err := rows2.Scan(
				&m.ID,
				&m.ModeName,
			)
			if err != nil {
				return nil, err
			}
			modes[m.ID] = m.ModeName
		}
		game.Modes = modes
		games = append(games, &game)
	}

	return games, nil
}

// Reusable private function
func (m *DBModels) getGameFromRow(ctx context.Context, row *sql.Row, id int) (*Game, error) {
	var game Game

	err := row.Scan(
		&game.ID,
		&game.Title,
		&game.ImageUrl,
		pq.Array(&game.Developers),
		pq.Array(&game.Publishers),
		&game.ReleaseDate,
		&game.Storage,
		&game.Likes,
		&game.CreatedAt,
		&game.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// get genres, if any
	query := `SELECT gg.genre_id, g.genre_name
				FROM games_genres gg
					LEFT JOIN genres g ON (g.id = gg.genre_id)
				WHERE gg.game_id = $1
			`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	genres := make(map[int]string)
	for rows.Next() {
		var g Genre
		err := rows.Scan(
			&g.ID,
			&g.GenreName,
		)
		if err != nil {
			return nil, err
		}
		genres[g.ID] = g.GenreName
	}
	game.Genres = genres

	// get modes, if any
	query = `SELECT gm.mode_id, m.mode_name
				FROM games_modes gm
					LEFT JOIN modes m ON (m.id = gm.mode_id)
				WHERE gm.game_id = $1
			`

	rows, err = m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modes := make(map[int]string)
	for rows.Next() {
		var m Mode
		err := rows.Scan(
			&m.ID,
			&m.ModeName,
		)
		if err != nil {
			return nil, err
		}
		modes[m.ID] = m.ModeName
	}
	game.Modes = modes

	return &game, nil
}
