package model

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
