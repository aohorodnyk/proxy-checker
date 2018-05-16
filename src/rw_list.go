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
		lines = addLine(lines, line)
	}

	return lines, nil
}

func writeFileListChannel(out chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.OpenFile(configuration.Result, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
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

func addLine(lines []string, line string) []string {
	if line == "" {
		return lines
	}
	for _, value := range lines {
		if value == line {
			return lines
		}
	}
	return append(lines, line)
}
