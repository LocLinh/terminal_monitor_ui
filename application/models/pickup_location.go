package models

type PickupLocationMessage struct {
	Payload PickupLocation `json:"payload"`
}

type PickupLocation struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	InsideId   int64  `json:"inside_id"`
	LocationId int64  `json:"location_id"`
	Zone       int64  `json:"zone"`
}
