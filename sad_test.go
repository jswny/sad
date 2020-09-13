package sad_test

import (
	"testing"

	testutils "github.com/jswny/sad/internal"

	"github.com/jswny/sad"
)

func TestGetFullName(t *testing.T) {
	opts := sad.Options{
		Name:    "foo",
		Channel: "beta",
	}

	fullName := sad.GetFullName(&opts)
	expected := "foo-beta"

	testutils.CompareStrings(expected, fullName, "full name", t)
}
