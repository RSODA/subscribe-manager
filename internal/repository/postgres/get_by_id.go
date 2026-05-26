package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (p *postgres) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	res := &domain.Subscription{}

	builder := squirrel.Select(idCol, serviceNameCol, userIDCol, priceCol, startDateCol, endDateCol).
		From(tableName).
		Where(squirrel.Eq{idCol: id}).
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrSubscriptionNotFound
		}
		p.l.Errorw("Error getting subscription by ID", "error", err)
		return nil, err
	}

	return res, nil
}
