package dbrepo

import (
	"errors"
	"time"

	"github.com/aiym182/booking/internal/models"
)

func (m *testDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into database.
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	// if the room id is 2 , then fail ; otherwise, pass

	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// inserts room restriction into database

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	if r.RoomID == 5 {
		return errors.New("some error")
	}
	return nil

}

//SearachAvailabilityByDates returns true if availability exists for roomID, and false if no availabiility exist
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	startDate := start.String()
	endDate := end.String()
	if startDate > endDate {
		return false, errors.New("some error")
	}
	return false, nil
}

// SearchAvailabiltyForAllRooms returns slice of available rooms , if any , for given range
func (m *testDBRepo) SearchAvailabiltyForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	startDate := start.String()
	endDate := end.String()
	if startDate > endDate {
		return nil, errors.New("some error")
	} else if startDate == endDate {
		return rooms, nil
	}

	rooms = append(rooms, models.Room{
		ID:        1,
		RoomName:  "Generals Quaters",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return rooms, nil

}

// GetRoomById gets a room by id from database

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {

	var u models.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {

	return nil
}
func (m *testDBRepo) Authenticate(email, password string) (int, string, error) {

	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil

}
