package now_order

import "context"

type INowOrderUsecase interface {
	ConsumeMessage(ctx context.Context, vals []byte) error
}
