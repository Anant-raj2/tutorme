
-- name: CreateTutor :one
INSERT INTO tutors (
  user_id, name,email, grade_level, role, gender, subject
) VALUES (
  gen_random_uuid(), $1, $2, $3, $4, $5, $6
)
RETURNING *;
