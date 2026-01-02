package repository

import (
	"context"
	"encoding/json"
	"fmt"
	no "terminal_monitor_ui/application/domains/now_order"
	"terminal_monitor_ui/application/domains/now_order/entity"
	"terminal_monitor_ui/application/models"
	"terminal_monitor_ui/constant"
	"terminal_monitor_ui/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NowOrderRepository struct {
	Db *gorm.DB
}

func InitNowOrderRepository(db *gorm.DB) no.INowOrderRepository {
	return &NowOrderRepository{
		Db: db,
	}
}

func (i *NowOrderRepository) CreateOrUpdateInbound(ctx context.Context, topic string, msg []byte) error {
	// zap.S().Infof("msg=%s \n", string(msg))
	order := &models.OrderMessage{}
	if string(msg) == "" {
		zap.S().Info("Drop msg empty")
		return nil
	}
	if err := json.Unmarshal(msg, order); err != nil {
		_ = logger.WriteFile("./logger/error.txt", err.Error()+"\n")
		zap.S().Errorf("Error mashal inbound item:%s \n", string(msg))
		return err
	}

	if order.Payload.CreatedAt < "2023-05-01T00:00:00Z" {
		zap.S().Info("Drop msg < 2023-05-01T00:00:00Z")
		return nil
	}
	if order.Payload.PackageCode == "" {
		zap.S().Info("Drop msg PackageCode empty")
		return nil
	}

	// not sync any data other than than that has category "48hInter"
	if order.Payload.Category != constant.INBOUND_CATEGORY_48H_INTER {
		zap.S().Info("Drop msg != 48hInter")
		return nil
	}

	inboundItem := entity.NowOrder{
		PackageCode:        order.Payload.PackageCode,
		ZoneCodeOriginal:   order.Payload.Zone,
		Status:             order.Payload.Status,
		Category:           order.Payload.Category,
		ShipperId:          order.Payload.ShipperId,
		PickupLocationId:   order.Payload.PickupLocationId,
		ReceiverLocationId: order.Payload.ReceiverLocationId,
		Width:              order.Payload.Width,
		Length:             order.Payload.Length,
		Height:             order.Payload.Height,
		Weight:             order.Payload.Weight,
		ItCode:             order.Payload.ItCode,
	}

	if order.Payload.Zone < 1 || order.Payload.Zone > 9 {
		inboundItem.ZoneCode = 0
	} else {
		inboundItem.ZoneCode = order.Payload.Zone
	}

	existingInbound := entity.NowOrder{}
	result := i.Db.Where(entity.NowOrder{PackageCode: inboundItem.PackageCode}).
		Attrs(inboundItem).
		FirstOrCreate(&existingInbound)
	if result.Error != nil {
		_ = logger.WriteFile("./logger/error.txt", fmt.Sprintf("[%s] Error: %v\n\n", inboundItem.PackageCode, result.Error.Error()))
		zap.S().Errorf("Error find or create inbound item:%v \n", result.Error)
		return result.Error
	}

	// if existing then update existing
	if result.RowsAffected == 0 {
		existingInbound = inboundItem
		if existingInbound.PreviousStatus != inboundItem.Status {
			existingInbound.PreviousStatus = existingInbound.Status
		}
		result := i.Db.Save(&existingInbound)
		if result.Error != nil {
			_ = logger.WriteFile("./logger/error.txt", fmt.Sprintf("[%s] Error: %v\n\n", inboundItem.PackageCode, result.Error.Error()))
			zap.S().Errorf("Error update inbound item:%v \n", result.Error)
			return result.Error
		}
	}

	zap.S().Infof("Success saving %s inbound item:%v \n", order.Payload.CreatedAt, inboundItem)
	return nil
}

func (i *NowOrderRepository) CreateOrUpdateNowOrder(ctx context.Context, order *models.Order) error {
	inboundItem := entity.NowOrder{
		PackageCode:        order.PackageCode,
		ZoneCodeOriginal:   order.Zone,
		Status:             order.Status,
		Category:           order.Category,
		ShipperId:          order.ShipperId,
		PickupLocationId:   order.PickupLocationId,
		ReceiverLocationId: order.ReceiverLocationId,
		Width:              order.Width,
		Length:             order.Length,
		Height:             order.Height,
		Weight:             order.Weight,
		ItCode:             order.ItCode,
	}

	if order.Zone < 1 || order.Zone > 9 {
		inboundItem.ZoneCode = 0
	} else {
		inboundItem.ZoneCode = order.Zone
	}

	existingInbound := entity.NowOrder{}
	result := i.Db.Where(entity.NowOrder{PackageCode: inboundItem.PackageCode}).
		Attrs(inboundItem).
		FirstOrCreate(&existingInbound)
	if result.Error != nil {
		zap.S().Errorf("Error find or create inbound item:%v \n", result.Error)
		return result.Error
	}

	// if existing then update existing
	if result.RowsAffected == 0 {
		// existingInbound = inboundItem
		if existingInbound.Status != inboundItem.Status {
			inboundItem.PreviousStatus = existingInbound.Status
		}
		result := i.Db.Save(&inboundItem)
		if result.Error != nil {
			zap.S().Errorf("Error update inbound item:%v \n", result.Error)
			return result.Error
		}
	}

	return nil
}
