package database

import (
	"database/sql"
	"time"

	"nicksrepo.com/nick/pkg/models"
)

type DbHandler interface {
	AvailabilitySearch(start_date, end_date string) ([]*models.SearchAvailabilityModel, error)
	InsertReservation(*models.Reservation) (int64, error)
	DeleteReservation(id int64) (sql.Result, error)
	InsertRoomRestriction(*models.RoomRestriction) (sql.Result, error)
	LastAvailabilitySearch(id string, start_date, end_date time.Time) (bool, error)
}
