package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kkvaleriy/istokAuthorization/internal/authorization/dtos"
)

const (
	UserTypeUser  = "USER"
	UserTypeAdmin = "ADMIN"
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
}
