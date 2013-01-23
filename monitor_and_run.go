package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"os/exec"
	"flag"
)

var debug bool
var directory string
var script string
var log_file string

func init() {
	flag.BoolVar(&debug, "debug", false, "Set to true output log to STDOUT.")
	flag.StringVar(&directory, "directory", "", "The directory to watch for changes. (Required)")
	flag.StringVar(&script, "script", "", "The path to the script to execute on changes.")
  flag.StringVar(&log_file, "log_file", "/tmp/monitor_and_run.log", "Where to create log file. Ignored when debug is true")
	flag.Parse()
}

func main() {
	block := make(chan int)

	if !debug {
		log_file, err := os.OpenFile(log_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer log_file.Close()

		log.SetOutput(log_file)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go monitor_downloads(watcher)

	<-block

	log.Println("I should not be here")
}

func monitor_downloads(watcher *fsnotify.Watcher) {
	run_script := true

	if "" == directory {
		flag.PrintDefaults()
		log.Fatal("No directory specified")
	}

	if "" == script {
		log.Println("No script specified")
		run_script = false
	}

	for {
		err := watcher.Watch(directory)
		if err != nil {
			log.Fatal(err)
		}

		select {
		case ev := <-watcher.Event:
			log.Println("event:", ev)
			if run_script {
				log.Println("Running script...")
				err := exec.Command(script).Run()
				if err != nil {
					log.Println("Error running script.", err)
				} else {
					log.Println("Script ran successfully.")
				}
			}
		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}
}
