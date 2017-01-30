package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/DHowett/go-plist"
	"github.com/kyokomi/archivex"
)

// Packager alfred workflow packager
type Packager struct {
	ExportFile  string `json:"exportFile"`
	Readme      string `json:"readme"`
	ReplaceInfo map[string]interface{}
	Resources   []string `json:"resources"`
}

func readPackager(packagerFilePath string) (Packager, error) {
	p := Packager{}
	data, err := ioutil.ReadFile(packagerFilePath)
	if err != nil {
		return Packager{}, err
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return Packager{}, err
	}
	return p, nil
}

func readInfoPlistWithReplace(p Packager, infoFilePath string) ([]byte, error) {
	infoMap, format, err := readInfoPlist(infoFilePath)
	if err != nil {
		return nil, err
	}

	// Replace info.plist
	for key, value := range p.ReplaceInfo {
		infoMap[key] = value
	}

	// replace readme
	readmeData, err := ioutil.ReadFile(p.Readme)
	if err != nil {
		return nil, err
	}
	infoMap["readme"] = string(readmeData)

	return plist.Marshal(infoMap, format)
}

func readInfoPlist(filePath string) (info map[string]interface{}, format int, err error) {
	var data []byte
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	format, err = plist.Unmarshal(data, &info)
	return info, format, err
}

func archiveAlfredWorkflow(p Packager, infoFileData []byte) {
	zip := new(archivex.ZipFile)
	zip.Create(p.ExportFile)
	zip.Add("info.plist", infoFileData)
	for _, resource := range p.Resources {
		zip.AddFile(resource)
	}
	zip.Close()
}

func main() {
	packagerFilePath := flag.String("p", "packager.json", "packager.json file path")
	infoFilePath := flag.String("i", "info.plist", "info.plist file path")
	flag.Parse()

	p, err := readPackager(*packagerFilePath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	data, err := readInfoPlistWithReplace(p, *infoFilePath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	archiveAlfredWorkflow(p, data)
}
