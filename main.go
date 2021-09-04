package main

import (
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Check provided args
func preRunE(cmd *cobra.Command, args []string, expectedArgsLen int) error {
	if len(args) < expectedArgsLen {
		return fmt.Errorf("invalid number of arguments, got: %d, want: %d", len(args), expectedArgsLen)
	}

	return nil
}

// Record from now to until duration will pass command handler
func runNow(cmd *cobra.Command, args []string) {
	var (
		b broadcast = broadcast{
			url:     args[0],
			start:   time.Now(),
			started: false,
		}
	)

	duration, err := time.ParseDuration(args[1])
	if err != nil {
		b.end, err = parseTime(
			args[1],
			b.started,
			time.Time{},
		)
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
	var (
		b broadcast = broadcast{
			url:     args[0],
			started: flagToBool(cmd.Flag("started").Value.String()),
		}
	)

	startDuration, err := time.ParseDuration(args[1])
	if err != nil {
		b.start, err = parseTime(
			args[1],
			b.started,
			time.Time{},
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b.start = time.Now().Add(startDuration)
	}

	// TODO: repeated code
	endDuration, err := time.ParseDuration(args[2])
	if err != nil {
		b.end, err = parseTime(
			args[2],
			b.started,
			b.start,
		)
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

	if !b.started && time.Now().Before(b.start) {
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

		log.WithField("duration", time.Until(b.start)).Info("waiting for broadcast start")
		time.Sleep(time.Until(b.start))

		done <- true
		bar.Finish()
	}

	record(&b)
}

func main() {
	cmdNow := &cobra.Command{
		Use:     "now [host] [end time/duration]",
		Aliases: []string{"n"},
		Short:   "Record broadcast from now",
		Example: "now example.com:2137/stream 2h13m7s",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := preRunE(cmd, args, 2); err != nil {
				return err
			}
			return nil
		},
		Run: runNow,
	}

	cmdBroadcast := &cobra.Command{
		Use:     "broadcast [host] [start time/duration] [end time/duration]",
		Aliases: []string{"b"},
		Short:   "Record next planned broadcast",
		Example: "broadcast example.com:2137/stream \"Fri 23:59\" \"Sat 6:00\"",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := preRunE(cmd, args, 2); err != nil {
				return err
			}
			return nil
		},
		Run: runBroadcast,
	}
	cmdBroadcast.Flags().BoolP("started", "s", false, "record broadcast even if it started but end time is until actual")

	rootCmd := &cobra.Command{Use: "radiocorder"}
	rootCmd.AddCommand(cmdNow, cmdBroadcast)

	// TODO: retry flag
	// rootCmd.PersistentFlags().BoolP("retry", "r", false, "retry recording after unplanned fatal until end")

	rootCmd.Execute()
}
