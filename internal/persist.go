package internal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/printzero/tint"
)

// FocusData is config required for focus to function
type FocusData struct {
	Editor string
	// map of repo url and array of todos
	TODOs          map[string][]string
	TODODependants map[int][]int
}

// SaveTODO saves a todo for the current repo to disk
func SaveTODO(fd FocusData, todo string) error {
	base, err := GetRepositoryURL()
	if err != nil {
		return nil
	}
	if fd.TODOs == nil && len(fd.TODOs[base]) == 0 {
		fd.TODOs = map[string][]string{
			base: {todo},
		}
	} else {
		fd.TODOs[base] = append(fd.TODOs[base], todo)
	}

	var buf bytes.Buffer
	encodeData(fd, &buf)
	err = writeFocusData(buf)
	if err != nil {
		return err
	}
	return nil
}

// RemoveTODO removes a TODO from the list of todos based on its index
func RemoveTODO(fd FocusData, atIndex int) error {
	var modTODOs []string
	base, err := GetRepositoryURL()
	if err != nil {
		return err
	}

	for i, todo := range fd.TODOs[base] {
		if atIndex != i+1 {
			modTODOs = append(modTODOs, todo)
		}
	}
	fd.TODOs[base] = modTODOs

	var buf bytes.Buffer
	encodeData(fd, &buf)
	err = writeFocusData(buf)
	if err != nil {
		return err
	}
	return nil
}

// ListTODOs displays a list of todos
func ListTODOs(fd FocusData) {
	base, _ := GetRepositoryURL()
	if len(fd.TODOs) == 0 {
		fmt.Println("no todos to list!")
		return
	}

	var result string
	for i, todo := range fd.TODOs[base] {
		result += tint.Init().Exp(fmt.Sprintf("@(%d.) %s\n", i+1, todo), tint.Green.Bold())
	}
	fmt.Println(result)
}

func encodeData(fd FocusData, buf *bytes.Buffer) error {
	enc := gob.NewEncoder(buf)
	err := enc.Encode(fd)
	if err != nil {
		return err
	}
	return nil
}

func decodeData() (FocusData, error) {
	var fd FocusData
	// TODO: get file from disk here
	b, err := ioutil.ReadFile("")
	if err != nil {
		return fd, err
	}

	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err = dec.Decode(&fd)
	if err != nil {
		return fd, err
	}
	return fd, nil
}

func getFocusDataFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".focus")
}

func writeFocusData(buf bytes.Buffer) error {
	cachePath := getFocusDataFilePath()
	err := ioutil.WriteFile(cachePath, buf.Bytes(), 0755)
	if err != nil {
		return err
	}
	return nil
}

// GetFocusData returns a focus cache file which is located inside
// your home dir if no focus config is found it creates a
// default config in this path and returns it
func GetFocusData() (FocusData, error) {
	var fd FocusData
	cachePath := getFocusDataFilePath()
	if f, _ := os.Stat(cachePath); f == nil {
		var editor string
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "nano"
		}

		fd := FocusData{
			Editor: editor,
		}

		var buf bytes.Buffer
		encodeData(fd, &buf)
		err := ioutil.WriteFile(cachePath, buf.Bytes(), 0755)
		if err != nil {
			return fd, err
		}
		return fd, nil
	}

	b, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return fd, err
	}

	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err = dec.Decode(&fd)
	if err != nil {
		return fd, err
	}
	return fd, nil
}
