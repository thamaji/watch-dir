package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func ShowVersion() {
	fmt.Println("v1.0.0")
}

func run(filter fsnotify.Op, args ...string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	for _, arg := range args {
		if err := watcher.Add(arg); err != nil {
			return err
		}
	}

	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGINT)

	for {
		select {
		case <-sigc:
			return nil

		case event := <-watcher.Events:
			if event.Op&filter != 0 {
				fmt.Fprintln(os.Stdout, event.Name)
			}

		case err := <-watcher.Errors:
			return err
		}
	}
}

type StringSlice []string

func (slice *StringSlice) String() string {
	return strings.Join(*slice, ",")
}

func (slice *StringSlice) Set(value string) error {
	*slice = append(*slice, value)
	return nil
}

func main() {
	var help, version bool
	var events StringSlice

	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&version, "v", false, "show version")
	flag.Var(&events, "e", "set watch event (CREATAE|WRITE|REMOVE|RENAME|CHMOD)")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: " + os.Args[0] + " [OPTIONS] DIR [DIR...]")
		fmt.Println()
		fmt.Println("Watch file system events")
		fmt.Println()
		fmt.Println("Options:")
		flag.CommandLine.PrintDefaults()
		fmt.Println()
	}

	flag.Parse()

	args := flag.Args()

	if help {
		flag.Usage()
		return
	}

	if version {
		ShowVersion()
		return
	}

	if len(args) <= 0 {
		flag.Usage()
		return
	}

	var filter fsnotify.Op
	for _, event := range events {
		switch strings.ToUpper(event) {
		case "CREATE":
			filter |= fsnotify.Create
		case "WRITE":
			filter |= fsnotify.Write
		case "REMOVE":
			filter |= fsnotify.Remove
		case "RENAME":
			filter |= fsnotify.Rename
		case "CHMOD":
			filter |= fsnotify.Chmod
		default:
			fmt.Fprintln(os.Stderr, "unknwon event: "+event)
			flag.Usage()
			os.Exit(1)
		}
	}

	if filter == 0 {
		filter = fsnotify.Create | fsnotify.Write | fsnotify.Remove | fsnotify.Rename | fsnotify.Chmod
	}

	if err := run(filter, flag.Args()...); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
