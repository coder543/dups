package main

import (
	"dups"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func addCmd(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
	cmd.Flags().Int64("min-size", 1024, "minimum file size to scan in bytes")
}

func commonSetup(args []string, minSize int64) ([][]dups.FileInfo, int64, int64) {
	if len(args) == 0 {
		log.Fatal("please provide a path: dups find path/to/directory")
	}
	path := dups.CleanPath(args[0])
	f, err := os.Stat(path)
	if err != nil {
		log.Fatal("can't find path:", err)
	}
	if !f.IsDir() {
		log.Fatal("please provide a directory path not a file path")
	}

	log.Println("scanning path for file metadata...")
	files, err := dups.GetFiles(path, minSize)
	if err != nil {
		log.Fatal("error while listing files:", err)
	}
	log.Printf("found %d files of min-size or larger\n", len(files))

	groups := dups.GroupFiles(files)

	totalFiles := int64(0)
	for _, group := range groups {
		totalFiles += int64(len(group))
	}
	log.Printf("found %d files for which there are one or more other files of the same size\n", totalFiles)

	hashes := dups.CollectHashes(groups)

	log.Println("scanning for duplicates ...")
	duplicates, totalFiles, totalDuplicates := dups.GetDuplicates(hashes)
	log.Printf("found %d files with total of %d duplicates\n", totalFiles, totalDuplicates)

	for _, fs := range duplicates {
		log.Printf("Path: %s \nSize: %d\n", fs[0].Path, fs[0].Size)
		for _, file := range fs[1:] {
			fmt.Println(file.Path)
		}
		log.Println("============================================================================")
	}

	return duplicates, totalFiles, totalDuplicates
}
