package handlers

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Error string `json:"error"`
}

type SubscribeResponse struct {
	Status bool `json:"status"`
}

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *AuthRequest) Validate() error {
	return validator.New().Struct(r)
}

type AuthResponse struct {
	Payload string `json:"payload"`
}

type SumResponse struct {
	Sum  int    `json:"sum"`
	Hash string `json:"hash"`
}
