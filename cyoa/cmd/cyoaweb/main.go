package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kinobi/gophercises.com/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "the port where CYOA webserver listen to")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("failed to open JSON story file: %s", err)
	}

	story, err := cyoa.JSONStory(f)
	if err != nil {
		log.Fatalf("failed to parse JSON story file: %s", err)
	}

	h := cyoa.NewHandler(story)
	log.Printf("listening on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
