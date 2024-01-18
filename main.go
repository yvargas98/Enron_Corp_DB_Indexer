package main

import (
	"Enron_Corp_DB_Indexer/indexer"
	"flag"
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
		log.Println("Enron Corp Directory DB Path is missing.")
		return
	}

	path := os.Args[1] + "/maildir/"

	log.Println("Start indexing, this might take a few minutes...")

	id := 0
	batchSize := 100
	var dataBatch []indexer.ECEmail
	userList, err := indexer.GetFolders(path)
	if err != nil {
		log.Println("Error: ", err)
	}
	for _, user := range userList {
		folders, err := indexer.GetFolders(path + user)
		if err != nil {
			log.Println("Error: ", err)
		}
		for _, folder := range folders {
			emailFiles, err := indexer.GetFiles(path + user + "/" + folder + "/")
			if err != nil {
				log.Println("Error: ", err)
			}
			for _, mail_file := range emailFiles {
				filePath := path + user + "/" + folder + "/" + mail_file
				data, err := indexer.ProcessFile(filePath, id)
				if err != nil {
					log.Printf("Error processing file %s: %s\n", filePath, err)
					continue
				}

				dataBatch = append(dataBatch, data)

				// Enviar un lote de 10 archivos
				if len(dataBatch) == batchSize {
					err := indexer.PostDataToOpenObserve(dataBatch)
					if err != nil {
						panic(err)
					}
					dataBatch = nil // Limpiar el lote después de enviar
				}
				id++
			}
		}
	}

	// Enviar los datos restantes, si los hay
	if len(dataBatch) > 0 {
		indexer.PostDataToOpenObserve(dataBatch)
	}
	log.Println("Indexing finished!!!!")

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
