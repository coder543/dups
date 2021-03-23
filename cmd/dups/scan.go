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
	"log"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Finds duplicate files in a given path but doesn't delete them.",
	Long: `Finds duplicate files in a given path but doesn't delete them.
You can add '>> file.txt' at the end to export the result into a text file
`,
	Run: func(cmd *cobra.Command, args []string) {
		minSize, _ := cmd.Flags().GetInt64("min-size")
		duplicates, totalFiles, totalDuplicates := commonSetup(args, minSize)

		totalBytes := int64(0)
		for _, fs := range duplicates {
			for _, file := range fs[1:] {
				totalBytes += file.Size
			}
		}
		log.Printf("found %d files with total of %d duplicates wasting %d bytes\n", totalFiles, totalDuplicates, totalBytes)
	},
}

func init() {
	addCmd(scanCmd)
}
