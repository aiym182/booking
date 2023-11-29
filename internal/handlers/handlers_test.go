package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aiym182/booking/internal/models"
)

// type PostData struct {
// 	key   string
// 	value string
// }

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quaters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-existent", "/green/eggs/and/ham", "GET", http.StatusNotFound},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"all res", "/admin/reservations-all", "GET", http.StatusOK},

	// {"post-search-avail", "/search-availability", "POST", []PostData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []PostData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"make-reservation", "/make-reservation", "POST", []PostData{
	// 	{key: "first_name", value: "john"},
	// 	{key: "last_name", value: "smith"},
	// 	{key: "email", value: "aiym182@gmail.com"},
	// 	{key: "phone", value: "432242323"},
	// }, http.StatusOK},
}

func TestHand(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)

	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Error(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {

	reservation := models.Reservation{

		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quaters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// rr stands for response recorder
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code : got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//test case where reservation is not in session (reset everything)

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code : got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code : got %d, wanted %d", rr.Code, http.StatusOK)

	}
}

func TestRepository_PostReservation(t *testing.T) {

	reservation := models.Reservation{

		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quaters",
		},
	}

	// reqBody := "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "invalid")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code : got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test where handler can't get reservation data from session
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing data from session : got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test where handler can't parse form due to empty body

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test where firstName is not valid

	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "J")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returned wrong response code for invalid firstName: got %d, wanted %d", rr.Code, http.StatusOK)

	}

	// test for failure to insert reservation into database

	reservation = models.Reservation{

		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Major's Suite",
		},
	}
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid firstName: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)

	}

	// test for failure to insert room restriction data into database

	reservation = models.Reservation{

		RoomID: 5,
		Room: models.Room{
			ID:       5,
			RoomName: "Some room",
		},
	}
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid firstName: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)

	}

}

func TestRepository_PostAvailability(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test where handler has empty body
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for no body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)

	}

	// test wehre handler has invalid startDate

	postedData = url.Values{}
	postedData.Add("start", "invalid")
	postedData.Add("end", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code invalid start date : got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)

	}

	// test where handler has invalid endDate

	postedData = url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "invalid")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)

	}

	// test where searchAvailabilityForAllRooms has error

	postedData = url.Values{}
	postedData.Add("start", "2050-01-02")
	postedData.Add("end", "2050-01-01")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for seachAvailabilityForAllRooms: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test where searchAvailabilityForAllRooms has no room

	postedData = url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-01")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code when there is 0 room alvaiable: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var j jsonResponse

	err := json.Unmarshal(rr.Body.Bytes(), &j)

	if err != nil {
		t.Error("unable to parse json")
	}

	// test where handler has no body
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &j)

	if err != nil {
		t.Error("unable to parse json")
	}

	// test where handler cannot conncect to database
	postedData = url.Values{}
	postedData.Add("start", "2050-01-02")
	postedData.Add("end", "2050-01-01")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &j)

	if err != nil {
		t.Error("unable to parse json")
	}

}

func TestRepository_ReservationSummary(t *testing.T) {

	reservation := models.Reservation{}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test where there is no session data

	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {

	reservation := models.Reservation{}
	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case when handler can't parse from url parameter

	req, _ = http.NewRequest("GET", "/choose-room/ho", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/ho"
	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code for wrong url parameter: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case where handler can't get session

	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code for no session: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}
func TestRepository_BookRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// testcase where database failed

	req, _ = http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code for failed database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"e",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestRepository_TestLogin(t *testing.T) {

	// test case when handler can't parse form (no body)

	req, _ := http.NewRequest("POST", "/user/login", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostShowLogin)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("failed PostShowLogin: expected code %d, but got %d", http.StatusOK, rr.Code)
	}

	// range through all tests

	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request

		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)

		req = req.WithContext(ctx)

		// set the header

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler

		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			//read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}

	}
}

