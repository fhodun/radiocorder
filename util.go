package main

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

)

// Parse broadcast url to name of broadcast file
func parseBroadcastUrl(broadcastUrl string) (string, error) {
	fileURL, err := url.Parse(broadcastUrl)
	if err != nil {
		return "", err
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1] + "_"

	return fileName, nil
}

// Parse weekday shortcut to time weekday type
func parseTime(s string) (time.Time, error) {
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

	t := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), timesInt[0], timesInt[1], 0, 0, timeNow.Location())
	for t.Weekday() != weekTime {
		t = t.AddDate(0, 0, 1)
	}

	return t, nil
}
