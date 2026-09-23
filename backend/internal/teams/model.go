package teams

import "time"

type Team struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Interests    []string  `json:"interests"`
	Skills       []string  `json:"skills"`
	Technologies []string  `json:"technologies"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