func TestRepository_AdminNewReservation(t *testing.T) {

	// test when AllNewReservation function from database repo has an error
	req, _ := http.NewRequest("GET", "/admin/reservations-new", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	// setting testError true

	Repo.App.TestError = true

	handler := http.HandlerFunc(Repo.AdminNewReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("failed %s: expected code %d, but got %d", "invalid-database", http.StatusOK, rr.Code)
	}
}
func TestRepository_AdminAllReservations(t *testing.T) {

	// test whenn AllReservation function from database repo has an error
	req, _ := http.NewRequest("GET", "/admin/reservations-all", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	// setting testError true
	Repo.App.TestError = true
	handler := http.HandlerFunc(Repo.AdminAllReservations)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("failed %s: expected code %d, but got %d", "invalid-database", http.StatusOK, rr.Code)
	}
}

var AdminPostShowReservation = []struct {
	name                 string
	url                  string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{name: "valid-data-form-new",
		url: "/admin/reservations/new/1",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"444-444-4444"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-new",
		expectedHTML:         "",
	},
	{
		name: "valid-data-form-all",
		url:  "/admin/reservations/all/1",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"444-444-4444"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-all",
		expectedHTML:         "",
	},
	{name: "valid-data-form-cal",
		url: "/admin/reservations/cal/1",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"444-444-4444"},
			"year":       {"2022"},
			"month":      {"01"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservation-calendar?y=2022&m=01",
		expectedHTML:         "",
	},

	{name: "invalid-body",
		url:                  "/admin/reservations/new/1/show",
		expectedResponseCode: http.StatusOK,
		expectedHTML:         "",
		expectedLocation:     "",
	},
	{name: "invalid-url-length",
		url: "/admin/reservations/new/",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"444-444-4444"},
			"year":       {"2022"},
			"month":      {"01"},
		},
		expectedResponseCode: http.StatusOK,
		expectedLocation:     "/admin/reservations-new",
		expectedHTML:         "",
	},
	{name: "invalid-id",
		url: "/admin/reservations/new/1001",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"444-444-4444"},
			"year":       {"2022"},
			"month":      {"01"},
		}, expectedResponseCode: http.StatusOK,
		expectedLocation: "/admin/reservations-new",
		expectedHTML:     ""},
	{name: "invalid-reservation",
		url: "/admin/reservations/new/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {""},
			"phone":      {"444-444-4444"},
		},
		expectedResponseCode: http.StatusOK,
		expectedLocation:     "/admin/reservations-new",
		expectedHTML:         ""},
}

func TestRepository_AdminPostShowReservation(t *testing.T) {

	for _, test := range AdminPostShowReservation {

		// test to see handler to have right response
		var req *http.Request

		if test.postedData != nil {
			req, _ = http.NewRequest("POST", test.url, strings.NewReader(test.postedData.Encode()))

		} else {
			req, _ = http.NewRequest("POST", test.url, nil)

		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.RequestURI = test.url
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.AdminPostShowReservation)
		handler.ServeHTTP(rr, req)

		// checking to see handler gives expected responseCode
		if rr.Code != test.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", test.name, test.expectedResponseCode, rr.Code)
		}
		// checking to see handler moves to expected location

		actualLoc, _ := rr.Result().Location()
		if actualLoc != nil && actualLoc.String() != test.expectedLocation {
			t.Errorf("failed %s: expected location %s, but got %s", test.name, test.expectedLocation, actualLoc.String())
		}

		if test.expectedHTML != "" {
			// read response body into string

			html := rr.Body.String()
			if !strings.Contains(html, test.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", test.name, test.expectedHTML)

			}
		}

	}

}

var AdminShowReservation = []struct {
	name                 string
	url                  string
	expectedResponseCode int
}{
	{
		"valid-new",
		"/admin/reservations/new/1/show",
		http.StatusOK,
	},
	{
		"valid-all",
		"/admin/reservations/all/1/show",
		http.StatusOK,
	},
	{
		"valid-cal",
		"admin/reservations/cal/1/show?y=2022&m=1",
		http.StatusOK,
	},
	{
		"invalid-url",
		"/admin/reservations/new/ho/show",
		http.StatusOK,
	},
	{
		"invalid-id",
		"/admin/reservations/new/1000/show",
		http.StatusOK,
	},
}

