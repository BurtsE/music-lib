package models

import "time"

type Song struct {
	Group       string
	Text        string
	Link        string
	ReleaseDate time.Time
}
