package fix

import (
	"errors"
	"testing"
)

func Test(t *testing.T) {
	err := errors.New("test")

	Hook(err, func(err error) bool {
		return true
	})

	Try(&err, func(err error) error {
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}
