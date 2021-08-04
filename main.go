package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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
	osArgsLen := len(os.Args)
	if osArgsLen < 2 {
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

	time.Sleep(broadcastDuration) // sleep for duration between current date

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Create http client
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Put content on file
	resp, err := client.Get(broadcastUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Set timer to save and disable recording
	timer := time.AfterFunc(auditionEnd.Sub(currentTime), func() {
		resp.Body.Close()
		file.Close()

		log.WithFields(log.Fields{"fileName": fileName}).Info("Recorded audio saved to file")

		os.Exit(0)
	})

	// Write data to file
	log.WithFields(log.Fields{"currentTime": time.Now().In(loc)}).Info("Recording started")
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Useless but must have :(
	fmt.Println(timer, size)
}
