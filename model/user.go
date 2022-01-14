package model

type User struct {
	ID       string `dynamo:"user_id,hash" json:"user_id,omitempty"`
	Name     string `dynamo:"nickname" json:"name,omitempty"`
	Nickname string `dynamo:"nickname,omitempty" json:"nickname,omitempty"`
}
