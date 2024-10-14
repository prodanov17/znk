package user

import (
	"fmt"

	"github.com/prodanov17/znk/internal/services/auth"
	"github.com/prodanov17/znk/internal/types"
	"github.com/prodanov17/znk/internal/utils"
)

type Service struct {
	store types.UserStore
}

func NewService(store types.UserStore) *Service {
	return &Service{store: store}
}

func (s *Service) GetUserByID(id int) (*types.User, error) {
	return s.store.GetUserByID(id)
}

func (s *Service) GetUserByEmail(email string) (*types.User, error) {
	return s.store.GetUserByEmail(email)
}

func (s *Service) UpdateUser(id int, userPayload *types.UpdateUserPayload) (*types.User, error) {
	if err := utils.ValidatePayload(userPayload); err != nil {
		return nil, err
	}

	user, err := s.store.GetUserByID(id)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if userPayload.Name != nil {
		user.Name = *userPayload.Name
	}
	if userPayload.Email != nil {
		user.Email = *userPayload.Email
	}
	if userPayload.PhoneNumber != nil {
		user.PhoneNumber = userPayload.PhoneNumber
	}
	if userPayload.Password != nil {
		hashedPassword, err := auth.HashPassword(*userPayload.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}
	if userPayload.UserType != nil {
		user.UserType = *userPayload.UserType
	}

	u, err := s.store.UpdateUser(user)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) RegisterUser(userPayload *types.RegisterUserPayload) (string, error) {
	if err := utils.ValidatePayload(userPayload); err != nil {
		return "", err
	}

	if userPayload.Password != userPayload.PasswordConfirmation {
		return "", fmt.Errorf("passwords do not match")
	}

	//check if user exists
	_, err := s.store.GetUserByEmail(userPayload.Email)
	if err == nil {
		return "", fmt.Errorf("user with email %s already exists", userPayload.Email)
	}

	hash, err := auth.HashPassword(userPayload.Password)

	if err != nil {
		return "", nil
	}

	user := types.User{
		Name:        userPayload.Name,
		Email:       userPayload.Email,
		Password:    hash,
		PhoneNumber: userPayload.PhoneNumber,
		UserType:    "customer",
	}

	u, err := s.store.CreateUser(&user)

	if err != nil {
		return "", err
	}

	return auth.CreateToken(u.ID)
}

func (s *Service) LoginUser(userPayload *types.LoginUserPayload) (string, error) {
	if err := utils.ValidatePayload(userPayload); err != nil {
		return "", err
	}

	user, err := s.store.GetUserByEmail(userPayload.Email)

	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if !auth.ComparePasswords(user.Password, userPayload.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	return auth.CreateToken(user.ID)
}
