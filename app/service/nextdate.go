package service

import (
	"errors"
	"main/app/model"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty string")
	} else if strings.Contains(repeat, "d ") {
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil {
			return "", err
		}
		if days > 400 {
			return "", errors.New("maximum days count must be 400")
		}

		parsedDate, err := time.Parse(model.DatePattern, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(0, 0, days)

		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, days)
		}

		return newDate.Format(model.DatePattern), nil
	} else if repeat == "y" {
		parsedDate, err := time.Parse(model.DatePattern, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(1, 0, 0)

		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}

		return newDate.Format(model.DatePattern), nil
	} else {
		return "", errors.New("repeat wrong format")
	}
}
