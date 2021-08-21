package main

import (
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
	// Duration
	duration time.Duration
	// Saved stream file name prefix
	fileNamePrefix string
}

// Find next broadcast start and end date
/*func (b *broadcast) findNextBroadcast(timeNow time.Time) {
	if (timeNow.Weekday() == time.Saturday) && (timeNow.Hour() < 6) {
		b.start = timeNow
	} else {
		b.start = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 23, 59, 0, 0, timeNow.Location())
		for b.start.Weekday() != time.Friday {
			b.start = b.start.AddDate(0, 0, 1)
		}
	}

	b.end = time.Date(b.start.Year(), b.start.Month(), b.start.Day(), 6, 0, 0, 0, b.start.Location())
	if b.end.Weekday() == time.Friday {
		b.end = b.end.AddDate(0, 0, 1)
	}
}*/

func (b *broadcast) record() error {
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
		file.Close()
		// if err := file.Close(); err != nil {
		// log.Warn(err)
		// }

		resp.Body.Close()
		// if err := resp.Body.Close(); err != nil {
		// log.Warn(err)
		// }
	})

	// Write data to file
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}
