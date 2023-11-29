package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/aiym182/booking/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// setting bool for testing purposes
func (m *PostgresDbRepo) SetBool(choose bool) bool {

	return choose
}

// InsertReservation inserts a reservation into database and returns the id of the row
func (m *PostgresDbRepo) InsertReservation(res models.Reservation) (int, error) {

	// if statement has something wrong, after 3 min it will automatically close.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date,
		room_id,created_at,updated_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

// inserts room restriction into database

func (m *PostgresDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	// if statement is something wrong, after 3 min it will automatically close.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id , reservation_id,
			created_at, updated_at, restriction_id)
			values ($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
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

//SearachAvailabilityByDates returns true if availability exists for roomID, and false if no availabiility exist
func (m *PostgresDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var numRows int

	query := `select count(id) from room_restrictions 
			where room_id = $1 and $2 < end_date and $3 > start_date`

	err := m.DB.QueryRowContext(ctx, query, roomID, start, end).Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabiltyForAllRooms returns slice of available rooms , if any , for given range
func (m *PostgresDbRepo) SearchAvailabiltyForAllRooms(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var rooms []models.Room
	query := `select
					r.id, r.room_name
		      from
					rooms r
			  where 
					r.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
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

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil

}

// GetRoomById gets a room by id from database

func (m *PostgresDbRepo) GetRoomByID(id int) (models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id =$1`

	row := m.DB.QueryRowContext(ctx, query, id)

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

//GetUserByID returns a user by id

func (m *PostgresDbRepo) GetUserByID(ID int) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `select id, first_name, last_name, email,password,access_level,created_at,updated_at
			  from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, ID)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *PostgresDbRepo) UpdateUser(u models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level =$4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

//Authenticate authenticates the user

func (m *PostgresDbRepo) Authenticate(email, password string) (int, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var id int
	var hashedPassword string

	query := `select id, password from users where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil

}

// All reservations returns a slice of all reservations from database
func (m *PostgresDbRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at,
	rm.id, rm.room_name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations from database
func (m *PostgresDbRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
	rm.id, rm.room_name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	where processed = 0
	order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,

			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

//GetReservationByID returns one reservation by ID
func (m *PostgresDbRepo) GetReservationByID(id int) (models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var res models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email ,r.phone,
		r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r 
		left join rooms rm on (r.room_id= rm.id)
		where r.id = $1
		`
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

//UpdateReservations updates reservation in the database
func (m *PostgresDbRepo) UpdateReservation(u models.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update reservations set first_name = $1, last_name = $2, email = $3, phone =$4, updated_at = $5 where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

// DeleteReservation deletes one reservation by id
func (m *PostgresDbRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `delete from reservations where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

//UpdateProcessedForReservation updates processed for reservation by id
func (m *PostgresDbRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update reservations set processed = $1 where id =$2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)

	if err != nil {
		return err
	}

	return nil
}

// all rooms returns all []modles.Room
func (m *PostgresDbRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var rooms []models.Room

	query := `select id, room_name, created_at, updated_at from rooms order by room_name`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return rooms, err
	}

	defer rows.Close()

	for rows.Next() {
		var rm models.Room

		err := rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.CreatedAt,
			&rm.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, rm)
	}

	return rooms, nil

}

// GetRestrictionsForRoomByDate returns restrictions for a room by date range

func (m *PostgresDbRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
			  from room_restrictions where $1 <end_date and $2 >=start_date and room_id =$3
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)

		if err != nil {
			return nil, err
		}

		restrictions = append(restrictions, r)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *PostgresDbRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at)
			  values ($1,$2,$3,$4,$5,$6)	
	`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())

	if err != nil {
		m.App.ErrorLog.Println("err")
		return err
	}

	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *PostgresDbRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `delete from room_restrictions where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		m.App.ErrorLog.Println("err")
		return err
	}

	return nil
}
