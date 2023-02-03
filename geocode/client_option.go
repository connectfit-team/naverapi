package geocode

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type QueryOption func(url.Values)

// WithCoordinate set coordinates to be the center of the search.
// If set, computes the distance from the `Query()` value to the coordinates.
func WithCoordinate(lon, lat float64) QueryOption {
	return func(data url.Values) {
		data.Set(coordinateConst, fmt.Sprintf("%f,%f", lon, lat))
	}
}

// WithLanguage set language query option.
// default : "kor"
func WithLanguage(lang Lang) QueryOption {
	return func(data url.Values) {
		data.Set(languageConst, string(lang))
	}
}

// WithHCode set filter condition.
// you can filter array values through repetitive codes.
// but if previous filter and the current filter are different,
// the previous filter will be removed.
func WithHCode(codes ...string) QueryOption {
	return func(data url.Values) {
		data.Set(filterConst, fmt.Sprintf("%s@%s", hCode, strings.Join(codes, ";")))
	}
}

// WithBCode set filter condition.
// you can filter array values through repetitive codes.
// but if previous filter and the current filter are different,
// the previous filter will be removed.
func WithBCode(codes ...string) QueryOption {
	return func(data url.Values) {
		data.Set(filterConst, fmt.Sprintf("%s@%s", bCode, strings.Join(codes, ";")))
	}
}

// WithPage decide page you want.
// default	: 1
func WithPage(page int) QueryOption {
	return func(data url.Values) {
		data.Set(pageConst, strconv.Itoa(page))
	}
}

// WithCount decide page unit.
// default 	: 10
// range 	: 1 ~ 100
func WithCount(count int) QueryOption {
	return func(data url.Values) {
		data.Set(countConst, strconv.Itoa(count))
	}
}
