package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
)

const (
	UserTypeUser     = "USER"
	UserTypeAdmin    = "ADMIN"
	minLenOfPassword = 8
)

type User struct {
	Name      string
	Lastname  string
	Nickname  string
	Email     string
	UserType  string
	IsActive  bool
	Phone     int
	UUID      uuid.UUID
	PassHash  [32]byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func SignUp(dto *dtos.CreateUserRequest) (*User, error) {

	if len(strings.TrimSpace(dto.Name)) < 1 {
		return nil, errors.New("the name of user is required")
	}
	if len(strings.TrimSpace(dto.Lastname)) < 1 {
		return nil, errors.New("the lastname of user is required")
	}
	if len(strings.TrimSpace(dto.Nickname)) < 1 {
		return nil, errors.New("the nickname of user is required")
	}
	if len(strings.TrimSpace(dto.Email)) < 1 || !strings.Contains(dto.Email, "@") {
		return nil, errors.New("the correct email of user is required")
	}
	if len(strings.TrimSpace(dto.Password)) < minLenOfPassword {
		return nil, fmt.Errorf("the password of user must be more than %v characters long", minLenOfPassword)
	}
	if dto.Phone < 7_000_000_00_00 || dto.Phone > 8_999_999_99_99 {
		return nil, errors.New("the correct email of user is required")
	}

	return &User{
		Name:      dto.Name,
		Lastname:  dto.Lastname,
		Nickname:  dto.Nickname,
		Email:     dto.Email,
		UserType:  UserTypeUser,
		IsActive:  true,
		Phone:     dto.Phone,
		UUID:      uuid.New(),
		PassHash:  sha256.Sum256([]byte(dto.Password)),
		CreatedAt: time.Now(),
	}, nil
}

func SignIn(dto *dtos.SignInRequest) (*User, error) {
	if (len(strings.TrimSpace(dto.Email)) < 1 || !strings.Contains(dto.Email, "@")) && dto.Phone < 7_000_000_00_00 {
		return nil, errors.New("the correct email or phone number of user is required")
	}
	if len(strings.TrimSpace(dto.Password)) < minLenOfPassword {
		return nil, fmt.Errorf("the password of user can't be empty or less then %v", minLenOfPassword)
	}
	return &User{Email: dto.Email, Phone: dto.Phone, PassHash: sha256.Sum256([]byte(dto.Password))}, nil
}
