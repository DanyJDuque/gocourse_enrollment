package enrollment_test

import (
	"context"

	"github.com/DanyJDuque/gocourse_domain/domain"
	"github.com/DanyJDuque/gocourse_enrollment/internal/enrollment"
)

type mockRespository struct {
	CreateMock func(ctx context.Context, erroll *domain.Enrollment) error
	GetAllMock func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error)
	UpdateMock func(ctx context.Context, id string, status *string) error
	CountMock  func(ctx context.Context, filters enrollment.Filters) (int, error)
}

func (m *mockRespository) Create(ctx context.Context, erroll *domain.Enrollment) error {
	return m.CreateMock(ctx, erroll)
}

func (m *mockRespository) GetAll(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
	return m.GetAllMock(ctx, filters, offset, limit)
}

func (m *mockRespository) Update(ctx context.Context, id string, status *string) error {
	return m.UpdateMock(ctx, id, status)
}
func (m *mockRespository) Count(ctx context.Context, filters enrollment.Filters) (int, error) {
	return m.CountMock(ctx, filters)
}
