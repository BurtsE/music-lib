package database

import (
	"MusicLibrary/internal/config"
	"MusicLibrary/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Postgres struct {
	db        *gorm.DB
	pageLimit int
}

func InitDB(cfg *config.Config) (*Postgres, error) {
	DSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Hostname,
		cfg.Postgres.Port,
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	p := &Postgres{
		db:        db,
		pageLimit: cfg.Service.PageLimit,
	}
	if err = p.migrate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Postgres) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (p *Postgres) migrate() error {
	return p.db.AutoMigrate(&models.Song{})
}

func (p *Postgres) CreateSong(song models.Song) error {
	return p.db.Create(&song).Error
}

func (p *Postgres) GetSong(songId uint) (models.Song, error) {
	song := models.Song{ID: songId}
	tx := p.db.First(&song)
	if tx.Error != nil {
		return models.Song{}, tx.Error
	}
	return song, nil
}

func (p *Postgres) GetSongs(group, name string, pageNumber int) (songs []models.Song, err error) {
	name = name + "%"
	group = group + "%"
	tx := p.db.Order("songs.name").
		Where("name LIKE ? AND group_name LIKE ?", name, group).
		Limit(p.pageLimit).
		Offset(p.pageLimit * pageNumber).
		Find(&songs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return
}

func (p *Postgres) DeleteSong(songId uint) error {
	return p.db.Delete(&models.Song{ID: songId}, songId).Error
}

func (p *Postgres) UpdateSong(song models.Song) error {
	fields := map[string]interface{}{}
	if song.Name != "" {
		fields["name"] = song.Name
	}
	if song.GroupName != "" {
		fields["group_name"] = song.GroupName
	}
	if song.Text != "" {
		fields["text"] = song.Text
	}
	if !song.ReleaseDate.Equal(time.Time{}) {
		fields["release_date"] = song.ReleaseDate
	}
	if song.Link != "" {
		fields["link"] = song.Link
	}
	return p.db.Model(&song).Updates(fields).Error
}
