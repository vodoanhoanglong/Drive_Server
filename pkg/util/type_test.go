package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGender(t *testing.T) {

	for i, ip := range []struct {
		Input   json.RawMessage
		IsError bool
		Output  Gender
	}{
		{
			[]byte("\"male\""),
			false,
			Male,
		}, {
			[]byte("\"female\""),
			false,
			Female,
		}, {
			[]byte("\"\""),
			true,
			Male,
		},
	} {
		var result Gender
		err := json.Unmarshal(ip.Input, &result)
		if ip.IsError {
			assert.EqualError(t, err, fmt.Sprintf("invalid gender `%s`", strings.Trim(string(ip.Input), "\"")), "%d", i)
		} else {
			assert.Nil(t, err, "%d", i)
			assert.Equal(t, ip.Output.String(), result.String(), "%d", i)
		}
	}
}

func TestDate(t *testing.T) {

	for i, ip := range []struct {
		Input   json.RawMessage
		IsError bool
		Output  *Date
	}{
		{
			[]byte("\"2020-01-01\""),
			false,
			&Date{2020, 1, 1},
		}, {
			[]byte("\"2020-02-29\""),
			false,
			&Date{2020, 2, 29},
		}, {
			[]byte("\"2020-02-30\""),
			true,
			nil,
		}, {
			[]byte("\"\""),
			true,
			nil,
		},
	} {
		var result Date
		err := json.Unmarshal(ip.Input, &result)
		if ip.IsError {
			assert.EqualError(t, err, fmt.Sprintf("invalid date `%s`", strings.Trim(string(ip.Input), "\"")), "%d", i)
		} else {
			assert.Nil(t, err, "%d", i)
			assert.Equal(t, ip.Output.String(), result.String(), "%d", i)
		}
	}
}

func TestUniqueStrings(t *testing.T) {

	fixtures := []struct {
		Input    []string
		Expected string
	}{
		{Input: []string{}, Expected: ""},
		{Input: []string{"a", "b", "c"}, Expected: "a,b,c"},
		{Input: []string{"a", "b", "b", "c"}, Expected: "a,b,c"},
		{Input: []string{"a", "b", "c", "a"}, Expected: "a,b,c"},
		{Input: []string{"c", "b", "c", "a", "b", "c"}, Expected: "a,b,c"},
	}

	for i, ss := range fixtures {
		sample := UniqueStrings{}
		for _, s := range ss.Input {
			sample.Add(s)
		}
		assert.Equal(t, ss.Expected, sample.String(), i)
	}
}
