package postgresql

import (
	"context"
	"errors"
	"fmt"

	models "sso/internal/domain/models"
	"sso/internal/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

var storagePath = "host=localhost password=gopher user=postgres dbname=sso port=5432 sslmode=disable"

func New() (*Storage, error) {
	const op = "storage.PostgreSQL.New"
	db, err := gorm.Open(postgres.Open(storagePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.PostgreSQL.SaveUser"
	user := models.User{
		Email:    email,
		PassHash: passHash,
	}

	if err := s.db.Create(&user).Error; err != nil {
		if s.db.Error != nil && s.db.Error.Error() == "UNIQUE constraint failed: users.email" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int64(user.ID), nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.PostgreSQL.User"
	var user models.User

	if err := s.db.Select("id, email, pass_hash").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.PostgreSQL.IsAdmin"
	var user models.User

	if err := s.db.Select("is_admin").Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return user.IsAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.PostgreSQL.App"
	var app models.App

	if err := s.db.Select("id, name, secret").Where("id = ?", appID).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
