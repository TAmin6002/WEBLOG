package repo

import (
	"database/sql"
	"weblog/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(u *sql.DB) *UserRepo {
	return &UserRepo{db: u}
}

func (r *UserRepo) AddUser(user *model.User) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := r.db.Exec(query, user.Username, user.Password)

	return err
}

func (r *UserRepo) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, id)

	return err
}

func (r *UserRepo) FindUserByID(id int) (*model.User, error) {
	user := &model.User{}

	query := "SELECT id, username, password, created_at FROM users WHERE id = $1"

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) FindUser(username string) (*model.User, bool, error) {
	user := &model.User{}

	query := "SELECT id, username, password, created_at FROM users WHERE username = $1"

	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}

func (r *UserRepo) FindUserIDsByUsernames(usernames []string) (map[string]int, []string, error) {
	found := make(map[string]int)
	var missing []string

	for _, uname := range usernames {
		user, ok, err := r.FindUser(uname)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			missing = append(missing, uname)
			continue
		}
		found[uname] = user.ID
	}

	return found, missing, nil
}