package entity

import "time"

type Block struct {
	Number    int
	Status    string
	UpdatedAt time.Time
}
