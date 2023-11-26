package main

import (
	"fmt"
	"log"
	"os"

	coordinationservice "scale.kv.store/pkg/NameNodeService"
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

	coordinator := coordinationservice.NewCoordinator()

	fmt.Println(coordinator.IndexRing.Ring[0])

	fmt.Println(coordinator.Put(5, 1).ToString())

	fmt.Println(coordinator.Put(5, 2).ToString())

	fmt.Println("Get 5 - 2")

	fmt.Println(coordinator.Get(5).ToString())

	fmt.Println(coordinator.Put(5, 3).ToString())

	fmt.Println("Get 5 - 3")

	fmt.Println(coordinator.Get(5).ToString())

	fmt.Println(coordinator.Put(5, 4).ToString())

	fmt.Println(coordinator.Put(6, 4).ToString())

	fmt.Println("Delete 5")

	fmt.Println(coordinator.Delete(5).ToString())

	fmt.Println("Get 5")

	fmt.Println(coordinator.Get(5).ToString())

	// fmt.Println("Buckets")

	// fmt.Println(coordinator.Buckets[0])
}
