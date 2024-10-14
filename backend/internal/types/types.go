package types

import (
	"time"
)

type UserService interface {
	GetUserByID(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	LoginUser(userPayload *LoginUserPayload) (string, error)
	RegisterUser(userPayload *RegisterUserPayload) (string, error)
	UpdateUser(id int, userPayload *UpdateUserPayload) (*User, error)
}

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	UserType    string    `json:"user_type"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(user *User) (*User, error)
}

type Establishment struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EstablishmentService interface {
	GetEstablishments() ([]*Establishment, error)
	GetEstablishment(id int) (*Establishment, error)
	CreateEstablishment(e EstablishmentPayload) (*Establishment, error)
}

type EstablishmentStore interface {
	GetEstablishments() ([]*Establishment, error)
	GetEstablishment(id int) (*Establishment, error)
	CreateEstablishment(e *Establishment) (*Establishment, error)
}

type Barber struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	PhoneNumber    *string   `json:"phone_number"`
	ProfilePicture *string   `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type BarberService interface {
	GetBarbers() ([]*Barber, error)
	GetBarber(id int) (*Barber, error)
	CreateBarber(b BarberPayload) (*Barber, error)
}

type BarberStore interface {
	GetBarbers() ([]*Barber, error)
	GetBarber(id int) (*Barber, error)
	CreateBarber(b *Barber) (*Barber, error)
}
