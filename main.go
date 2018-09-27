package main

import (
	"github.com/juju/errors"
	"log"
	"os"
	"strconv"
)

func readPort() (int, error) {
	var port, err = strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		return 0, errors.Annotate(err, "COULD_NOT_PARSE_PORT")
	}
	return port, nil
}

func main() {
	var port, err = readPort()
	if err != nil {
		err = errors.Trace(err)
		log.Fatalf("Could not read the HTTP listening port.\n%v", err)
	}
	var server = Make(port)
	err = server.ListenAndServe()
	log.Fatalf("The HTTP server has shut down unexpectedly: %v", err)
}
