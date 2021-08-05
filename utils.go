package main

import (
	"net/url"
	"strings"
	"time"
)

// Find next broadcast start and end date
func FindNextBroadcast(currentTime time.Time) (time.Time, time.Time) {
	auditionStart := currentTime
	for auditionStart.Weekday() != time.Friday {
		auditionStart = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, 23, 59, 0, 0, currentTime.Location())
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
