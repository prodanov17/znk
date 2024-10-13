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

type EstablishmentPayload struct {
	Name        string  `json:"name" validate:"required"`
	Address     string  `json:"address" validate:"required"`
	PhoneNumber *string `json:"phone_number"`
}

type BarberPayload struct {
	UserID         int     `json:"user_id" validate:"required"`
	ProfilePicture *string `json:"profile_picture"`
}
