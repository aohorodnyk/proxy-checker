package main

import (
	"bufio"
	"log"
	"os"
	"sync"
)

func readFileList(fileName string) (lines []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// This is our buffer now
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func writeFileListChannel(out chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.OpenFile(configuration.FileNameResult, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	for proxyURL := range out {
		_, err = file.WriteString(proxyURL + "\n")
		if err != nil {
			log.Fatalln(err)
		}
	}
}
