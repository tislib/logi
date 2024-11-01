package main

import (
	"strings"
)

var result float64

func ParseJson(d string, debug bool) (interface{}, error) {
	s := NewScanner(strings.NewReader(d), debug)
	parser := yyNewParser()

	parser.Parse(s)

	if s.err != nil {
		return nil, s.err
	}
	return result, nil
}

func statement(data interface{}) {
	result = toFloat64(data)
}

func toFloat64(i interface{}) float64 {
	switch v := i.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	case int:
		return float64(v)
	}

	return 0
}
