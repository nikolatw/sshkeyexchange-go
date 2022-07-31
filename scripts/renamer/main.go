package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func main() {
	binariesFolders, err := os.ReadDir("keyexchange")
	handleError(err)

	err = os.Mkdir("keyexchange/bin", 0o777)
	handleError(err)

	for _, maybeBinaryFolder := range binariesFolders {
		if maybeBinaryFolder.Type() == fs.ModeDir {
			binaryFolderName := "keyexchange/" + maybeBinaryFolder.Name()
			binaryFolder, err := os.ReadDir(binaryFolderName)
			handleError(err)

			bin, err := ioutil.ReadFile(binaryFolderName + "/" + binaryFolder[0].Name())
			handleError(err)

			extension := path.Ext(binaryFolder[0].Name())

			err = ioutil.WriteFile("keyexchange/bin/"+strings.TrimPrefix(maybeBinaryFolder.Name(), "sshkeyexchange-go_")+extension, bin, 0644)
			handleError(err)
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
