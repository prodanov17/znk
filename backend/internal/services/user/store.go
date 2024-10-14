package user

import (
	"database/sql"
	"fmt"

	"github.com/prodanov17/znk/internal/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) UpdateUser(user *types.User) (*types.User, error) {
	_, err := s.db.Exec("UPDATE users SET name = ?, email = ?, password = ?, phone_number = ?, user_type = ? WHERE id = ?", user.Name, user.Email, user.Password, user.PhoneNumber, user.UserType, user.ID)
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(user.ID)
}

func (s *Store) CreateUser(user *types.User) (*types.User, error) {
	res, err := s.db.Exec("INSERT INTO users (name, email, password, phone_number, user_type) VALUES (?,?,?,?,?)", user.Name, user.Email, user.Password, user.PhoneNumber, user.UserType)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(int(id))
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.UserType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
