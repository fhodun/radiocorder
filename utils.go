package main

import (
	"net/url"
	"strings"
	"time"
)

// Find next broadcast
func FindNextBroadcast(currentTime time.Time, loc *time.Location) time.Time {
	auditionEnd := currentTime
	for auditionEnd.Weekday() != time.Friday {
		auditionEnd = time.Date(auditionEnd.Year(), auditionEnd.Month(), auditionEnd.Day()+1, 23, 59, 0, 0, loc)
	}

	return auditionEnd
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
