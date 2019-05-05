package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bake/qual-o-mat/qualomat"
	"github.com/namsral/flag"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	data := flag.String("data", "qual-o-mat-data", "Path to qom-data directory")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	s := newServer(qualomat.New(*data))
	log.Fatal(http.ListenAndServe(addr, s))
}
