package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/RSODA/subscribe-manager/internal/domain"
)

func (p *postgres) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	res := &domain.Subscription{}

	builder := squirrel.Update(tableName).
		Set(serviceNameCol, sub.ServiceName).
		Set(userIDCol, sub.UserID).
		Set(priceCol, sub.Price).
		Set(startDateCol, sub.StartDate).
		Set(endDateCol, sub.EndDate).
		Where(squirrel.Eq{idCol: sub.ID}).
		Suffix("RETURNING " + idCol + ", " + serviceNameCol + ", " + userIDCol + ", " + priceCol + ", " + startDateCol + ", " + endDateCol).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		p.l.Errorw("Error building query", "error", err)
		return nil, err
	}

	err = p.db.QueryRow(ctx, query, args...).Scan(
		&res.ID,
		&res.ServiceName,
		&res.UserID,
		&res.Price,
		&res.StartDate,
		&res.EndDate,
	)
	if err != nil {
		p.l.Errorw("Error updating subscription", "error", err)
		return nil, err
	}

	return res, nil
}
