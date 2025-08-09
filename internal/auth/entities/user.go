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
	if dto.Phone < 7_000_000_00_00 || dto.Phone > 8_999_999_99_99 {
		return nil, errors.New("the correct email of user is required")
	}
	if err := isValidePassword(dto.Password); err != nil {
		return nil, err
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
	if err := isValidePassword(dto.Password); err != nil {
		return nil, err
	}
	return &User{Email: dto.Email, Phone: dto.Phone, PassHash: sha256.Sum256([]byte(dto.Password))}, nil
}

func (u *User) RefreshToken(ttl time.Duration) *RToken {
	issuedAt := time.Now().Unix()
	expiresAt := time.Unix(issuedAt, 0).Add(ttl).Unix()
	return &RToken{
		UUID:      uuid.New(),
		UserUUID:  u.UUID,
		Nickname:  u.Nickname,
		CreatedAt: issuedAt,
		ExpiresAt: expiresAt,
	}
}

func Empty() *User {
	return &User{}
}


func isValidePassword(pass string) error {
	if len(strings.TrimSpace(pass)) < minLenOfPassword {
		return fmt.Errorf("the password of user can't be empty or less then %v", minLenOfPassword)
	}
	return nil
}
