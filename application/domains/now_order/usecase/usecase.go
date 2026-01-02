package usecase

import (
	"context"
	"encoding/json"
	no "terminal_monitor_ui/application/domains/now_order"
	"terminal_monitor_ui/application/models"
	"terminal_monitor_ui/constant"
	"terminal_monitor_ui/logger"

	"go.uber.org/zap"
)

type NowOrderUsecase struct {
	NoRepo no.INowOrderRepository
}

func InitUsecase(
	noRepo no.INowOrderRepository,
) no.INowOrderUsecase {
	return &NowOrderUsecase{
		NoRepo: noRepo,
	}
}

func (s *NowOrderUsecase) ConsumeMessage(ctx context.Context, vals []byte) error {
	msgNo := &models.OrderMessage{}
	if string(vals) == "" {
		zap.S().Info("Drop message empty")
		return nil
	}
	if err := json.Unmarshal(vals, msgNo); err != nil {
		_ = logger.WriteFile("./logger/error.txt", err.Error()+"\n")
		zap.S().Errorf("Error mashal shipping order :%s \n", string(vals))
		return err
	}

	if msgNo.Payload.CreatedAt < "2024-03-14T00:00:00Z" {
		zap.S().Info("Drop msg < 2024-03-14")
		return nil
	}
	if msgNo.Payload.PackageCode == "" {
		zap.S().Info("Drop msg PackageCode empty")
		return nil
	}

	// not sync any data other than than that has category "48hInter"
	if msgNo.Payload.Category != constant.INBOUND_CATEGORY_48H_INTER {
		zap.S().Info("Drop msg != 48hInter")
		return nil
	}

	err := s.NoRepo.CreateOrUpdateNowOrder(ctx, &msgNo.Payload)
	if err != nil {
		zap.S().Errorf("Error update inbound item:%s \n", string(vals))
		return err
	}

	zap.S().Infof("success save order: %v", msgNo.Payload.PackageCode)

	return nil
}
