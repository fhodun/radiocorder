package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
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
	// Define is broadcast started
	started bool
}

func (b broadcast) createFile() (*os.File, error) {
	name := fmt.Sprintf("%s_%s_%d-%d",
		b.fileNamePrefix,
		b.start.Format("2006-01-02"),
		b.start.Hour(), b.start.Minute(),
	)

	for i := 1; func() bool {
		_, err := os.Stat(name + ".ogg")
		return os.IsExist(err)
	}(); i++ {
		log.Infof("File name \"%s\" already exists, retrying", fmt.Sprint(name, ".ogg"))
		name = fmt.Sprintf("%s_%d", name, i)
	}

	file, err := os.Create(name + ".ogg")
	if err != nil {
		return nil, err
	}

	return file, nil
}

func record(b *broadcast) {
	var (
		file     *os.File
		bar      *pb.ProgressBar
		fileInfo fs.FileInfo
		done     chan bool = make(chan bool, 1)
	)

	log.Info("starting recording")

	// Get audio from host
	resp, err := http.Get(b.url)
	if err != nil {
		log.Fatal(err)
	}

	file, err = b.createFile()
	if err != nil {
		log.Fatal(err)
	}

	// Create timer to close file after audition end
	time.AfterFunc(time.Until(b.end), func() {
		done <- true
		bar.Finish()

		fileInfo, err = file.Stat()
		if err != nil {
			log.Warn(err)
		}

		if err := resp.Body.Close(); err != nil {
			log.Warn(err)
		}

		if err := file.Close(); err != nil {
			log.Warn(err)
		}
	})

	// Create progress bar from duration between actual time and broadcast end
	bar = pb.New(int(time.Until(b.end).Seconds()))
	bar.Start()

	ticker := time.NewTicker(1 * time.Second)
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
		defer bar.Finish()
		done <- true

		log.Warn(err)
	}

	log.Infof("written broadcast to file with name: %s and size of %d bytes", fileInfo.Name(), fileInfo.Size())
}
