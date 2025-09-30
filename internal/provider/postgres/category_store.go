package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

var (
	//go:embed queries/get_category_by_id.sql
	getCategoryByIDQuery string

	//go:embed queries/get_categories.sql
	getCategoriesQuery string
)

type CategoryStore struct {
	db *pgxpool.Pool
}

func NewCategoryStore(db *pgxpool.Pool) *CategoryStore {
	return &CategoryStore{
		db: db,
	}
}

func (s *CategoryStore) GetByID(ctx context.Context, id int64) (model.Category, error) {
	row := s.db.QueryRow(ctx, getCategoryByIDQuery, id)

	category := model.Category{}

	err := row.Scan(&category.ID, &category.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Category{}, fmt.Errorf("category not found: %w", model.ErrNotFound)
	}
	if err != nil {
		return model.Category{}, fmt.Errorf("select category by id: %w", err)
	}

	return category, nil
}

func (s *CategoryStore) GetCategories(ctx context.Context) ([]model.Category, error) {
	rows, err := s.db.Query(ctx, getCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("select categories: %w", err)
	}

	categories, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Category, error) {
		var category model.Category

		err := row.Scan(&category.ID, &category.Name)
		if err != nil {
			return model.Category{}, err
		}

		return category, nil
	})
	if err != nil {
		return nil, fmt.Errorf("select categories: %w", err)
	}

	return categories, nil
}
