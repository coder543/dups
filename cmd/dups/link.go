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

	"github.com/spf13/cobra"
)

// link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Finds duplicate files in a given path and deletes them.",
	Long: `Finds duplicate files in a given path and deletes them.
You can add '>> file.txt' at the end to export the result into a text file
`,
	Run: func(cmd *cobra.Command, args []string) {
		minSize, _ := cmd.Flags().GetInt64("min-size") // minimum file size to scan
		duplicates, _, _ := commonSetup(args, minSize)

		if len(duplicates) > 0 {
			totalSize, totalLinked, err := dups.LinkDuplicates(duplicates)
			if err != nil {
				log.Fatal("error linking duplicate files:", err)
			}
			log.Printf("converted %d files into hard links, regaining %d bytes of storage.\n", totalLinked, totalSize)
		} else {
			log.Println("no duplicate files found.")
		}
	},
}

func init() {
	addCmd(linkCmd)
}