func TestRepository_AdminShowReservation(t *testing.T) {

	for _, test := range AdminShowReservation {
		req, _ := http.NewRequest("GET", test.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		req.RequestURI = test.url
		handler := http.HandlerFunc(Repo.AdminShowReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", test.name, test.expectedResponseCode, rr.Code)
		}
	}
}

var AdminReservationCalendar = []struct {
	name                 string
	url                  string
	testError            bool
	expectedResponseCode int
}{
	{"no-month-and-year",
		"/admin/reservation-calendar",
		false,
		http.StatusOK,
	},
	{"month-and-year-specified",
		"/admin/reservation-calendar/?y=2022&m=01",
		false,
		http.StatusOK,
	},
	{
		"database-allrooms-error",
		"/admin/reservation-calendar",
		true,
		http.StatusOK,
	},
	{
		"database-GetRestrictionsForRoomByDate-error",
		"/admin/reservation-calendar/?y=2022&m=02",
		false,
		http.StatusOK,
	},
}

func TestRepository_AdminReservationCalendar(t *testing.T) {

	for _, test := range AdminReservationCalendar {

		req, _ := http.NewRequest("GET", test.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// rr stands for response recorder
		rr := httptest.NewRecorder()

		Repo.App.TestError = test.testError
		handler := http.HandlerFunc(Repo.AdminReservationCalendar)

		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", test.name, test.expectedResponseCode, rr.Code)

		}

	}
}

var AdminProcessReservation = []struct {
	name               string
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		"Valid-process-reservation-new",
		"/admin/process-reservation/new/1/do",
		http.StatusSeeOther,
		"/admin/reservations-new",
	},
	{
		"Valid-process-reservation-all",
		"/admin/process-reservation/all/1/do",
		http.StatusSeeOther,
		"/admin/reservations-all"},
	{
		"Valid-process-reservation-cal",
		"/admin/process-reservation/cal/1/do?y=2022&m=1",
		http.StatusSeeOther,
		"/admin/reservation-calendar?y=2022&m=1",
	},
	{
		"Invalid-database",
		"/admin/process-reservation/new/10000/do",
		http.StatusOK,
		"",
	},
}

func TestRepository_AdminProcessReservation(t *testing.T) {

	for _, test := range AdminProcessReservation {
		req, _ := http.NewRequest("GET", test.url, nil)

		ctx := getCtx(req)

		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		Repo.App.TestError = false

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)

		log.Println(req.URL.Query().Get("src"))
		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", test.name, http.StatusSeeOther, rr.Code)
		}

		// actualLoc, _ := rr.Result().Location()

		// if actualLoc.String() != test.expectedLocation {
		// 	t.Errorf("failed %s: expected location %s, but got %s", test.name, test.expectedLocation, actualLoc.String())

		// }
	}

}

var AdminDeleteReservation = []struct {
	name               string
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		"Valid-delete-reservation-new",
		"/admin/delete-reservation/new/1/do",
		http.StatusSeeOther,
		"/admin/reservations-new",
	},
	{
		"Valid-delete-reservation-all",
		"/admin/delete-reservation/all/1/do",
		http.StatusSeeOther,
		"/admin/reservations-all"},
	{
		"Valid-delete-reservation-cal",
		"/admin/delete-reservation/cal/1/do?y=2022&m=1",
		http.StatusSeeOther,
		"/admin/reservation-calendar?y=2022&m=1",
	},
}

func TestRepository_AdminDeleteReservation(t *testing.T) {

	for _, test := range AdminDeleteReservation {
		req, _ := http.NewRequest("GET", test.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		Repo.App.TestError = false

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", test.name, http.StatusSeeOther, rr.Code)
		}
	}
}

