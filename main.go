package main

import (
	"flag"
	"fmt"

	service "github.com/ommz/graphql-test/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	// accept -n as an optional CLI flag
	// nFlag is a pointer. Its value defaults to 5
	nFlag := flag.Uint64("n", 5, "number of projects to fetch")
	flag.Parse()

	// zero-check it to avoid panic(s) in service layer
	if *nFlag == 0 {
		log.Fatal("-n value cannot be less than 1")
	}

	namesCSV, forkSum := service.CallGraphQLAPI(*nFlag)

	fmt.Println("Names: ", namesCSV, "\r\nForks Sum: ", forkSum)
}
