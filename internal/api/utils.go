package api

import (
	"MusicLibrary/internal/models"
	"encoding/json"
	"fmt"
	"strings"
)

// splitText Делит песню на куплеты, считая,
// что они разделены пустыми строками, возвращая json
func splitText(text string) ([]byte, error) {
	res := make(map[string]interface{})
	verses := strings.SplitAfter(text, "\n\n")
	for i := range verses {
		res[fmt.Sprintf("verse_%d", i+1)] = verses[i]
	}
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// songFromConstructors Совмещает данные полученные из реквеста с обогащенной информацией из api
func songFromConstructors(song songConstructor, details detailsConstructor) models.Song {
	return models.Song{
		Name:        song.Name,
		GroupName:   song.Group,
		Text:        details.Text,
		Link:        details.Link,
		ReleaseDate: details.ReleaseDate.Time,
	}
}

func fillSongParams(params songParams) models.Song {
	song := models.Song{}
	if params.Name != nil {
		song.Name = *params.Name
	}
	if params.GroupName != nil {
		song.GroupName = *params.GroupName
	}
	if params.Text != nil {
		song.Text = *params.Text
	}
	if params.Link != nil {
		song.Link = *params.Link
	}
	if params.ReleaseDate != nil {
		song.ReleaseDate = *params.ReleaseDate
	}
	return song
}
