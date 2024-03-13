package models

import "time"

type User struct {
	FirstName   string
	LastName    string
	Email       string
	PassWord    string
	AccessLevel int
}

type Reservation struct {
	RoomId    int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
}

type RoomRestriction struct {
	Id            int
	RoomId        int
	ReservationId int
	RestrictionId int
	StartDate     time.Time
	EndDate       time.Time
	Created_at    time.Time
	Updated_at    time.Time
}

type SearchAvailabilityModel struct {
	Room_id   int
	Room_name string
}
