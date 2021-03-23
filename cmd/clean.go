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
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Finds duplicate files in a given path and deletes them.",
	Long: `Finds duplicate files in a given path and deletes them.
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
		singleCore, _ := cmd.Flags().GetBool("single-core") // single core option
		minSize, _ := cmd.Flags().GetInt64("min-size")      // minimum file size to scan
		flat, _ := cmd.Flags().GetBool("flat")
		if !flat {
			fmt.Println("scanning path ...")
		}
		files, err := dups.GetFiles(path, minSize)
		if err != nil {
			log.Fatal("error while listing files:", err)
		}
		if !flat {
			fmt.Printf("found %d files. calculating hashes using sha256 algorithm with multicore: %t", len(files), !singleCore)
		}
		groups, totalFiles := dups.GroupFiles(files)
		hashes := dups.CollectHashes(groups, singleCore, flat, totalFiles)
		if !flat {
			fmt.Println("scanning for duplicates ...")
		}
		duplicates, totalFiles, totalDuplicates := dups.GetDuplicates(hashes)
		if !flat {
			fmt.Printf("found %d files with total of %d duplicates\n", totalFiles, totalDuplicates)
			for _, fs := range duplicates {
				fmt.Printf("Path: %s \nSize: %d\n", fs[0].Path, fs[0].Size)
				for _, file := range fs[1:] {
					fmt.Println(file.Path)
				}
				fmt.Println("============================================================================")
			}
		} else {
			for _, fs := range duplicates {
				for _, file := range fs[1:] {
					fmt.Println(file.Path)
				}
			}
		}
		if len(duplicates) > 0 {
			totalSize, totalDeleted, err := dups.RemoveDuplicates(duplicates)
			if err != nil {
				log.Fatal("error deleting duplicate files:", err)
			}
			fmt.Printf("removed %d files with the total size of %d bytes.\n", totalDeleted, totalSize)
		} else {
			fmt.Println("no duplicate files found.")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolP("flat", "f", false, "flat output, no extra info (only prints duplicate files")
	cleanCmd.Flags().BoolP("single-core", "s", false, "use single cpu core")
	cleanCmd.Flags().Int64("min-size", 1024, "minimum file size to scan in bytes")
}
