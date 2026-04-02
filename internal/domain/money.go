package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	resoulotion = 100 // scraped site uses only pounds so a const is enough
)

type Money uint64

// Creates a new money instance, this function is build to handle pound book prices and nothing more
func NewMoney(s string) (Money, error) {
	p := strings.Split(s, ".")
	if len(p) != 2 {
		return 0, errors.New("money format is invalid")
	}
	major := p[0]
	minor, found := strings.CutSuffix(p[1], "£")
	if !found {
		return 0, errors.New("money format is invalid")
	}

	ma, err := strconv.Atoi(major)
	if err != nil {
		return 0, errors.Join(errors.New("money format is invalid"), err)
	}

	mi, err := strconv.Atoi(minor)

	if err != nil {
		return 0, errors.Join(errors.New("money format is invalid"), err)
	}

	if ma < 0 || mi < 0 {
		return 0, errors.New("money format is invalid")
	}

	return Money(ma*resoulotion + mi), nil
}

func (m Money) String() string {
	return fmt.Sprintf("%d.%d£", m/resoulotion, m%resoulotion)
}
