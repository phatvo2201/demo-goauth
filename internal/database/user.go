package database

import "database/sql"

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := "SELECT id, email, password, role FROM users WHERE email = $1"
	row := db.QueryRow(query, email)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

type User struct {
	ID       int
	Email    string
	Password string
	Role     string
}
