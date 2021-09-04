package main

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Parse string to bool
func flagToBool(f string) bool {
	var b bool

	switch f {
	case "true":
		b = true
	case "false":
		b = false
	}

	return b
}

// Parse broadcast url to name of broadcast file
func parseBroadcastUrl(broadcastUrl string) (string, error) {
	fileURL, err := url.Parse(broadcastUrl)
	if err != nil {
		return "", err
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	return fileName, nil
}

// Parse weekday shortcut to time
func parseTime(s string, started bool, after time.Time) (time.Time, error) {
	daysOfWeek := map[string]time.Weekday{
		"mon": time.Monday,
		"tue": time.Tuesday,
		"wed": time.Wednesday,
		"thu": time.Thursday,
		"fri": time.Friday,
		"sat": time.Saturday,
		"sun": time.Sunday,
	}

	s = strings.ToLower(s)
	d := strings.SplitN(s, ", ", 2)
	timesString := strings.SplitN(d[1], ":", 2)
	var timesInt [2]int
	for index, tString := range timesString {
		timesInt[index], _ = strconv.Atoi(tString)
	}

	weekTime, ok := daysOfWeek[d[0]] // TODO: can make it better
	if !ok {
		return time.Time{}, errors.New("invalid weekday name")
	}

	timeNow := time.Now()

	t := time.Date(
		timeNow.Year(),
		timeNow.Month(),
		timeNow.Day(),
		timesInt[0],
		timesInt[1],
		0,
		0,
		timeNow.Location(),
	)

	// Add one day to valid date
	if after.IsZero() || timeNow.Before(after) {
		if t.Weekday() == weekTime && !started {
			t = t.AddDate(0, 0, 1)
		}
	}

	// Set t weekday to next day named as in args
	for t.Weekday() != weekTime {
		t = t.AddDate(0, 0, 1)
	}

	return t, nil
}
