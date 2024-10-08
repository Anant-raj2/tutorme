// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

const createTutor = `-- name: CreateTutor :one
INSERT INTO tutors (
  user_id, name,email, grade_level, gender, subject
) VALUES (
  gen_random_uuid(), $1, $2, $3, $4, $5
)
RETURNING user_id, name, email, grade_level, gender, subject
`

type CreateTutorParams struct {
	Name       string
	Email      string
	GradeLevel int32
	Gender     string
	Subject    string
}

func (q *Queries) CreateTutor(ctx context.Context, arg CreateTutorParams) (Tutor, error) {
	row := q.db.QueryRow(ctx, createTutor,
		arg.Name,
		arg.Email,
		arg.GradeLevel,
		arg.Gender,
		arg.Subject,
	)
	var i Tutor
	err := row.Scan(
		&i.UserID,
		&i.Name,
		&i.Email,
		&i.GradeLevel,
		&i.Gender,
		&i.Subject,
	)
	return i, err
}
