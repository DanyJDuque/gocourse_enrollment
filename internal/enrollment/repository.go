package enrollment

import (
	"context"
	"log"

	"github.com/DanyJDuque/gocourse_domain/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, enroll *domain.Enrollment) error
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(l *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: l,
		db:  db,
	}
}

func (r *repo) Create(ctx context.Context, enroll *domain.Enrollment) error {
	if err := r.db.WithContext(ctx).Create(enroll).Error; err != nil {
		r.log.Printf("error: %v", err)
		return err
	}
	r.log.Println("enrollment created with id: ", enroll.ID)
	return nil
}
