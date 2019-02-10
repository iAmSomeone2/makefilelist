// makefilelist is designed to make a text file containing a list of all of the
// files in a directory.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// String constants for the command-line flags
const fileFlagDesc string = "The destination of the output file."
const recursiveFlagDesc string = "Set this flag to traverse all files in the directory."
const extFlagDesc string = "Only lists the files with the specified file extensions.\nMultiple extensions can be separated with ','"

func main() {
	// All flag pointers must be dereferenced to be used
	directoryPtr := flag.String("dir", "./", "The directory to make the list for.")
	recursiveTruePtr := flag.Bool("recursive", false, recursiveFlagDesc)
	fileNamePtr := flag.String("out", "./file_list.txt", fileFlagDesc)
	extPtr := flag.String("ext", "", extFlagDesc) // Read through this flag to get comma separated values

	flag.Parse()

	fmt.Println("\nReading from: " + *directoryPtr + "\n")

	fmt.Println("Working...")

	foundFiles := traverseFolder(*directoryPtr, *recursiveTruePtr)
	// Get all of the files in a slice

	if *extPtr != "" {
		// If the ext flag is set, filer out everything that doesn't match it
		// Run filter function
		exts := strings.Split(*extPtr, ",")
		fmt.Println("Extensions to look for:")
		fmt.Println(exts)
		foundFiles = filterExt(exts, foundFiles)
	}

	err := writeList(*fileNamePtr, foundFiles)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nAll found files are listed in: " + *fileNamePtr)
}

/*
	Returns a string slice with all files that don't match the provided extension
	filtered out.
*/
func filterExt(exts []string, fileNames []string) []string {
	var filteredFileNames []string

	for _, fileName := range fileNames {
		/*
			Read through all of the file names and remove any that don't
			 have the provided extension.
		*/
		for _, ext := range exts {
			if path.Ext(fileName) == ext {
				filteredFileNames = append(filteredFileNames, fileName)
			}
		}
	}

	return filteredFileNames
}

/*
	Creates the output file and writes the results to it.
*/
func writeList(fileName string, namesToWrite []string) error {

	/*
		Check if file exists before creation.
		If it does, delete the original file and create a new one.
	*/

	if _, err := os.Stat(fileName); err == nil {
		// The file exists and there were no errors.
		fmt.Println("\n" + fileName + " was found. Replacing with a new version...")
		os.Remove(fileName)
	} else {
		// The system can't determine if the file exists or not
		fmt.Println("\n" + fileName + " couldn't be found. Creating...")
		// return err
	}

	outputFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer outputFile.Close() // This will close the file after the function is finished.

	// Create a formatted string for the list then push it to the file
	numNames := len(namesToWrite)
	var outputString string
	for index, name := range namesToWrite {
		if index == numNames-1 {
			outputString += name
		} else {
			outputString += name + "\n"
		}
	}

	_, err = outputFile.WriteString(outputString)
	if err != nil {
		return err
	}

	return nil
}

/*
	This function returns a slice with all of the file names in the given folder.
	If runRecursive is true, the function will recursively read through any folders
	found in the given directory.
*/
func traverseFolder(folder string, runRecursive bool) []string {

	const pathSeparator string = string(os.PathSeparator)

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println("ERROR: Couldn't enter: " + folder)
		log.Fatal(err)
	}

	var directories [5000]string
	var realFiles [10000]string
	fileIdx := 0
	dirIdx := 0
	// Essentially a for-each loop
	for _, file := range files {
		if !file.IsDir() {
			realFiles[fileIdx] = path.Clean(folder + pathSeparator + file.Name())
			fileIdx++
		} else {
			directories[dirIdx] = path.Clean(folder + pathSeparator + file.Name())
			dirIdx++
		}
	}

	// Stop going if there are no more directories.
	if dirIdx == 0 {
		runRecursive = false
	}

	fileNames := realFiles[0:fileIdx]
	// The recursive flag was set, so parse all directories inside the provided one.
	if runRecursive {
		// Each directory that was found needs to be traversed.
		for _, dir := range directories {
			if dir != "" {
				nextDir := traverseFolder(dir, true)
				fileNames = append(fileNames, nextDir...)
			}
		}
		return fileNames
	}

	return fileNames
}
