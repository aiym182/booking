package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {

	r := httptest.NewRequest("POST", "/somewhere", nil)

	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("Got invalid form")
	}

}
func TestFrom_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/somewhere", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required field is missing")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/somewhere", nil)
	r.PostForm = postedData

	form = New(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("Form does not have requried fields when it is supposed to ")
	}

}

func TestHas(t *testing.T) {
	r := httptest.NewRequest("POST", "/somewhere", nil)
	form := New(r.PostForm)

	h := form.Has("Whatever", r)

	if h {
		t.Error("Form shows has field when it is not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	h = form.Has("a", r)

	if !h {
		t.Error("Form doesn't have any field when it is supoosed to ")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/somewhere", nil)

	form := New(r.PostForm)

	if form.MinLen("anything", 3) {
		t.Error("Form shows has 3 fields when it has none")
	}

	postedData := url.Values{}
	postedData.Add("some_field", "some value")

	form = New(postedData)

	form.MinLen("some_field", 100)

	if form.Valid() {
		t.Error("Shows minlength of 100 met when data is shorter")
	}

}

func TestValidateEmail(t *testing.T) {

	r := httptest.NewRequest("POST", "/something", nil)

	form := New(r.PostForm)

	form.ValidateEmail("X")
	if form.Valid() {
		t.Error("Email is valid when its not supposed to be validated")
	}
	postedData := url.Values{}
	postedData.Add("email", "aiym182@gmail.com")
	form = New(postedData)

	form.ValidateEmail("email")

	if !form.Valid() {
		t.Error("Email is not valid when it is supposed to ")
	}

}
