package makefilelist

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

/*
	FilterExt returns a string slice with all files that don't match the provided extension
	filtered out.
*/
func FilterExt(exts []string, fileNames []string) []string {
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
func WriteList(fileName string, namesToWrite []string) error {

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
func TraverseFolder(folder string, runRecursive bool) []string {

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
				nextDir := TraverseFolder(dir, true)
				fileNames = append(fileNames, nextDir...)
			}
		}
		return fileNames
	}

	return fileNames
}
