package user

import (
	"account/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}

func (r *repo) FindOneByPhone(c context.Context, phone string) (*domain.User, error) {
	var u domain.User

	var role domain.Role

	sql := `
	SELECT id, username, phone, password, photo, description, created_at, role_id,
		(SELECT name FROM roles WHERE id = role_id) as role_name,
		(SELECT permissions FROM roles WHERE id = role_id) as role_permissions
	FROM users WHERE phone = $1
	`

	row := r.db.QueryRow(c, sql, phone)
	err := row.Scan(
		&u.ID, &u.Username, &u.Phone, &u.Password,
		&u.Photo, &u.Description, &u.CreatedAt,
		&role.ID, &role.Name, &role.Permissions,
	)
	if err != nil {
		return nil, err
	}

	u.Role = &role

	return &u, nil
}

func (r *repo) FindOneByUsername(c context.Context, username string) (*domain.User, error) {
	var u domain.User

	var role domain.Role

	sql := `
	SELECT id, username, phone, password, photo, description, created_at, role_id,
		(SELECT name FROM roles WHERE id = role_id) as role_name,
		(SELECT permissions FROM roles WHERE id = role_id) as role_permissions
	FROM users WHERE username = $1
	`

	row := r.db.QueryRow(c, sql, username)
	err := row.Scan(
		&u.ID, &u.Username, &u.Phone, &u.Password,
		&u.Photo, &u.Description, &u.CreatedAt,
		&role.ID, &role.Name, &role.Permissions,
	)
	if err != nil {
		return nil, err
	}

	u.Role = &role

	return &u, nil
}

func (r *repo) FindOneByID(c context.Context, userId int) (*domain.User, error) {
	var u domain.User

	var role domain.Role

	sql := `
	SELECT id, username, phone, password, photo, description, created_at, role_id,
		(SELECT name FROM roles WHERE id = role_id) as role_name,
		(SELECT permissions FROM roles WHERE id = role_id) as role_permissions
	FROM users WHERE id = $1
	`

	row := r.db.QueryRow(c, sql, userId)
	err := row.Scan(
		&u.ID, &u.Username, &u.Phone, &u.Password,
		&u.Photo, &u.Description, &u.CreatedAt,
		&role.ID, &role.Name, &role.Permissions,
	)
	if err != nil {
		return nil, err
	}

	u.Role = &role

	return &u, nil
}

func (r *repo) Create(c context.Context, user *domain.User) (int, error) {
	var userId int
	sql := "INSERT INTO users (username, phone, password) VALUES ($1,$2,$3) RETURNING id;"
	err := r.db.QueryRow(c, sql, user.Username, user.Phone, user.Password).Scan(&userId)
	return userId, err
}

func (r *repo) Save(c context.Context, user *domain.User) error {
	sql := `
		UPDATE users SET id = $1, username = $2, phone = $3,
			password = $4, photo = $5, description = $6, role_id = $7
		WHERE id = $1;
		`

	_, err := r.db.Exec(c, sql, user.ID, user.Username, user.Phone, user.Password, user.Photo, user.Description, user.Role.ID)
	return err
}
