/*
Copyright © 2020 NAME HERE navid2zp@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"dups"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Finds duplicate files in a given path but doesn't delete them.",
	Long: `Finds duplicate files in a given path but doesn't delete them.
You can add '>> file.txt' at the end to export the result into a text file
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("please provide a path: dups find path/to/directory")
			return
		}
		path := dups.CleanPath(args[0])
		f, err := os.Stat(path)
		if err != nil {
			log.Fatal("can't find path:", err)
		}
		if !f.IsDir() {
			log.Fatal("please provide a directory path not a file path")
		}
		singleCore, _ := cmd.Flags().GetBool("single-core")
		minSize, _ := cmd.Flags().GetInt64("min-size")

		log.Println("scanning path ...")
		files, totalSize, err := dups.GetFiles(path, minSize)
		if err != nil {
			log.Fatal("error while listing files:", err)
		}

		log.Printf("found %d files. calculating hashes using sha256 algorithm with multicore: %t\n", len(files), !singleCore)
		groups, totalFiles := dups.GroupFiles(files)
		hashes := dups.CollectHashes(groups, singleCore, totalFiles, totalSize)

		log.Println("scanning for duplicates ...")
		duplicates, totalFiles, totalDuplicates := dups.GetDuplicates(hashes)
		for _, fs := range duplicates {
			log.Printf("Path: %s \nSize: %d\n", fs[0].Path, fs[0].Size)
			for _, file := range fs[1:] {
				log.Println(file.Path)
			}
			log.Println("============================================================================")
		}
		log.Printf("found %d files with total of %d duplicates wasting %d bytes\n", totalFiles, totalDuplicates, totalSize)
	},
}

func init() {
	addCmd(scanCmd)
}
