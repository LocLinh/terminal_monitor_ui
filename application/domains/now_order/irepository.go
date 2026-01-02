package now_order

import (
	"context"
	"terminal_monitor_ui/application/models"
)

type INowOrderRepository interface {
	CreateOrUpdateInbound(ctx context.Context, topic string, msg []byte) error
	CreateOrUpdateNowOrder(ctx context.Context, order *models.Order) error
}
