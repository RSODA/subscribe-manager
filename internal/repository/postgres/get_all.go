package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/RSODA/subscribe-manager/internal/domain"
)

func (p *postgres) GetAll(ctx context.Context) ([]*domain.Subscription, error) {
	res := []*domain.Subscription{}

	builder := squirrel.Select(idCol, serviceNameCol, userIDCol, priceCol, startDateCol, endDateCol).
		From(tableName).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		p.l.Errorw("Error building query", "error", err)
		return nil, err
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		p.l.Errorw("Error querying database", "error", err, "query", query, "args", args)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.Subscription
		err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.UserID,
			&s.Price,
			&s.StartDate,
			&s.EndDate,
		)
		if err != nil {
			p.l.Errorw("Error scanning row", "error", err)
			return nil, err
		}
		res = append(res, &s)
	}

	if err := rows.Err(); err != nil {
		p.l.Errorw("Error iterating over rows", "error", err)
		return nil, err
	}

	return res, nil
}
