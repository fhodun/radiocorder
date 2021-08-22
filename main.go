package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Record from now to until duration will pass command handler
func recordNow(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		log.Fatalf("invalid number of arguments, got: %d, want: %d", len(args), 2)
	}

	var (
		err error

		timeNow time.Time = time.Now()
		b       broadcast = broadcast{
			url:   args[0],
			start: timeNow,
		}
	)

	duration, err := time.ParseDuration(args[1])
	if err != nil {
		b.end, err = parseTime(args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b.end = timeNow.Add(duration)
	}

	b.fileNamePrefix, err = parseBroadcastUrl(b.url)
	if err != nil {
		log.Fatal(err)
	}

	if err := b.record(); err != nil {
		log.Fatal(err)
	}
}

// Record from now to until end time will come command handler
func listenToBroadcast(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		log.Fatalf("invalid number of arguments, got: %d, want: %d", len(args), 3)
	}

	/*start, err := parseTime(args[1])
	if err != nil {
		log.Fatal(err)
	}

	end, err := parseTime(args[2])
	if err != nil {
		log.Fatal(err)
	}*/

	var (
		timeNow time.Time = time.Now()
		b       broadcast = broadcast{
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
		b.start = timeNow.Add(startDuration)
	}

	// TODO: repeated code
	endDuration, err := time.ParseDuration(args[2])
	if err != nil {
		b.end, err = parseTime(args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b.end = timeNow.Add(endDuration)
	}

	b.fileNamePrefix, err = parseBroadcastUrl(b.url)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Until(b.start))

	if err := b.record(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	cmdNow := &cobra.Command{
		Use:     "now [host] [end time/duration]",
		Aliases: []string{"n"},
		Short:   "Record broadcast from now",
		Example: "now example.com:2137/stream 2h13m7s",
		Run:     recordNow,
	}

	cmdBroadcast := &cobra.Command{
		Use:     "broadcast [host] [start time/duration] [end time/duration]", // TODO: add [.../duration]
		Aliases: []string{"b"},
		Short:   "Record broadcast",
		Example: "broadcast example.com:2137/stream \"Fri 23:59\" \"Sat 6:00\"",
		Run:     listenToBroadcast,
	}

	rootCmd := &cobra.Command{Use: "radiocorder"}
	rootCmd.AddCommand(cmdNow, cmdBroadcast)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
