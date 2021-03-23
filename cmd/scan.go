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
	"bufio"
	"dups"
	"fmt"
	"log"
	"os"
	"strings"

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
		flat, _ := cmd.Flags().GetBool("flat")

		if !flat {
			fmt.Println("scanning path ...")
		}
		files, err := dups.GetFiles(path, minSize)
		if err != nil {
			log.Fatal("error while listing files:", err)
		}
		if !flat {
			fmt.Printf("found %d files. calculating hashes using sha256 algorithm with multicore: %t\n", len(files), !singleCore)
		}
		groups, totalFiles := dups.GroupFiles(files)
		hashes := dups.CollectHashes(groups, singleCore, flat, totalFiles)
		if !flat {
			fmt.Println("scanning for duplicates ...")
		}
		duplicates, totalFiles, totalDuplicates := dups.GetDuplicates(hashes)
		if !flat {
			totalBytes := int64(0)
			for _, fs := range duplicates {
				fmt.Printf("Path: %s \nSize: %d\n", fs[0].Path, fs[0].Size)
				for _, file := range fs[1:] {
					fmt.Println(file.Path)
					totalBytes += fs[0].Size
				}
				fmt.Println("============================================================================")
			}
			fmt.Printf("found %d files with total of %d duplicates wasting %d bytes\n", totalFiles, totalDuplicates, totalBytes)
		} else {
			for _, fs := range duplicates {
				for _, file := range fs[1:] {
					fmt.Println(file.Path)
				}
			}
		}
		if !flat {
			if len(duplicates) > 0 {
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Println("Listing completed.")
				fmt.Println("Would you like to delete duplicates? (y/n)")
				scanner.Scan()
				text := scanner.Text()
				lowered := strings.ToLower(text)
				if lowered == "y" || lowered == "yes" {
					totalSize, totalDeleted, err := dups.RemoveDuplicates(duplicates)
					if err != nil {
						log.Fatal("error deleting duplicate files:", err)
					}
					fmt.Printf("removed %d files with the total size of %d bytes\n", totalDeleted, totalSize)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().BoolP("flat", "f", false, "flat output, no extra info (only prints duplicate files")
	scanCmd.Flags().BoolP("single-core", "s", false, "use single cpu core")
	scanCmd.Flags().Int64("min-size", 1024, "minimum file size to scan in bytes")
}
