package postgres

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
)

func (p *postgres) TotalCost(ctx context.Context, userID *string, serviceName *string, from *time.Time, to *time.Time) (int, error) {
	builder := squirrel.Select("COALESCE(SUM(price), 0)").
		From(tableName).
		PlaceholderFormat(squirrel.Dollar)

	if userID != nil {
		builder = builder.Where(squirrel.Eq{userIDCol: *userID})
	}
	if serviceName != nil {
		builder = builder.Where(squirrel.Eq{serviceNameCol: *serviceName})
	}
	if from != nil {
		builder = builder.Where(squirrel.GtOrEq{startDateCol: *from})
	}
	if to != nil {
		builder = builder.Where(squirrel.LtOrEq{startDateCol: *to})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		p.l.Errorw("Error building query", "error", err)
		return 0, err
	}

	var totalCost int
	err = p.db.QueryRow(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		p.l.Errorw("Error querying total cost", "error", err)
		return 0, err
	}

	return totalCost, nil
}
