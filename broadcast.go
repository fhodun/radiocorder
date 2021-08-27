package main

import (
	"github.com/cheggaaa/pb/v3"
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
	var (
		// Create file name from prefix and start date
		fileName string = b.fileNamePrefix + b.start.Format("2006-01-02") + ".ogg"
		file     *os.File
		bar      *pb.ProgressBar
	)

	// Get audio from host
	resp, err := http.Get(b.url)
	if err != nil {
		return err
	}

	// Create blank file
	file, err = os.Create(fileName)
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

	// Create progress bar from duration between actual time and broadcast end
	bar = pb.New(int(time.Until(b.end).Seconds()))
	bar.Start()

	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				bar.Increment()
			}
		}
	}()

	// Write data to file
	if _, err := io.Copy(file, resp.Body); err != nil {
		log.Warn(err)
	}

	done <- true
	bar.Finish()

	return nil
}
