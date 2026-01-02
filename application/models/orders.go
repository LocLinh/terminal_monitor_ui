package models

type OrderMessage struct {
	Payload Order `json:"payload"`
}

// Model of kafka message
type Order struct {
	// Id                           int64   `json:"id"`
	// Code                         string  `json:"code"`
	// SenderId                     int64   `json:"sender_id"`
	PickupLocationId   int64 `json:"pickup_location_id"`
	ReceiverLocationId int32 `json:"receiver_location_id"`
	// ReceiverName                 string  `json:"receiver_name"`
	// ReceiverPhone                string  `json:"receiver_phone"`
	// ReceiverAddress              string  `json:"receiver_address"`
	// CityId                       int64   `json:"city_id"`
	// DistrictId                   int64   `json:"district_id"`
	// WardId                       int64   `json:"ward_id"`
	Weight float64 `json:"weight"`
	Height int64   `json:"height"`
	// Cod                          int64   `json:"cod"`
	// PackageValue                 int32   `json:"package_value"`
	// VehicleType                  string  `json:"vehicle_type"`
	ShipperId int64  `json:"shipper_id,omitempty"`
	Category  string `json:"category,omitempty"`
	// CategoryType                 int64   `json:"category_type"`
	// PromotionId                  int64   `json:"promotion_id"`
	// PaymentParty                 int64   `json:"payment_party"`
	// Caution                      int64   `json:"caution"`
	// Note                         string  `json:"note"`
	Status    int64  `json:"status,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt                    string  `json:"updated_at"`
	// DeletedAt                    string  `json:"deleted_at"`
	// Longitude                    string  `json:"longitude"`
	// Latitude                     string  `json:"latitude"`
	// ReceivedAt                   int64   `json:"received_at"`
	// PartnerExpectedTime          string  `json:"partner_expected_time"`
	// PartnerExpectedTimeUpdatedAt string  `json:"partner_expected_time_updated_at"`
	// TransportCost                int64   `json:"transport_cost"`
	// OthersCost                   int64   `json:"others_cost"`
	// Discount                     int64   `json:"discount"`
	// TotalCost                    int64   `json:"total_cost"`
	Width  int64 `json:"width"`
	Length int64 `json:"length"`
	// PackageType                  int64   `json:"package_type"`
	PackageCode string `json:"package_code,omitempty"`
	ItCode      int64  `json:"it_code"`
	// ReturningTime                int64   `json:"returning_time"`
	// Ordering                     int64   `json:"ordering"`
	// Review                       string  `json:"review"`
	// PackageCod                   int64   `json:"package_cod"`
	// AdminRating                  int64   `json:"admin_rating"`
	// ClientRating                 int64   `json:"client_rating"`
	// RatingCreatedAt              string  `json:"rating_created_at"`
	// RatingAt                     string  `json:"rating_at"`
	Zone int64 `json:"zone,omitempty"`
	// LateNoDeduction              int64   `json:"late_no_deduction"`
	// AccumulativeCount            int64   `json:"accumulative_count"`
	// AccumulativeSenderCount      int64   `json:"accumulative_sender_count"`
	// PrintShiftOffAt              string  `json:"print_shift_off_at"`
	// ZoneLongitude                string  `json:"zone_longitude"`
	// ZoneLatitude                 string  `json:"zone_latitude"`
	// AssignedShipperBy            int64   `json:"assigned_shipper_by"`
	// WaitedAt                     string  `json:"waited_at"`
	// UpdateReceivedTime           int64   `json:"update_received_time"`
	// TransshipmentPointId         int64   `json:"transshipment_point_id"`
	// TransferAt                   string  `json:"transfer_at"`
	// FirstZone                    int64   `json:"first_zone"`
	// IsRevoke                     int64   `json:"is_revoke"`
	Deleted string `json:"__deleted,omitempty"`
}
