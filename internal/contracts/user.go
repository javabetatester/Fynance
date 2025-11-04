package contracts

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Plan     string `json:"plan" binding:"omitempty,oneof=FREE BASIC PRO"`
}

type UserUpdateRequest struct {
	Name  string `json:"name" binding:"omitempty"`
	Email string `json:"email" binding:"omitempty,email"`
	Plan  string `json:"plan" binding:"omitempty,oneof=FREE BASIC PRO"`
}

type UserDeletionResponse struct {
	Message string `json:"message"`
}
