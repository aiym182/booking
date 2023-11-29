package dbrepo

import (
	"errors"
	"time"

	"github.com/aiym182/booking/internal/models"
)

// setting bool for testing purpose
func (m *testDBRepo) SetBool(choose bool) bool {

	return choose
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

	if email == "me@here.ca" {
		return 1, "", nil
	}
	return 0, "", errors.New("some error")
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation
	if m.App.TestError {
		return nil, errors.New("some error")
	}

	return reservations, nil

}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	reservations := []models.Reservation{{ID: "1"}}
	if m.App.TestError {

		return nil, errors.New("some error")
	}
	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	var reservation models.Reservation

	if id > 100 {
		return reservation, errors.New("some error")
	}

	return reservation, nil
}

func (m *testDBRepo) UpdateReservation(u models.Reservation) error {

	if u.Email == "" {
		return errors.New("some error")
	}
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	if id > 1000 {
		return errors.New("some error")
	}

	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room

	if m.App.TestError {
		return rooms, errors.New("some errors")
	}

	rooms = append(rooms, models.Room{
		ID:        1,
		RoomName:  "Generals Quaters",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()},
	)
	// rooms = append(rooms, models.Room{
	// 	ID:        1000,
	// 	RoomName:  "Generals Quaters",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// })

	return rooms, nil
}

func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var Restrictions []models.RoomRestriction

	startDate := start.Month().String()

	if startDate == "February" {

		return nil, errors.New("some error")
	}

	if roomID > 1 {

		Restrictions = append(Restrictions, models.RoomRestriction{
			ID:            2,
			StartDate:     start,
			EndDate:       end,
			RoomID:        roomID,
			ReservationID: 0,
			RestrictionID: 1,
			Room:          models.Room{},
			Reservation:   models.Reservation{},
			Restriction:   models.Restriction{},
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
		return Restrictions, nil
	}

	Restrictions = append(Restrictions, models.RoomRestriction{
		ID:            1,
		StartDate:     start,
		EndDate:       end,
		RoomID:        roomID,
		ReservationID: 1,
		RestrictionID: 1,
		Room:          models.Room{},
		Reservation:   models.Reservation{},
		Restriction:   models.Restriction{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})

	return Restrictions, nil
}

func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *testDBRepo) DeleteBlockByID(id int) error {

	return nil
}
