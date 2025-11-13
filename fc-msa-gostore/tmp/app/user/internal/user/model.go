package user

import "time"

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *UserResponse) FromUser(u *User) {
	r.ID = u.ID
	r.Name = u.Name
	r.Email = u.Email
	r.CreatedAt = u.CreatedAt
	r.UpdatedAt = u.UpdatedAt
}

type UserWithTokenResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
