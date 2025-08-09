package dtos

// @Description New user password.
type UpdateUserPasswordRequest struct {
	ID       *UUIDRequest `json:"-"`
	Password string       `json:"password" example:"mySuperPass" validate:"required,min=8"`
}
