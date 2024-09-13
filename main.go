package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// channel to read reload request
var reloadChan = make(chan bool)

func main() {
	// serve the static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// End point for reload notification
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		// Read the incoming request
		<-reloadChan
		w.Write([]byte("reload")) // sends reload info to /reload route
	})

	// Concurrent function to watch file changes
	go scanFileChanges()

	// start the server
	log.Println(`Starting the serve at: 8000`)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// scan for file changes
func scanFileChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Listen for events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("events:", event)
			if event.Op&fsnotify.Write == fsnotify.Write && isTrackedFile(event.Name) {
				log.Println("Modified File:", event.Name)
				reloadChan <- true
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error :", err)
		}
	}

}

// Check if this is correct file extension or not
func isTrackedFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))

	return ext == ".html" || ext == ".css" || ext == ".js"
}
