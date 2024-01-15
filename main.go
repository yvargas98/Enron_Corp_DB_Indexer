package main

import (
	"Enron_Corp_DB_Indexer/indexer"
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "cpu_profile.pprof", "write cpu profile to cpu_profile.pprof")
var memprofile = flag.String("memprofile", "mem_profile.pprof", "write memory profile to mem_profile.pprof")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

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
				data, err := indexer.ProcessFile(filePath, id)
				if err != nil {
					fmt.Printf("Error processing file %s: %s\n", filePath, err)
					continue
				}
				id++
				indexer.PostDataToOpenObserve(data)
			}
		}
	}
	fmt.Println("Indexing finished!!!!")

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
