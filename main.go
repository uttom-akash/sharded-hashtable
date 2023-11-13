package main

import (
	"fmt"
	"log"
	"os"

	"scale.kv.store/internal/models"
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

	coordinator := models.NewCoordinator()

	fmt.Println(coordinator.IndexRing[0])

	fmt.Println(coordinator.Put(5, 5))

	fmt.Println(coordinator.Put(5, 5))

	fmt.Println(coordinator.Put(5, 5))

	fmt.Println(coordinator.Put(5, 5))

	fmt.Println("Get 5")

	fmt.Println(coordinator.Get(5).Value)

	// fmt.Println("Delete 5")

	// fmt.Println(coordinator.Delete(5))

	fmt.Println("Get 5")

	fmt.Println(coordinator.Get(5))

	// fmt.Println("Buckets")

	// fmt.Println(coordinator.Buckets[0])
}