// var AdminPostReservationCalendar = []struct {
// 	name                 string
// 	url                  string
// 	postedData           url.Values
// 	expectedResponseCode int
// 	expectedLocation     string
// 	expectedHTML         string
// }{
// 	{
// 		name: "Valid-Admin-Post-Reservation-Calendar",
// 		url:  "/admin/reservation-calendar",
// 		postedData: url.Values{
// 			"year":  {time.Now().Format("2006")},
// 			"month": {time.Now().Format("01")},
// 			fmt.Sprintf("add_block_1_%s", time.Now().AddDate(0, 0, 2).Format("2006-01-2")): {"1"},
// 		},
// 		expectedResponseCode: http.StatusSeeOther,
// 		expectedLocation:     "/admin/reservation-calendar?y=2022&m=1",
// 		expectedHTML:         "",
// 	},
// 	{
// 		name:                 "invalid-body",
// 		url:                  "/admin/reservation-calendar",
// 		expectedResponseCode: http.StatusOK,
// 		expectedLocation:     "",
// 		expectedHTML:         "",
// 	},
// }

// func TestRepository_AdminPostReservationCalendar(t *testing.T) {

// 	for _, test := range AdminPostReservationCalendar {

// 		// test to see handler to have right response
// 		var req *http.Request

// 		if test.postedData != nil {
// 			req, _ = http.NewRequest("POST", test.url, strings.NewReader(test.postedData.Encode()))
// 		} else {
// 			req, _ = http.NewRequest("POST", test.url, nil)
// 		}
// 		ctx := getCtx(req)
// 		now := time.Now()
// 		blockMap := make(map[string]int)
// 		reservationMap := make(map[string]int)

// 		currentYear, currentMonth, _ := now.Date()
// 		currentLocation := now.Location()

// 		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
// 		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

// 		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
// 			reservationMap[d.Format("2006-01-2")] = 0
// 			blockMap[d.Format("2006-01-2")] = 0
// 		}

// 		// if test.blocks > 0 {
// 		// 	bm[firstOfMonth.Format("2006-01-2")] = e.blocks
// 		// }

// 		// if e.reservations > 0 {
// 		// 	rm[lastOfMonth.Format("2006-01-2")] = e.reservations
// 		// }

// 		session.Put(ctx, "block_map_1", blockMap)
// 		session.Put(ctx, "reservation_map_1", reservationMap)

// 		req = req.WithContext(ctx)
// 		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 		rr := httptest.NewRecorder()
// 		Repo.App.TestError = false

// 		handler := http.HandlerFunc(Repo.AdminPostReservationCalendar)

// 		handler.ServeHTTP(rr, req)

// 		if rr.Code != test.expectedResponseCode {
// 			t.Errorf("failed %s: expected code %d, but got %d", test.name, http.StatusSeeOther, rr.Code)
// 		}

// 	}
// }

var adminPostReservationCalendarTests = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
	blocks               int
	reservations         int
}{
	{
		name: "cal",
		postedData: url.Values{
			"year":  {time.Now().Format("2006")},
			"month": {time.Now().Format("01")},
			fmt.Sprintf("add_block_1_%s", time.Now().AddDate(0, 0, 2).Format("2006-01-2")): {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
	{
		name:                 "cal-blocks",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		blocks:               1,
	},
	{
		name:                 "cal-res",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		reservations:         1,
	},
}

func TestPostReservationCalendar(t *testing.T) {
	for _, e := range adminPostReservationCalendarTests {
		Repo.App.TestError = false
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/admin/reservation-calendar", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/admin/reservation-calendar", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		now := time.Now()
		bm := make(map[string]int)
		rm := make(map[string]int)
		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2006-01-2")] = 0
			bm[d.Format("2006-01-2")] = 0
		}

		if e.blocks > 0 {
			bm[firstOfMonth.Format("2006-01-2")] = e.blocks
		}

		if e.reservations > 0 {
			rm[lastOfMonth.Format("2006-01-2")] = e.reservations
		}

		session.Put(ctx, "block_map_1", bm)
		session.Put(ctx, "reservation_map_1", rm)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostReservationCalendar)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

	}
}

func getCtx(req *http.Request) context.Context {

	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))

	if err != nil {
		log.Println(err)
	}

	return ctx
}
