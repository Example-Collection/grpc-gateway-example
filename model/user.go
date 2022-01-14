package model

import "time"

type User struct {
	ID        string    `dynamo:"user_id,hash" json:"user_id,omitempty"`
	Name      string    `dynamo:"name,omitempty" json:"name,omitempty"`
	Nickname  string    `dynamo:"nickname,omitempty" json:"nickname,omitempty"`
	CreatedAt time.Time `dynamo:"created_at,omitempty" json:"created_at,omitempty"`
}
