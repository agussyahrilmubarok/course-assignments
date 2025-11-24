package user

import "time"

type SignUpParam struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ur *UserResponse) From(user *User) {
	ur.ID = user.ID.Hex()
	ur.Name = user.Name
	ur.Email = user.Email
	ur.CreatedAt = user.CreatedAt
	ur.UpdatedAt = user.UpdatedAt
}

type UserWithTokenResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
