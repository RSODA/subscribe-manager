package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
)

func (p *postgres) Delete(ctx context.Context, id string) error {
	builder := squirrel.Delete(tableName).
		Where(squirrel.Eq{idCol: id}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		p.l.Errorw("Error building delete query", "error", err)
		return err
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		p.l.Errorw("Error executing delete query", "error", err)
		return err
	}

	return nil
}
