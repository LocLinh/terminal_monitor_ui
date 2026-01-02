package entities

import "time"

// import "gorm.io/gorm"

type Inbound struct {
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

func (m *Inbound) TableName() string {
	return "HasakiSystem"
}

type Store struct {
	StoreId         int64
	StoreName       string
	StoreLocationId int64
	StoreServiceId  int64
	StoreStockId    int64
	StoreStatus     int64
	LocationId      int64
	StoreConfig     int64
	Properties      string
	Zone            int64
}

func (m *Store) TableName() string {
	return "HasakiStore"
}

type NowLocation struct {
	PickupLocationId int64
	Name             string `json:"name"`
	InsideId         int64  `json:"inside_id"`
	LocationId       int64  `json:"location_id"`
	Zone             int64  `json:"zone"`
}

func (m *NowLocation) TableName() string {
	return "NowLocation"
}
