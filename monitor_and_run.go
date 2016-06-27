package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

				// exec.Command cannot be reused so declare it in the
				// loop
				cmd := exec.Command(script)
				// Connect to cmd's stdout
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Fatal(err)
				}
				// Start cmd
				if err := cmd.Start(); err != nil {
					log.Println("Error running script.", err)
				} else {
					// It started
					log.Println("Script ran successfully.")
					out, err := ioutil.ReadAll(stdout)
					if err != nil {
						log.Println("error: ioutil.ReadAll failed: ", err)
					} else {
						// Print cmd's stdout to STDOUT
						fmt.Println(string(out))
					}

					if err := cmd.Wait(); err != nil {
						log.Println("error: cmd.Wait() failed: ", err)
					}
				}
			}
		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}
}
