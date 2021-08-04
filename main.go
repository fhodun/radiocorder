package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Find next broadcast
func findNextBroadcast(currentTime time.Time, loc *time.Location) time.Time {
	auditionEnd := currentTime
	for int(auditionEnd.Weekday()) != 5 {
		auditionEnd = time.Date(auditionEnd.Year(), auditionEnd.Month(), auditionEnd.Day()+1, 23, 59, 0, 0, loc)
	}

	return auditionEnd
}

// Parse broadcast url to name of broadcast file
func parseBroadcastUrl(broadcastUrl string) (fileName string) {
	// Build fileName from fullPath
	fileURL, err := url.Parse(broadcastUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Parse to fileName
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1] + ".ogg"

	return
}

func main() {
	if osArgsLen := len(os.Args); osArgsLen < 2 {
		log.WithFields(log.Fields{"len(os.Args)": osArgsLen}).Fatal("No arguments provided")
	}

	broadcastUrl := string(os.Args[1])
	fileName := parseBroadcastUrl(broadcastUrl)

	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		log.Fatal(err)
	}
	currentTime := time.Now().In(loc)

	auditionStart := findNextBroadcast(currentTime, loc)
	auditionEnd := time.Date(auditionStart.Year(), auditionStart.Month(), auditionStart.Day()+1, 6, 0, 0, 0, loc)
	broadcastDuration := auditionStart.Sub(currentTime)
	log.WithFields(log.Fields{"date": auditionStart, "broadcastDuration": broadcastDuration}).Info("Found next broadcast date")

	// Sleep for duration between current date and broadcast
	time.Sleep(broadcastDuration)

	// Get audio from host
	resp, err := http.Get(broadcastUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"currentTime": time.Now().In(loc)}).Info("Recording started")

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Create timer to close file after audition end
	_ = time.AfterFunc(auditionEnd.Sub(time.Now().In(loc)), func() {
		defer resp.Body.Close()
		defer log.WithFields(log.Fields{"fileName": fileName}).Info("Recorded audio saved to file")
		file.Close()
		
		os.Exit(0)
	})

	// Write data to file
	if _, err := io.Copy(file, resp.Body); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Connection with host closed")
	}
}
