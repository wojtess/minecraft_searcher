package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type JsonEntry struct {
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
	Ports     []struct {
		Port   int    `json:"port"`
		Proto  string `json:"proto"`
		Status string `json:"status"`
		Reason string `json:"reason"`
		TTL    int    `json:"ttl"`
	} `json:"ports"`
}

func (s scanner) inputFileReaderJson(inputFile *os.File) <-chan WorkerInputData {
	buf, err := ioutil.ReadAll(inputFile)

	if err != nil {
		s.errrorsLogger <- fmt.Errorf("error when reading file in inputFileReaderJson: %s", err)
		return nil
	}

	var data []JsonEntry

	err = json.Unmarshal(buf, &data)

	if err != nil {
		s.errrorsLogger <- fmt.Errorf("error when creating json in inputFileReaderJson: %s", err)
		return nil
	}
	size := len(data) - 1
	output := make(chan WorkerInputData, s.threads)
	go func() {
		defer fmt.Printf("producer stoped\n")
		i := int(0)
		for {
			if i >= len(data) {
				break
			}
			output <- WorkerInputData{
				Ip:       data[i].IP,
				Ordinary: i,
			}
			if i%1500 == 0 {
				fmt.Printf("[%d/%d]putting %s\n", i, size, data[i].IP)
			}
			i++
		}
		for i := 0; i < s.threads; i++ {
			output <- WorkerInputData{
				Ordinary: -1,
			}
		}
	}()
	return output
}

func (s scanner) inputFileReader(inputFile *os.File) <-chan WorkerInputData {
	reader := bufio.NewReader(inputFile)
	reader.Discard(s.ordinary * 4)
	info, _ := inputFile.Stat()
	size := info.Size() / 4
	output := make(chan WorkerInputData, s.threads)
	go func() {
		i := int(0)
		for {
			buffer := make([]byte, 4)
			readed, err := reader.Read(buffer)
			if err != nil {
				s.errrorsLogger <- fmt.Errorf("inputFileReader failed: %s", err)
			}
			if readed != len(buffer) {
				s.errrorsLogger <- fmt.Errorf("inputFileReader failed: reader bytes length is diffrent than buffer size")
			}
			ip := fmt.Sprintf("%d.%d.%d.%d", buffer[0], buffer[1], buffer[2], buffer[3])
			output <- WorkerInputData{
				Ip:       ip,
				Ordinary: i,
			}
			if i%1500 == 0 {
				fmt.Printf("[%d/%d]putting %s\n", i, size, ip)
			}
			i++
		}
	}()
	return output
}
