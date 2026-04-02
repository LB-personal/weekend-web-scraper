package domain

import (
	"errors"
	"strings"
)

type Rating struct {
	value uint8
}

func NewRateing(v uint8) (Rating, error) {
	if v >= 1 && v <= 5 {
		return Rating{
			value: v,
		}, nil
	}

	return Rating{}, errors.New("rating can only be btween 1 to 5")
}

func (r Rating) String() string {
	var b strings.Builder
	for range r.value {
		b.WriteRune('*')
	}

	return b.String()
}
