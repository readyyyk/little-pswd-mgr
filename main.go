package main

import (
	"encoding/json"
	"fmt"
	"github.com/rodaine/table"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func logError(err error) {
	if err != nil {
		panic(err)
	}
}

var delitmer = "@"
var dataPath = "/.tokens"

var tokensData dataS
var dataFile os.File

type record struct {
	Data string `json:"data"`
	User string `json:"user"`
	Host string `json:"host"`
}

type uhPair struct {
	User string `json:"user,omitempty"`
	Host string `json:"host,omitempty"`
}

type dataS struct {
	Data []record
}

/*type dataI interface {
	read()
	add(data record)
	delete(data uhPair)
	log()
}*/

func (d *dataS) read() {
	dataJSON, err := os.ReadFile(dataPath)
	if err != nil {
		_, err := os.Create(dataPath)
		logError(err)
	}
	if len(dataJSON) == 0 {
		dataJSON = []byte("[]")
	}
	//dataJSON, err := io.ReadAll(&dataFile)
	//dataJSON, err := io.ReadAll(&dataFile)
	//fmt.Println(dataJSON, dataFile)

	logError(json.Unmarshal(dataJSON, &d.Data))
}
func (d *dataS) add(data record) {
	d.Data = append(d.Data, data)

	dataJSON, err := json.Marshal(d.Data)
	logError(err)
	//_, err = dataFile.Write(dataJSON)
	err = os.WriteFile(dataPath, dataJSON, 0644)
	logError(err)
}
func (d *dataS) delete(data uhPair) (deleted bool) {
	for i, el := range d.Data {
		if reflect.DeepEqual(uhPair{el.User, el.Host}, data) {
			d.Data = append(d.Data[:i], d.Data[i+1:]...)
			return true
		}
	}
	return false
}
func (d *dataS) log() {
	// github.com/jedib0t/go-pretty/v6/table
	tbl := table.New("data", "user", "@", "host")
	tbl.WithPadding(3)
	if len(tokensData.Data) == 0 {
		fmt.Println("No data found")
		return
	}
	for _, data := range tokensData.Data {
		tbl.AddRow(data.Data, data.User, delitmer, data.Host)
	}
	tbl.Print()
}

func init() {
	e, _ := os.Executable()
	dataPath = filepath.Dir(e) + dataPath

	//tempFile, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_RDONLY, os.ModePerm)
	//if err != nil {
	//	tempFile, _ = os.Create(dataPath)
	//	_, _ = tempFile.Write([]byte("[{}]"))
	//}
	//dataFile = *tempFile

	tokensData.read()
}
func main() {
	if len(os.Args) == 1 {
		tokensData.log()
		return
	}

	if os.Args[1] == "--add" || os.Args[1] == "-a" {
		fmt.Println(len(os.Args[2]))
		spaceIndex := strings.Index(os.Args[2], " ")
		atIndex := strings.Index(os.Args[2], "@")
		tokensData.add(record{
			Data: os.Args[2][:spaceIndex],
			User: os.Args[2][spaceIndex+1 : atIndex],
			Host: os.Args[2][atIndex+1:],
		})
		return
	} else if os.Args[1] == "--delete" || os.Args[1] == "-d" {
		atIndex := strings.Index(os.Args[2], "@")
		tokensData.delete(uhPair{
			User: os.Args[2][:atIndex],
			Host: os.Args[2][atIndex+1:],
		})
		return
	}
}
