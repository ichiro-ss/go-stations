package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	var todo model.TODO
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// if subject == "" {
	// 	return nil, err.(sqlite3.ErrError)
	// }

	res, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}
	created_id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	row := s.db.QueryRowContext(ctx, confirm, created_id)

	if err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}
	todo.ID = created_id
	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	var todos []*model.TODO
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	var rows *sql.Rows
	var err error
	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()
	for rows.Next() {
		var todo model.TODO
		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	if len(todos) == 0 {
		return []*model.TODO{}, nil
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	var todo model.TODO
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}
	if num, err := res.RowsAffected(); num == 0 || err != nil {
		return nil, &model.ErrNotFound{time.Now(), "not found"}
	}
	row := s.db.QueryRowContext(ctx, confirm, id)
	if err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}
	todo.ID = id

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
