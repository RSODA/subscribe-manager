package postgres

import (
	"github.com/RSODA/subscribe-manager/internal/db"
	"github.com/RSODA/subscribe-manager/internal/repository"
	"go.uber.org/zap"
)

type postgres struct {
	l  *zap.SugaredLogger
	db db.DB
}

const (
	tableName = "subscriptions"

	idCol          = "id"
	userIDCol      = "user_id"
	priceCol       = "price"
	serviceNameCol = "service_name"
	startDateCol   = "start_date"
	endDateCol     = "end_date"
)

func NewPostgresRepository(db db.DB, l *zap.SugaredLogger) repository.SubscriptionRepository {
	return &postgres{db: db, l: l}
}
