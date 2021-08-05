package main

import (
	"net/url"
	"strings"
	"time"
)

// Find next broadcast start and end date
func FindNextBroadcast(timeNow time.Time) (time.Time, time.Time) {
	auditionStart := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 23, 59, 0, 0, timeNow.Location())
	for auditionStart.Weekday() != time.Friday {
		auditionStart = auditionStart.AddDate(0, 0, 1)
	}

	auditionEnd := time.Date(auditionStart.Year(), auditionStart.Month(), auditionStart.Day()+1, 6, 0, 0, 0, auditionStart.Location())

	return auditionStart, auditionEnd
}

// Parse broadcast url to name of broadcast file
func ParseBroadcastUrl(broadcastUrl string) (string, error) {
	// Build fileName from fullPath
	fileURL, err := url.Parse(broadcastUrl)
	if err != nil {
		return "", err
	}

	// Parse to fileName
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1] + ".ogg"

	return fileName, nil
}
