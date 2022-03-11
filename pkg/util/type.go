package util

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

// MessageOutput represent simple message response
type MessageOutput struct {
	Message string `json:"message"`
}

// AffectedRowsOutput represent affected rows response
type AffectedRowsOutput struct {
	AffectedRows int `json:"affected_rows"`
}

func ParseGender(input string) (*Gender, error) {
	var result Gender

	switch input {
	case string(Male):
		result = Male
	case string(Female):
		result = Female
	default:
		return nil, fmt.Errorf("invalid gender `%s`", input)
	}

	return &result, nil
}

func (g *Gender) UnmarshalJSON(b []byte) error {
	result, err := ParseGender(strings.Trim(string(b), "\""))

	if err != nil {
		return err
	}
	*g = *result
	return nil
}

func (g Gender) String() string {
	return string(g)
}

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d *Date) UnmarshalJSON(b []byte) error {
	result, err := ParseDate(strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	*d = *result
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func ParseDate(input string) (*Date, error) {
	parts := strings.Split(input, "-")
	invalidError := fmt.Errorf("invalid date `%s`", input)

	if len(parts) != 3 {
		return nil, invalidError
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil || year < 0 || year > 9999 {
		return nil, invalidError
	}

	if year < 0 {
		return nil, invalidError
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil || month <= 0 || month > 12 {
		return nil, invalidError
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil || day <= 0 ||
		((month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12) && day > 31) ||
		((month == 4 || month == 6 || month == 9 || month == 11) && day > 30) ||
		(month == 2 && year%4 > 0 && day > 28) ||
		(month == 2 && year%4 == 0 && day > 29) {
		return nil, invalidError
	}

	return &Date{
		Year:  year,
		Month: month,
		Day:   day,
	}, nil
}

type PGArrayString []string

func (pas *PGArrayString) UnmarshalJSON(b []byte) error {
	result, err := DecodePostgresArray(strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	*pas = result
	return nil
}

func (pas *PGArrayString) MarshalJSON() ([]byte, error) {
	if pas == nil {
		return nil, nil
	}
	return json.Marshal(EncodePostgresArray(*pas))
}

// UniqueStrings is the special array string that only store unique values
type UniqueStrings map[string]bool

// Add append new value or skip if it's existing
func (us UniqueStrings) Add(values ...string) {
	for _, s := range values {
		if _, ok := us[s]; !ok {
			us[s] = true
		}
	}
}

// IsEmpty check if the array is empty
func (us UniqueStrings) IsEmpty() bool {
	return len(us) == 0
}

// Value return
func (us UniqueStrings) Value() []string {
	results := make([]string, 0, len(us))
	for k := range us {
		results = append(results, k)
	}
	return results
}

// String implement string interface
func (us UniqueStrings) String() string {
	results := us.Value()
	sort.Strings(results)
	return strings.Join(results, ",")
}
