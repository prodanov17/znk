package types

type RegisterUserPayload struct {
	Name                 string  `json:"name" validate:"required"`
	Email                string  `json:"email" validate:"required,email"`
	Password             string  `json:"password" validate:"required"`
	PasswordConfirmation string  `json:"password_confirmation" validate:"required"`
	PhoneNumber          *string `json:"phone_number"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserPayload struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" validate:"email"`
	Password    *string `json:"password"`
	UserType    *string `json:"user_type"`
	PhoneNumber *string `json:"phone_number"`
}

type CreateRoomPayload struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type CreateGamePayload struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}
