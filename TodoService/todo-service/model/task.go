package model

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description,omitempty"`
	Creator     string         `json:"creator"`
	Assignee    string         `json:"assignee"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	StartedAt   sql.NullTime   `json:"started_at,omitempty"`
	CompletedAt sql.NullTime   `json:"completed_at,omitempty"`
	DeletedAt   sql.NullTime   `json:"deleted_at,omitempty"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Creator     string `json:"creator" binding:"required"`
	Assignee    string `json:"assignee" binding:"required"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_progress completed deleted"`
}
