package dups

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

// FileInfo represents a file containing os.FileInfo and file path
type FileInfo struct {
	Path string
	Size int64
}

type barWriter struct {
	bar   *pb.ProgressBar
	inner io.Writer
}

func (w *barWriter) Write(p []byte) (n int, err error) {
	n, err = w.inner.Write(p)
	w.bar.Add(n)
	return n, err
}

// GetFileHash returns sha256 hash of a file
func GetFileHash(path string, bar *pb.ProgressBar) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()
	h := sha256.New()

	bw := &barWriter{
		bar:   bar,
		inner: h,
	}

	limited := io.LimitReader(f, 128*1024)

	if _, err := io.Copy(bw, limited); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// GetFiles finds and returns all the files in the given path
// It will also returns any file in sub-directories if "full=true"
func GetFiles(root string, minSize int64) ([]FileInfo, error) {
	var filesInfos []FileInfo
	cleanedPath := CleanPath(root)

	err := filepath.Walk(cleanedPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// Ignore files less than minimum size
		size := info.Size()
		if size < minSize {
			return nil
		}

		filesInfos = append(filesInfos, FileInfo{
			Path: path,
			Size: size,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return filesInfos, nil
}

// GroupFiles groups files based on their file size
// This will help avoid unnecessary hash calculations since files with different file sizes can't be duplicates
func GroupFiles(files []FileInfo) (map[int64][]FileInfo, int64) {
	groups := make(map[int64][]FileInfo)
	fileCount := int64(0)
	for _, file := range files {
		groups[file.Size] = append(groups[file.Size], file)
		fileCount++
	}
	for bucket, files := range groups {
		numFiles := len(files)
		if numFiles < 2 {
			fileCount -= int64(numFiles)
			delete(groups, bucket)
		}
	}
	return groups, fileCount
}

// CollectHashes returns hashes for the given group files if there is more than one file with the same size
// A hash will be the key and a list of FileInfo for files that share the hash as the value
// "singleThread=false" will force all the function to use one thread only
// minSize is the minimum file size to scan
// "flat=true" will tell the function not to print out any data other than the path to duplicate files
// algorithm is the algorithm to calculate the hash with
func CollectHashes(fileGroups map[int64][]FileInfo, singleThread bool, fileCount int64) map[string][]FileInfo {
	hashes := map[string][]FileInfo{}
	var lock = sync.Mutex{}

	// TODO: create second bar when initial duplicates are found,
	// then do full file hashes to compare those duplicates
	bar := createBar(128 * 1024 * fileCount)
	defer bar.Finish()

	if singleThread {
		// All groups will have more than one file
		for _, group := range fileGroups {
			for _, file := range group {
				hash, err := GetFileHash(file.Path, bar)
				if err != nil {
					log.Printf("Encountered error hashing file %q: %s", file.Path, err)
				} else {
					hashes[hash] = append(hashes[hash], file)
				}
			}
		}

		return hashes
	}

	numWorkers := runtime.GOMAXPROCS(0)
	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)

	workChan := make(chan FileInfo, numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for file := range workChan {
				hash, err := GetFileHash(file.Path, bar)
				if err != nil {
					log.Printf("Encountered error hashing file %q: %s", file.Path, err)
				} else {
					lock.Lock()
					hashes[hash] = append(hashes[hash], file)
					lock.Unlock()
				}
			}

			wg.Done()
		}()
	}

	for _, group := range fileGroups {
		// All groups will have more than one file
		for _, file := range group {
			file := file
			workChan <- file
		}
	}

	close(workChan)
	wg.Wait()

	return hashes
}

// GetDuplicates scans the given map of hashes and finds the one with duplicates
// It will return a slice containing slices with each slice containing paths to duplicate files
// It will also returns the total of duplicate files and the total of files that have duplicates
func GetDuplicates(hashes map[string][]FileInfo) ([][]FileInfo, int64, int64) {
	var duplicateFiles [][]FileInfo
	// total duplicate files
	total := int64(0)
	// Total number of files with duplicates
	totalFiles := int64(0)
	for _, files := range hashes {
		if len(files) > 1 {
			totalFiles++
			// for original file which will be counted in the next for
			total--
			var duplicates []FileInfo
			for _, file := range files {
				total++
				duplicates = append(duplicates, file)
			}
			duplicateFiles = append(duplicateFiles, duplicates)
		}
	}
	return duplicateFiles, totalFiles, total
}

// RemoveDuplicates removes duplicates
// It will keep the first file in a duplicate set and removes any other files in the set
// It will return the sum of deleted file sizes and total number of deleted files
func RemoveDuplicates(fileSets [][]FileInfo) (int64, int64, error) {
	totalSize := int64(0)
	totalDeleted := int64(0)
	for _, files := range fileSets {
		for i, file := range files {
			if i > 0 {
				totalSize += file.Size
				totalDeleted++
				err := os.Remove(file.Path)
				if err != nil {
					return totalSize, totalDeleted, err
				}

			}
		}
	}
	return totalSize, totalDeleted, nil
}

// LinkDuplicates links duplicates
// It will keep the first file in a duplicate set and make hard links from any other files in the set to that file.
// It will return the sum of linked file sizes and total number of linked files.
func LinkDuplicates(fileSets [][]FileInfo) (int64, int64, error) {
	totalSize := int64(0)
	totalLinked := int64(0)
	for _, files := range fileSets {
		linkPath := files[0].Path
		for _, file := range files[1:] {
			totalSize += file.Size
			totalLinked++

			err := os.Remove(file.Path)
			if err != nil {
				return totalSize, totalLinked, err
			}

			err = os.Link(linkPath, file.Path)
			if err != nil {
				return totalSize, totalLinked, err
			}
		}
	}
	return totalSize, totalLinked, nil
}
