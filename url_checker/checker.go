package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type URLResult struct {
	URL          string
	Status       string
	ResponseTime string
	Error        string
}

func ReadURLs(file *os.File) ([]string, error) {
	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func CheckURLs(urls []string, concurrency int) []URLResult {
	results := make([]URLResult, len(urls))
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, concurrency)
	progress := make(chan int)

	total := len(urls)

	go func() {
		checked := 0
		for range progress {
			checked++
			fmt.Printf("\rProgress: %d/%d", checked, total)
		}
	}()

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			start := time.Now()
			resp, err := http.Get(url)
			duration := time.Since(start).Milliseconds()

			result := URLResult{URL: url}

			if err != nil {
				result.Status = "N/A"
				result.ResponseTime = "N/A"
				result.Error = err.Error()
				log.Printf("Error checking %s: %v\n", url, err)
			} else {
				defer resp.Body.Close()
				result.Status = fmt.Sprintf("%d", resp.StatusCode)
				result.ResponseTime = fmt.Sprintf("%d", duration)
			}

			results[i] = result
			progress <- 1
		}(i, url)
	}

	wg.Wait()
	close(progress)

	fmt.Println()
	return results
}

func WriteReport(filename string, results []URLResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"URL", "Status", "ResponseTime(ms)", "Error"})

	for _, r := range results {
		writer.Write([]string{r.URL, r.Status, r.ResponseTime, r.Error})
	}
	return nil
}

func PrintSummary(results []URLResult) {
	var failed, total int
	var totalTime int64

	for _, r := range results {
		total++
		if r.Error != "" {
			failed++
		} else {
			var ms int64
			fmt.Sscanf(r.ResponseTime, "%d", &ms)
			totalTime += ms
		}
	}

	fmt.Println("\n=== Summary ===")
	fmt.Printf("Total URLs: %d\n", total)
	fmt.Printf("Failed: %d\n", failed)
	if total-failed > 0 {
		avg := totalTime / int64(total-failed)
		fmt.Printf("Average Response Time: %d ms\n", avg)
	}
}
