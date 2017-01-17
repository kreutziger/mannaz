package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

const VERSION string = "0.1.5"

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(s string) error {
	sf.value = s
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

var fileVar stringFlag
var stdioVar bool
var versionVar bool
var typeVar bool
var outputFileVar stringFlag

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkVersion(version bool) {
	if version {
		fmt.Println(VERSION)
		os.Exit(0)
	}
}

func readFile(file string) []byte {
	if file == "" {
		fmt.Println("no file was given")
		os.Exit(66)
	}

	dat, err := ioutil.ReadFile(file)
	check(err)
	return dat
}

func writeFile(file string, conv []byte, stdio bool) {
	if stdio {
		fmt.Println(string(conv))
	}

	if file != "" {
		err := ioutil.WriteFile(file, conv, 0644)
		check(err)
	}
}

func toJSON(dat []byte, ext string) interface{} {
	err := error(nil)
	var packer interface{}

	switch ext {
	case "yml", "YML", "yaml", "YAML":
		dat, err = yaml.YAMLToJSON(dat)
		check(err)
		fallthrough
	case "json", "JSON":
		fallthrough
	case "":
		err = json.Unmarshal(dat, &packer)
		check(err)
	default:
		fmt.Printf("file has no valid filetype: %s\n", ext)
		os.Exit(22)
	}

	return packer
}

func fmtJSON(manifest interface{}, toYaml bool) []byte {
	json, err := json.MarshalIndent(manifest, "", "  ")
	check(err)

	if toYaml {
		json, err = yaml.JSONToYAML(json)
		check(err)
	}

	return json
}

func init() {
	flag.Var(&fileVar, "file", "file to convert")
	flag.BoolVar(&stdioVar, "stdio", false, "show output on stdio")
	flag.BoolVar(&versionVar, "version", false, "show version and exit")
	flag.BoolVar(&typeVar, "type", false, "output format is YAML instead of JSON")
	flag.Var(&outputFileVar, "output", "filename to write result to")
}

func main() {

	flag.Parse()

	checkVersion(versionVar)

	data := readFile(fileVar.value)
	fileExt := filepath.Ext(fileVar.value[:1])
	manifest := toJSON(data, filepath.Ext(fileVar.value)[1:])

	json := fmtJSON(manifest, typeVar)

	outputName := outputFileVar.value
	if !outputFileVar.set {
		outputName = fileVar.value[:len(fileVar.value)-len(fileExt)] + ".json"
	}
	writeFile(outputName, json, stdioVar)
}
