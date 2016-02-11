// GOGH.

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

const mainJournalFileBasename string = "GOGH.md"
const newJournalFileTempFolder string = "/tmp"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	writeAJournalEntry()
}

func writeAJournalEntry() {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	filename := filepath.Join(newJournalFileTempFolder, timestamp+".md")

	makeJournalEntryFile(filename, timestamp)
	editJournalEntryFileWithTextMate(filename)
	prependJournalEntryToMainJournal(filename)
	removeTemporaryJournalEntryFile(filename)
}

func makeJournalEntryFile(filename string, timestamp string) {
	text := []byte(timestamp + "\n\n")
	err := ioutil.WriteFile(filename, text, 0644)
	check(err)
}

func editJournalEntryFileWithTextMate(filename string) {
	// open file in textmate.
	// wait till file is closed.
	// position caret at line 3.
	cmd := exec.Command("mate", "-w", "-l 3", filename)
	err := cmd.Start()
	check(err)

	err = cmd.Wait()
	check(err)
}

func mainJournalContents() []byte {
	mainJournalContents, err := ioutil.ReadFile(mainJournalFilename())
	if err != nil {
		mainJournalContents = []byte("") // First run. Initialize main journal.
	}
	return mainJournalContents
}

func prependJournalEntryToMainJournal(journalEntryFilename string) {
	newJournalEntryContents, err := ioutil.ReadFile(journalEntryFilename)
	check(err)

	newContents := append(newJournalEntryContents, []byte("\n\n")...)
	newContents = append(newContents, mainJournalContents()...)

	f, err := os.Create(mainJournalFilename())
	defer f.Close()
	check(err)

	f.Write(newContents)
	f.Sync()
}

func removeTemporaryJournalEntryFile(filename string) {
	err := os.Remove(filename)
	check(err)
}

func mainJournalFilename() string {
	return filepath.Join(homeDirectory(), mainJournalFileBasename)
}

func homeDirectory() string {
	usr, err := user.Current()
	check(err)
	return usr.HomeDir
}
