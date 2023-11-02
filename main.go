package main

import (
	"fmt"
	"log"
	"os"
	"scale.kv.store/internal/models"
	"strconv"
)

func ConfigureLogger() *os.File {

	logfile, logFErr := os.Create("./storage/logs/log.log")

	if logFErr != nil {
		log.Fatalf("Error opening log file: %v", logFErr)
	}

	log.SetOutput(logfile)
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime) // Todo: add it later | log.LUTC

	return logfile
}

func main() {

	logfile := ConfigureLogger()
	defer logfile.Close()

	shard := models.NewShard()

	fmt.Println("shard version: " + strconv.FormatUint(uint64(shard.Version.Get()), 10))

	fmt.Println(shard.Put(5, 5).Value)

	fmt.Println(shard.Get(5).Value)

	fmt.Println(shard.Delete(5))

	fmt.Println(shard.Get(5))

	fmt.Println(shard.Buckets[0])
}
