package api

import (
	"fmt"
	"time"
)

type songConstructor struct {
	Group string `json:"group"`
	Name  string `json:"name"`
}

type detailsConstructor struct {
	Text        string     `json:"text"`
	Link        string     `json:"link"`
	ReleaseDate CustomTime `json:"releaseDate"`
}

const customLayout = "02.01.2006"

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return fmt.Errorf("invalid time format: %s", str)
	}
	str = str[1 : len(str)-1]
	parsedTime, err := time.Parse(customLayout, str)
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}

type songParams struct {
	Name        *string    `json:"name"`
	GroupName   *string    `json:"group"`
	Text        *string    `json:"text"`
	Link        *string    `json:"link"`
	ReleaseDate *time.Time `json:"release_date"`
}
