package ftgs

import (
	"fmt"
	"io/fs"
	"mime"
	"os"
	"path"
	"strings"

	sl "github.com/Averianov/cisystemlog"
)

var ToGo map[string][]byte

func check(err error) {
	if err != nil {
		sl.L.Warning("%s\n", err.Error())
	}
}

func init() {
	sl.CreateLogs("converter", "./log", 4, 0)
}

func ConvertDirectory(inPath, outPath, typecheck string) (err error) {
	inPath = checkPath(inPath)
	outPath = checkPath(outPath)

	var files []fs.DirEntry
	if files, err = os.ReadDir(inPath); err == nil {
		sl.L.Debug("## readDirectory: %s\n", inPath)
		for _, f := range files {
			if strings.Contains(f.Name(), typecheck) {
				fileName := f.Name()
				ConvertFile(inPath, outPath, fileName)
			}
		}
	}
	return
}

func ConvertFile(inPath, outPath, filename string) (err error) {
	var data []byte
	var contentType string
	sl.L.Debug("## ReadFileToByte - fileName: %s\n", inPath+filename)
	data, contentType, err = readFileToBytes(inPath, filename)
	sl.L.Debug("contentType: %s", contentType)
	check(err)
	if err == nil {
		err = saveByteInGo(outPath, filename, data)
		check(err)
	}
	return
}

func saveByteInGo(outPath, fileName string, data []byte) (err error) {
	var f *os.File
	path := outPath + strings.ToLower(fileName) + ".go"

	sl.L.Info("## SaveInGo - path: %v\n", path)
	//sl.L.Info("## SaveInGo - data: %v\n", data)

	f, err = os.Create(path)
	check(err)
	defer f.Close()

	_, err = f.WriteString(`package memfd

import (
	ftgc "github.com/Averianov/ftgc"
)

ftgc.ToGo[` + strings.ToUpper(strings.Replace(strings.Replace(fileName, ".", "", -1), "-", "", -1)) + `] = []byte{` + convertBytes(data) + `}`)

	return
}

func convertBytes(data []byte) string {
	var str strings.Builder
	for _, b := range data {
		fmt.Fprintf(&str, "%d,", b)
	}
	return str.String()[:str.Len()-1]
}

func readFileToBytes(inPath, file string) ([]uint8, string, error) {
	pathfile := inPath + file
	reqFile, err := os.Open(pathfile)
	if err != nil {
		return nil, "", err
	}

	var fi os.FileInfo
	fi, err = reqFile.Stat() // read file data
	contentType := mime.TypeByExtension(path.Ext(pathfile))
	var bytes = make([]uint8, fi.Size())
	reqFile.Read(bytes)
	return bytes, contentType, err
}

func checkPath(path string) string {
	if path == "" {
		path = "./"
	} else {
		runes := []rune(path)
		if runes[len(runes)-1] != '/' {
			path = path + "/"
		}
	}
	return path
}
