package user

import "github.com/google/uuid"

type RToken struct {
	UUID      uuid.UUID
	UserUUID  uuid.UUID
	Nickname  string
	CreatedAt int64
	ExpiresAt int64
}
