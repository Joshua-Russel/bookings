package repository

import (
	"github.com/Joshua-Russel/bookings/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservations(reservation models.Reservation) (int, error)
	InsertRoomRestrictions(room models.RoomRestriction) error
	SearchAvailabilityByDatesAndRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(roomId int) (models.Room, error)
}
