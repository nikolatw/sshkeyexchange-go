package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	binariesFolders, err := os.ReadDir("keyexchange")
	handleError(err)
	for _, maybeBinaryFolder := range binariesFolders {
		if maybeBinaryFolder.Type() == fs.ModeDir {
			binaryFolderName := "keyexchange/" + maybeBinaryFolder.Name()
			binaryFolder, err := os.ReadDir(binaryFolderName)
			handleError(err)

			bin, err := ioutil.ReadFile(binaryFolderName + "/" + binaryFolder[0].Name())
			handleError(err)

			err = ioutil.WriteFile("keyexchange/"+strings.TrimPrefix(maybeBinaryFolder.Name(), "sshkeyexchange-go_"), bin, 0644)
			handleError(err)
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
