package config

import (
	"fmt"
	"log/slog"
	"slices"
)

// asString returns the string value as-is without modification.
func asString(s string) (string, error) {
	return s, nil
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
