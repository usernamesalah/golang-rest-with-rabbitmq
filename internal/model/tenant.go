package model

import (
	"time"
)

type Tenant struct {
	ID        int
	ClientID  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
