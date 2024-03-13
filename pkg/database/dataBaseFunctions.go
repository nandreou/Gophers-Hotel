package database

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"nicksrepo.com/nick/pkg/models"
)

func (db *DB) AvailabilitySearch(start_date, end_date string) ([]*models.SearchAvailabilityModel, error) {
	query := `SELECT rooms.id, rooms.room_name FROM rooms WHERE rooms.id not in (SELECT room_restrictions.room_id FROM room_restrictions WHERE start_date < $2 AND end_date > $1)`

	var res []*models.SearchAvailabilityModel

	rows, err := db.SQL.Query(query, start_date, end_date)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tmp models.SearchAvailabilityModel
		if err := rows.Scan(
			&tmp.Room_id,
			&tmp.Room_name,
		); err != nil {
			return nil, err
		} else {
			res = append(res, &tmp)
		}
	}
	return res, nil
}

func (db *DB) InsertReservation(res *models.Reservation) (int64, error) {
	query := `INSERT INTO reservations (room_id, first_name, last_name, email, phone, start_date, end_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var psqlId int64

	err := db.SQL.QueryRow(query,
		res.RoomId,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		time.Now(),
		time.Now(),
	).Scan(&psqlId)

	if err != nil {
		return 0, err
	}

	return psqlId, nil
}

func (db *DB) DeleteReservation(id int64) (sql.Result, error) {
	query := `DELETE FROM reservations WHERE id = $1`
	return db.SQL.Exec(query, id)
}

func (db *DB) InsertRoomRestriction(rstr *models.RoomRestriction) (sql.Result, error) {
	query := `INSERT INTO room_restrictions (room_id, reservation_id, restriction_id, start_date, end_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	return db.SQL.Exec(query,
		rstr.RoomId,
		rstr.ReservationId,
		rstr.RestrictionId,
		rstr.StartDate,
		rstr.EndDate,
		time.Now(),
		time.Now(),
	)
}

func (db *DB) LastAvailabilitySearch(id string, start_date, end_date time.Time) (bool, error) {
	query := `SELECT room_restrictions.room_id FROM room_restrictions WHERE start_date < $2 AND end_date > $1 AND room_id = $3`

	var (
		tmp int
	)

	db.SQL.QueryRow(query, start_date, end_date, id).Scan(&tmp)

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return false, err
	}

	if idInt == tmp {
		log.Println(idInt)
		return false, nil
	}

	return true, nil
}
