package config

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"
)

// asString returns the string value as-is without modification.
func asString(s string) (string, error) {
	return s, nil
}

// asBool parses a string representation into a bool.
func asBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// asInt parses a string representation into a int.
func asInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// asInt64 parses a string representation into a int64.
func asInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// asDuration parses a string representation into a time.Duration.
func asDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// asLevel parses a string representation into a slog.Level.
func asLevel(s string) (slog.Level, error) {
	var l slog.Level
	if err := l.UnmarshalText([]byte(s)); err != nil {
		return 0, err
	}
	return l, nil
}

// oneOf returns a validation function that checks
// if a given string matches one of the allowed values.
func oneOf[T ~string](allowed ...T) func(string) (T, error) {
	return func(s string) (T, error) {
		v := T(s)
		if !slices.Contains(allowed, v) {
			var zero T
			return zero, fmt.Errorf("%q is not allowed, want one of %v", s, allowed)
		}
		return v, nil
	}
}
