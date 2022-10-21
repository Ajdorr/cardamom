package testing_ext

import (
	"reflect"
	"testing"
)

type Results struct {
	T            *testing.T
	IsSuccessful bool
}

func (r *Results) FailNowIfUnsuccessful() {
	if !r.IsSuccessful {
		r.T.FailNow()
	}
}

func TestEqual(t *testing.T, a, b any, failureMessage string, messageArgs ...any) *Results {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		t.Errorf(failureMessage, messageArgs...)
		return &Results{T: t, IsSuccessful: false}
	}

	if a != b {
		t.Errorf(failureMessage, messageArgs...)
		return &Results{T: t, IsSuccessful: false}
	}
	return &Results{T: t, IsSuccessful: true}
}
