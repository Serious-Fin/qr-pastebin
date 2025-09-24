package shares

import (
	"testing"
	"time"
)

func TestExpireInCreation1(t *testing.T) {
	expireAt := time.Now().AddDate(1, 2, 3)
	want := "Expires in 1 year 2 months and 3 days"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation2(t *testing.T) {
	expireAt := time.Now().AddDate(2, 1, 3)
	want := "Expires in 2 years 1 month and 3 days"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation3(t *testing.T) {
	expireAt := time.Now().AddDate(3, 2, 1)
	want := "Expires in 3 years 2 months and 1 day"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation4(t *testing.T) {
	expireAt := time.Now().AddDate(0, 1, 2)
	want := "Expires in 1 month and 2 days"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation5(t *testing.T) {
	expireAt := time.Now().AddDate(1, 0, 2)
	want := "Expires in 1 year and 2 days"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation6(t *testing.T) {
	expireAt := time.Now().AddDate(1, 2, 0)
	want := "Expires in 1 year and 2 months"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation7(t *testing.T) {
	expireAt := time.Now().AddDate(1, 0, 0)
	want := "Expires in 1 year"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation8(t *testing.T) {
	expireAt := time.Now().AddDate(0, 1, 0)
	want := "Expires in 1 month"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation9(t *testing.T) {
	expireAt := time.Now().AddDate(0, 0, 1)
	want := "Expires in 1 day"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation10(t *testing.T) {
	expireAt := time.Now().Add(time.Second)
	want := "Expires today"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}

func TestExpireInCreation11(t *testing.T) {
	expireAt := time.Now().Add(-time.Second)
	want := "Already expired"
	got := createExpireInTextFromDate(expireAt)
	if got != want {
		t.Errorf(`expected "%s", got "%s"`, want, got)
	}
}
