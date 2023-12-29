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

	if len(os.Args) < 2 {
		fmt.Println("Enron Corp Directory DB Path is missing.")
		return
	}

	path := os.Args[1] + "/maildir/"

	fmt.Println("Indexing started...")

	id := 0
	userList := indexer.GetFolders(path)
	for _, user := range userList {
		folders := indexer.GetFolders(path + user)
		for _, folder := range folders {
			emailFiles := indexer.GetFiles(path + user + "/" + folder + "/")
			for _, mail_file := range emailFiles {
				sysFile, _ := os.Open(path + user + "/" + folder + "/" + mail_file)
				lines := bufio.NewScanner(sysFile)
				id++
				indexer.PostDataToZincSearch(indexer.FormatData(lines, id))
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
