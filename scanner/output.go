package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"minecraft_searcher/redis"
)

// func (s scanner) outputFileWriter(outputFile string, overwrite bool) chan<- string {
// 	input := make(chan string, s.threads)
// 	go func() {
// 		if overwrite {
// 			os.Remove(outputFile)
// 		}
// 		file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 		if err != nil {
// 			s.errrorsLogger <- err
// 		}
// 		defer file.Close()
// 		for {
// 			dataToWrite := <-input
// 			file.WriteString(fmt.Sprintf("%s\n", dataToWrite))
// 		}
// 	}()
// 	return input
// }

func (s scanner) workerOutputParaser() chan<- WorkerOutputData {
	input := make(chan WorkerOutputData, s.threads)
	go func() {
		for {
			data := <-input
			jsonData, err := json.Marshal(data)
			if err != nil {
				s.errrorsLogger <- fmt.Errorf("error when decoding json in function workerOutputParaser: %s", err)
			}
			err = redis.GetRedis().LPush(context.Background(), "servers", jsonData).Err()
			if err != nil {
				s.errrorsLogger <- fmt.Errorf("error when putting data to redis in function workerOutputParaser: %s", err)
			}
			// _ := string(bytes)
		}
	}()
	return input
}
