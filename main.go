package main

import (
	"Enron_Corp_DB_Indexer/indexer"
	"bufio"
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "cpu_profile.pprof", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "mem_profile.pprof", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	fmt.Println("Starting indexer!")
	indexerData, err := indexer.CreateIndexerFromJsonFile("./index.json")
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("Deleting index if exists...")
	// deleted := indexer.DeleteIndexOnZincSearch("enron_corp")
	// if deleted != nil {
	// 	fmt.Println("Index doesn't exist. Creating...")
	// }

	sent := indexer.CreateIndexOnZincSearch(indexerData)
	if sent != nil {
		log.Fatal(sent)
	}

	log.Println("Index created successfully.")

	if len(os.Args) < 2 {
		fmt.Println("Enron Corp Directory DB Path is missing.")
		return
	}

	path := os.Args[1] + "/maildir/"

	fmt.Println("Start indexing, this might take a few minutes...")

	id := 0
	userList := indexer.GetFolders(path)
	for _, user := range userList {
		folders := indexer.GetFolders(path + user)
		for _, folder := range folders {
			emailFiles := indexer.GetFiles(path + user + "/" + folder + "/")
			for _, mail_file := range emailFiles {
				filePath := path + user + "/" + folder + "/" + mail_file
				sysFile, err := os.Open(filePath)
				if err != nil {
					fmt.Printf("Error opening file %s: %s\n", filePath, err)
					continue
				}
				lines := bufio.NewScanner(sysFile)
				id++
				indexer.PostDataToOpenObserve(indexer.FormatData(lines, id))
				sysFile.Close()
			}
		}
	}
	fmt.Println("Indexing finished!!!!")

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
