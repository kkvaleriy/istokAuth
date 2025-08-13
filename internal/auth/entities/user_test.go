package user

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
)

func TestSignUp(t *testing.T) {
	tests := []struct {
		name    string
		dto     *dtos.CreateUserRequest
		want    *User
		wantErr string
	}{
		{
			name: "valid user",
			dto: &dtos.CreateUserRequest{
				Name:     "John",
				Lastname: "Doe",
				Nickname: "johndoe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    79123456789,
			},
			want: &User{
				Name:     "John",
				Lastname: "Doe",
				Nickname: "johndoe",
				Email:    "john@example.com",
				UserType: UserTypeUser,
				IsActive: true,
				Phone:    79123456789,
				PassHash: sha256.Sum256([]byte("password123")),
			},
			wantErr: "",
		},
		{
			name: "empty name",
			dto: &dtos.CreateUserRequest{
				Name:     "",
				Lastname: "Doe",
				Nickname: "johndoe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    79123456789,
			},
			want:    nil,
			wantErr: "the name of user is required",
		},
		{
			name: "empty lastname",
			dto: &dtos.CreateUserRequest{
				Name:     "John",
				Lastname: "",
				Nickname: "johndoe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    79123456789,
			},
			want:    nil,
			wantErr: "the lastname of user is required",
		},
		{
			name: "empty nickname",
			dto: &dtos.CreateUserRequest{
				Name:     "John",
				Lastname: "Doe",
				Nickname: "",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    79123456789,
			},
			want:    nil,
			wantErr: "the nickname of user is required",
		},
		{
			name: "invalid email",
			dto: &dtos.CreateUserRequest{
				Name:     "John",
				Lastname: "Doe",
				Nickname: "johndoe",
				Email:    "invalid",
				Password: "password123",
				Phone:    79123456789,
			},
			want:    nil,
			wantErr: "the correct email of user is required",
		},
		{
			name: "short password",
			dto: &dtos.CreateUserRequest{
				Name:     "John",
				Lastname: "Doe",
				Nickname: "johndoe",
				Email:    "john@example.com",
				Password: "pass",
				Phone:    79123456789,
			},
			want:    nil,
			wantErr: fmt.Sprintf("the password of user can't be empty or less then %v", minLenOfPassword),
		},
		{
			name: "invalid phone",
			dto: &dtos.CreateUserRequest{
				Name:     "John",
				Lastname: "Doe",
				Nickname: "johndoe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    123456789,
			},
			want:    nil,
			wantErr: "the correct email of user is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := SignUp(tt.dto)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("SignUp() unexpected error: %v", err)
				return
			}

			tt.want.UUID = got.UUID
			tt.want.CreatedAt = got.CreatedAt
			tt.want.UpdatedAt = got.UpdatedAt
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignUp() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSignIn(t *testing.T) {
	tests := []struct {
		name    string
		dto     *dtos.SignInRequest
		want    *User
		wantErr string
	}{
		{
			name: "valid email and password",
			dto: &dtos.SignInRequest{
				Email:    "john@example.com",
				Password: "password123",
				Phone:    0,
			},
			want: &User{
				Email:    "john@example.com",
				PassHash: sha256.Sum256([]byte("password123")),
			},
			wantErr: "",
		},
		{
			name: "valid phone and password",
			dto: &dtos.SignInRequest{
				Email:    "",
				Password: "password123",
				Phone:    79123456789,
			},
			want: &User{
				Phone:    79123456789,
				PassHash: sha256.Sum256([]byte("password123")),
			},
			wantErr: "",
		},
		{
			name: "invalid email and phone",
			dto: &dtos.SignInRequest{
				Email:    "invalid",
				Password: "password123",
				Phone:    123456789,
			},
			want:    nil,
			wantErr: "the correct email or phone number of user is required",
		},
		{
			name: "short password",
			dto: &dtos.SignInRequest{
				Email:    "john@example.com",
				Password: "pass",
				Phone:    0,
			},
			want:    nil,
			wantErr: fmt.Sprintf("the password of user can't be empty or less then %v", minLenOfPassword),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := SignIn(tt.dto)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("SignIn() unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignIn() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name string
		user *User
		ttl  time.Duration
	}{
		{
			name: "valid ttl",
			user: &User{
				UUID:     uuid.New(),
				Nickname: "johndoe",
			},
			ttl: 10 * time.Minute,
		},
		{
			name: "zero ttl",
			user: &User{
				UUID:     uuid.New(),
				Nickname: "johndoe",
			},
			ttl: 0,
		},
		{
			name: "negative ttl",
			user: &User{
				UUID:     uuid.New(),
				Nickname: "johndoe",
			},
			ttl: -10 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.user.RefreshToken(tt.ttl)
			if got.UUID == uuid.Nil {
				t.Errorf("RefreshToken() UUID is empty")
			}
			if got.UserUUID != tt.user.UUID {
				t.Errorf("RefreshToken() UserUUID = %v, want %v", got.UserUUID, tt.user.UUID)
			}
			if got.Nickname != tt.user.Nickname {
				t.Errorf("RefreshToken() Nickname = %v, want %v", got.Nickname, tt.user.Nickname)
			}

			expectedExpiresAt := time.Unix(got.CreatedAt, 0).Add(tt.ttl).Unix()
			if got.ExpiresAt != expectedExpiresAt {
				t.Errorf("RefreshToken() ExpiresAt = %v, want %v", got.ExpiresAt, expectedExpiresAt)
			}
		})
	}
}
