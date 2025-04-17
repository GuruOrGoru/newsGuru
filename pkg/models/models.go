package models

import (
	"context"
	"time"

	"github.com/guruorgoru/newsguru/pkg/logs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type News struct {
	NewsID      int64     `db:"news_id"`
	Title       string    `db:"title"`
	Body        string    `db:"body"`
	AuthorName  string    `db:"author_name"`
	Category    string    `db:"category"`
	PublishedAt time.Time `db:"published_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
type NewsModel struct {
	DB *pgxpool.Pool
}

func (n *NewsModel) Insert(news News) (int64, error) {
	query := `INSERT INTO news (title, body, author_name, category, published_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6)
            RETURNING news_id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newsID int64
	err := n.DB.QueryRow(ctx, query,
		news.Title,
		news.Body,
		news.AuthorName,
		news.Category,
		news.PublishedAt,
		news.UpdatedAt,
	).Scan(&newsID)

	if err != nil {
		return 0, err
	}

	return newsID, nil
}

func (n *NewsModel) GetByID(id int64) (*News, error) {
	query := `SELECT news_id, title, body, author_name, category, published_at, updated_at
            FROM news WHERE news_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	news := &News{}

	err := n.DB.QueryRow(ctx, query, id).Scan(
		&news.NewsID,
		&news.Title,
		&news.Body,
		&news.AuthorName,
		&news.Category,
		&news.PublishedAt,
		&news.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (n *NewsModel) GetAll() ([]*News, error) {
	query := `SELECT news_id, title, body, author_name, category, published_at, updated_at
            FROM news`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := n.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newss := []*News{}

	for rows.Next() {
		news := &News{}

		err := rows.Scan(
			&news.NewsID,
			&news.Title,
			&news.Body,
			&news.AuthorName,
			&news.Category,
			&news.PublishedAt,
			&news.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		newss = append(newss, news)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return newss, err
}
func (n *NewsModel) Delete(id int) error {
	query := `DELETE FROM news WHERE news_id=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd, err := n.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return logs.SErrorNotFound
	}
	return nil
}
