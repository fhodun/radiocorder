package main

import (
	"time"

	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Record from now to until duration will pass command handler
func runNow(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		log.Fatalf("invalid number of arguments, got: %d, want: %d", len(args), 2)
	}

	var (
		b broadcast = broadcast{
			url:   args[0],
			start: time.Now(),
		}
	)

	duration, err := time.ParseDuration(args[1])
	if err != nil {
		b.end, err = parseTime(args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b.end = b.start.Add(duration)
	}

	b.fileNamePrefix, err = parseBroadcastUrl(b.url)
	if err != nil {
		log.Fatal(err)
	}

	record(&b)
}

// Record from now to until end time will come command handler
func runBroadcast(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		log.Fatalf("invalid number of arguments, got: %d, want: %d", len(args), 3)
	}

	var (
		b broadcast = broadcast{
			url: args[0],
		}
	)

	startDuration, err := time.ParseDuration(args[1])
	if err != nil {
		b.start, err = parseTime(args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b.start = time.Now().Add(startDuration)
	}

	// TODO: repeated code
	endDuration, err := time.ParseDuration(args[2])
	if err != nil {
		b.end, err = parseTime(args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b.end = b.start.Add(endDuration)
	}

	b.fileNamePrefix, err = parseBroadcastUrl(b.url)
	if err != nil {
		log.Fatal(err)
	}

	// Create progress bar from duration between actual time and broadcast start
	bar := pb.New(int(time.Until(b.start).Seconds()))
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

	time.Sleep(time.Until(b.start))

	done <- true
	bar.Finish()

	record(&b)
}

func main() {
	cmdNow := &cobra.Command{
		Use:     "now [host] [end time/duration]",
		Aliases: []string{"n"},
		Short:   "Record broadcast from now",
		Example: "now example.com:2137/stream 2h13m7s",
		Run:     runNow,
	}

	cmdBroadcast := &cobra.Command{
		Use:     "broadcast [host] [start time/duration] [end time/duration]", // TODO: add [.../duration]
		Aliases: []string{"b"},
		Short:   "Record broadcast",
		Example: "broadcast example.com:2137/stream \"Fri 23:59\" \"Sat 6:00\"",
		Run:     runBroadcast,
	}

	rootCmd := &cobra.Command{Use: "radiocorder"}
	rootCmd.AddCommand(cmdNow, cmdBroadcast)

					

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
