// Code generated by sqlc. DO NOT EDIT.

package database

import (
	"time"

	"github.com/jackc/pgtype"
)

type Product struct {
	ID         pgtype.UUID `json:"id"`
	CreateTime time.Time   `json:"create_time"`
	Name       string      `json:"name"`
	Kind       string      `json:"kind"`
}
