package user

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/prodanov17/znk/internal/services/auth"
	"github.com/prodanov17/znk/internal/types"
	"github.com/prodanov17/znk/internal/utils"
)

func TestUserUpdate(t *testing.T) {
	store := NewMockUserStore(&sql.DB{})
	service := NewService(store)

	existingUser := &types.User{
		ID:          1,
		Name:        "Jane Doe",
		Email:       "jane@doe.com",
		PhoneNumber: utils.StringPtr("1234567890"),
		Password:    "password",
		UserType:    "customer",
	}

	store.GetUserByIDFunc = func(id int) (*types.User, error) {
		if id != existingUser.ID {
			return nil, fmt.Errorf("user not found")
		}
		return existingUser, nil
	}

	store.UpdateUserFunc = func(user *types.User) (*types.User, error) {
		u, _ := store.GetUserByIDFunc(user.ID)

		u.Name = user.Name
		u.Email = user.Email
		u.PhoneNumber = user.PhoneNumber
		u.Password = user.Password

		return u, nil
	}

	t.Run("update user correctly", func(t *testing.T) {
		userPayload := &types.UpdateUserPayload{
			Name:        utils.StringPtr("John Doe"),
			Email:       utils.StringPtr("john@doe.com"),
			PhoneNumber: utils.StringPtr("1234567890"),
			Password:    utils.StringPtr("password"),
		}

		updatedUser, err := service.UpdateUser(existingUser.ID, userPayload)
		if err != nil {
			t.Errorf("UpdateUser returned an error: %v", err)
		}

		if updatedUser.Name != *userPayload.Name {
			t.Errorf("Expected name to be %s, got %s", *userPayload.Name, updatedUser.Name)
		}
		if updatedUser.Email != *userPayload.Email {
			t.Errorf("Expected email to be %s, got %s", *userPayload.Email, updatedUser.Email)
		}
		if *updatedUser.PhoneNumber != *userPayload.PhoneNumber {
			t.Errorf("Expected phone number to be %s, got %s", *userPayload.PhoneNumber, *updatedUser.PhoneNumber)
		}
		if !auth.ComparePasswords(updatedUser.Password, *userPayload.Password) {
			t.Errorf("Passwords do not match")
		}

	})

	t.Run("update user with invalid payload", func(t *testing.T) {
		userPayload := &types.UpdateUserPayload{
			Name:        utils.StringPtr("John Doe"),
			Email:       utils.StringPtr(""),
			PhoneNumber: utils.StringPtr("1234567890"),
			Password:    utils.StringPtr("password"),
		}

		_, err := service.UpdateUser(existingUser.ID, userPayload)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("update user with invalid id", func(t *testing.T) {
		userPayload := &types.UpdateUserPayload{
			Name:        utils.StringPtr("John Doe"),
			Email:       utils.StringPtr("john@doe.com"),
			PhoneNumber: utils.StringPtr("1234567890"),
			Password:    utils.StringPtr("password"),
		}

		_, err := service.UpdateUser(0, userPayload)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("partial update", func(t *testing.T) {
		userPayload := &types.UpdateUserPayload{
			Email: utils.StringPtr("newemail@doe.com"),
		}

		updatedUser, err := service.UpdateUser(existingUser.ID, userPayload)
		if err != nil {
			t.Errorf("UpdateUser returned an error: %v", err)
		}

		if updatedUser.Email != *userPayload.Email {
			t.Errorf("Expected email to be %s, got %s", *userPayload.Email, updatedUser.Email)
		}

		if updatedUser.Name != existingUser.Name {
			t.Errorf("Expected name to remain %s, got %s", existingUser.Name, updatedUser.Name)
		}

		if *updatedUser.PhoneNumber != *existingUser.PhoneNumber {
			t.Errorf("Expected phone number to remain %s, got %s", *existingUser.PhoneNumber, *updatedUser.PhoneNumber)
		}

		if updatedUser.Password != existingUser.Password {
			t.Errorf("Expected password to remain the same")
		}
	})

}

func TestRegisterUser(t *testing.T) {
	store := NewMockUserStore(&sql.DB{})
	service := NewService(store)

	existingUser := &types.User{
		ID:          1,
		Name:        "Jane Doe",
		Email:       "jane@doe.com",
		PhoneNumber: utils.StringPtr("1234567890"),
	}

	store.GetUserByEmailFunc = func(email string) (*types.User, error) {
		if email == existingUser.Email {
			return existingUser, nil
		}

		return nil, fmt.Errorf("user not found")
	}

	store.CreateUserFunc = func(user *types.User) (*types.User, error) {
		user.ID = 1
		return user, nil
	}

	tests := []struct {
		name        string
		userPayload *types.RegisterUserPayload
		wantErr     bool
	}{
		{name: "register user with invalid email", userPayload: &types.RegisterUserPayload{
			Name:                 "John Doe",
			Email:                "jjane@doe",
			Password:             "password",
			PasswordConfirmation: "password",
		}, wantErr: true},
		{name: "register user correctly", userPayload: &types.RegisterUserPayload{
			Name:                 "John Doe",
			Email:                "john@doe.com",
			Password:             "password",
			PasswordConfirmation: "password",
		}, wantErr: false},
		{name: "register with existing email", userPayload: &types.RegisterUserPayload{
			Name:                 "John Doe",
			Email:                "jane@doe.com",
			Password:             "password",
			PasswordConfirmation: "password",
		}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.RegisterUser(tt.userPayload)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

type mockUserStore struct {
	db *sql.DB

	UpdateUserFunc     func(user *types.User) (*types.User, error)
	GetUserByIDFunc    func(id int) (*types.User, error)
	CreateUserFunc     func(user *types.User) (*types.User, error)
	GetUserByEmailFunc func(email string) (*types.User, error)
}

func NewMockUserStore(db *sql.DB) *mockUserStore {
	return &mockUserStore{db: db}
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(id)
	}
	return nil, nil
}

func (m *mockUserStore) UpdateUser(user *types.User) (*types.User, error) {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(user)
	}
	return nil, nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(email)
	}
	return nil, nil
}

func (m *mockUserStore) CreateUser(user *types.User) (*types.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return nil, nil
}
