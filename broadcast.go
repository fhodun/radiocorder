package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

type broadcast struct {
	// Stream url
	url string
	// Start time
	start time.Time
	// End time
	end time.Time
	// Saved stream file name prefix
	fileNamePrefix string
}

func (b broadcast) record() error {
	// Get audio from host
	resp, err := http.Get(b.url)
	if err != nil {
		return err
	}

	// Create file name from prefix and start date
	fileName := b.fileNamePrefix + b.start.Format("2006-01-02") + ".ogg"

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	// Create timer to close file after audition end
	time.AfterFunc(time.Until(b.end), func() {
		if err := file.Close(); err != nil {
			log.Warn(err)
		}

		if err := resp.Body.Close(); err != nil {
			log.Warn(err)
		}
	})

	// Write data to file
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}
