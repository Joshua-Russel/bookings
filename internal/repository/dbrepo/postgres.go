package dbrepo

import (
	"context"
	"github.com/Joshua-Russel/bookings/internal/models"
	"time"
)

func (psql *postgresDBRepo) AllUsers() bool {
	return true
}
func (psql *postgresDBRepo) InsertReservations(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	var resId int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
                      end_date, room_id, created_at, updated_at) values 
					($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := psql.DB.QueryRowContext(ctx, stmt, res.FirstName,
		res.LastName, res.Email,
		res.Phone,
		res.StartDate, res.EndDate,
		res.RoomID, time.Now(), time.Now()).Scan(&resId)
	if err != nil {
		return 0, err
	}
	return resId, nil
}

func (psql *postgresDBRepo) InsertRoomRestrictions(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
                               created_at, updated_at, restriction_id) values 
							   ($1,$2,$3,$4,$5,$6,$7);`
	_, err := psql.DB.ExecContext(ctx, stmt, r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (psql *postgresDBRepo) SearchAvailabilityByDatesAndRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select count(id) 
				from room_restrictions 
				where room_id = $1 and 
				      $2 < end_date and $3 > start_date;`
	row := psql.DB.QueryRowContext(ctx, stmt, roomId, start, end)
	var id_count int
	err := row.Scan(&id_count)
	if err != nil {
		return false, err
	}
	if id_count == 0 {
		return true, nil
	}
	return false, nil

}
func (psql *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `
select r.id, r.room_name
from rooms r 
where r.id not in (select room_id from room_restrictions rr where $1 < rr.end_date  and $2 > rr.start_date);`
	rows, err := psql.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}
func (psql *postgresDBRepo) GetRoomById(roomId int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := ` select id, room_name,created_at, updated_at from rooms where id = $1;`
	row := psql.DB.QueryRowContext(ctx, query, roomId)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil

}
