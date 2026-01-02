package entity

import "time"

type NowOrder struct {
	PackageCode        string    `gorm:"primaryKey;column:maKienHang"`
	ZoneCode           int64     `gorm:"column:maZone"`
	ZoneCodeOriginal   int64     `gorm:"column:zone_code_original"`
	Status             int64     `gorm:"column:trangThai"`
	PreviousStatus     int64     `gorm:"column:trangThaiTruoc"`
	Category           string    `gorm:"column:category"`
	ShipperId          int64     `gorm:"column:shipper_id"`
	PickupLocationId   int64     `gorm:"column:pickup_location_id"`
	ReceiverLocationId int32     `gorm:"column:receiver_location_id"`
	Width              int64     `gorm:"column:width"`
	Length             int64     `gorm:"column:length"`
	Height             int64     `gorm:"column:height"`
	Weight             float64   `gorm:"column:weight"`
	ItCode             int64     `gorm:"column:it_code;default:(-)"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (m *NowOrder) TableName() string {
	return "HasakiSystem"
}
