package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	filePath := flag.String("f", "urls.txt", "Path to file containing URLs")
	concurrency := flag.Int("c", 10, "Concurrency limit (max number of simultaneous checks)")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer file.Close()

	urls, err := ReadURLs(file)
	if err != nil {
		log.Fatalf("Failed to read URLs: %v", err)
	}

	results := CheckURLs(urls, *concurrency)

	err = WriteReport("report.csv", results)
	if err != nil {
		log.Fatalf("Failed to write report: %v", err)
	}

	PrintSummary(results)
}
