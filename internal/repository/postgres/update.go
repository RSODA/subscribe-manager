package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *postgres) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	res := &domain.Subscription{}

	builder := squirrel.Update(tableName).
		Where(squirrel.Eq{idCol: sub.ID}).
		Suffix("RETURNING " + idCol + ", " + serviceNameCol + ", " + userIDCol + ", " + priceCol + ", " + startDateCol + ", " + endDateCol).
		PlaceholderFormat(squirrel.Dollar)

	if sub.UserID != uuid.Nil {
		builder = builder.Set(userIDCol, sub.UserID)
	}
	if sub.ServiceName != "" {
		builder = builder.Set(serviceNameCol, sub.ServiceName)
	}
	if sub.Price > 0 {
		builder = builder.Set(priceCol, sub.Price)
	}
	if !sub.StartDate.IsZero() {
		builder = builder.Set(startDateCol, sub.StartDate)
	}
	if sub.EndDate != nil && !sub.EndDate.IsZero() {
		builder = builder.Set(endDateCol, sub.EndDate)
	}

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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrSubscriptionNotFound
		}
		p.l.Errorw("Error updating subscription", "error", err)
		return nil, err
	}

	return res, nil
}
