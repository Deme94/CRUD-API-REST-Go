package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type DBModels struct {
	DB *sql.DB
}

// GetAllGames returns one game and error, if any
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
		query = `SELECT gg.genre_id, g.genre_name
					FROM games_genres gg
						LEFT JOIN genres g on (g.id = gg.genre_id)
					WHERE gg.game_id = $1
				`

		rows, _ = m.DB.QueryContext(ctx, query, game.ID)

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
						LEFT JOIN modes m on (m.id = gm.mode_id)
					WHERE gm.game_id = $1
				`

		rows, _ = m.DB.QueryContext(ctx, query, game.ID)
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

		games = append(games, &game)
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
	query = `SELECT gg.genre_id, g.genre_name
				FROM games_genres gg
					LEFT JOIN genres g on (g.id = gg.genre_id)
				WHERE gg.game_id = $1
			`

	rows, _ := m.DB.QueryContext(ctx, query, id)

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
					LEFT JOIN modes m on (m.id = gm.mode_id)
				WHERE gm.game_id = $1
			`

	rows, _ = m.DB.QueryContext(ctx, query, id)
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
