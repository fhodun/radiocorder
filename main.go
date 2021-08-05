package main

import (
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	if osArgsLen := len(os.Args); osArgsLen < 2 {
		log.WithFields(log.Fields{"len(os.Args)": osArgsLen}).Fatal("No arguments provided")
	}

	broadcastUrl := string(os.Args[1])
	fileName, err := ParseBroadcastUrl(broadcastUrl)
	if err != nil {
		log.WithFields(log.Fields{"broadcastUrl": broadcastUrl}).Fatal(err)
	}

	currentTime := time.Now()

	auditionStart, auditionEnd := FindNextBroadcast(currentTime)
	broadcastDuration := auditionStart.Sub(currentTime)
	log.WithFields(log.Fields{"date": auditionStart, "broadcastDuration": broadcastDuration}).Info("Found next broadcast date")

	// Sleep for duration between current date and broadcast
	time.Sleep(broadcastDuration)

	// Get audio from host
	resp, err := http.Get(broadcastUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"currentTime": time.Now()}).Info("Recording started")

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Create timer to close file after audition end
	_ = time.AfterFunc(time.Until(auditionEnd), func() {
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
