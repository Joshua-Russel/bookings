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
	GetUserByID(userId int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, password string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(res models.Reservation) error
	UpdateProcessedForReservation(id, processed int) error
	DeleteReservation(id int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)

	InsertBlockForRoom(id int, startDate time.Time) error
	DeleteBlockByID(id int) error
}
