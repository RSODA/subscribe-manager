package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/RSODA/subscribe-manager/internal/domain"
)

func (p *postgres) Create(ctx context.Context, s *domain.Subscription) (*domain.Subscription, error) {
	builder := squirrel.Insert(tableName).
		Columns(serviceNameCol, userIDCol, priceCol, startDateCol, endDateCol).
		Values(s.ServiceName, s.UserID, s.Price, s.StartDate, s.EndDate).
		Suffix("RETURNING " + idCol + ", " + serviceNameCol + ", " + userIDCol + ", " + priceCol + ", " + startDateCol + ", " + endDateCol).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		p.l.Errorw("Error building query", "error", err)
		return nil, err
	}

	result := &domain.Subscription{}
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&result.ID,
		&result.ServiceName,
		&result.UserID,
		&result.Price,
		&result.StartDate,
		&result.EndDate,
	)
	if err != nil {
		p.l.Errorw("Error creating subscription", "error", err)
		return nil, err
	}

	return result, nil
}
