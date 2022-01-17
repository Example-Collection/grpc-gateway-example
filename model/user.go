package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        string    `dynamo:"user_id,hash" json:"user_id,omitempty"`
	Name      string    `dynamo:"name,omitempty" json:"name,omitempty"`
	Nickname  string    `dynamo:"nickname,omitempty" json:"nickname,omitempty"`
	CreatedAt time.Time `dynamo:"created_at,omitempty" json:"created_at,omitempty"`
}

func (user *User) New() *User {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	return user
}
