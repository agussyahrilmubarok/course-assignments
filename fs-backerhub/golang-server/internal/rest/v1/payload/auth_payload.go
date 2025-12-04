package payloadV1

type SignUpRequest struct {
	Name       string `json:"name" binding:"required,min=3,max=100"`
	Occupation string `json:"occupation" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6,max=100"`
}

type SignUpResponse struct {
	ID string `json:"id"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignInResponse struct {
	Token string `json:"token"`
}
