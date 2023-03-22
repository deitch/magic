package main

import (
	"log"
	"os"

	"github.com/deitch/magic/pkg/magic"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <file>", os.Args[0])
	}
	file := os.Args[1]
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("error opening file %s: %v", file, err)
	}
	defer f.Close()
	info, err := magic.GetType(f)
	if err != nil {
		log.Fatalf("error getting type of file %s: %v", file, err)
	}
	log.Printf("file %s is type %v", file, info)
}
