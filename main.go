package main

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"path/filepath"
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

func (r *uhPair) eq(pair uhPair) bool {
	return r.User == pair.User && r.Host == pair.Host
}

type dataS struct {
	Data []record
}

/*
	type dataI interface {
		read()
		add(data record)
		delete(data uhPair)
		log()
	}
*/
func (d *dataS) rewriteFile() {
	dataJSON, err := json.Marshal(d.Data)
	logError(err)
	err = os.WriteFile(dataPath, dataJSON, 0777)
	logError(err)
}
func (d *dataS) read() {
	dataJSON, err := os.ReadFile(dataPath)
	if err != nil {
		_, err := os.Create(dataPath)
		logError(err)
	}
	if len(dataJSON) == 0 {
		dataJSON = []byte("[]")
	}
	logError(json.Unmarshal(dataJSON, &d.Data))
}
func (d *dataS) add(data record) {
	d.Data = append(d.Data, data)
	d.rewriteFile()
}
func (d *dataS) delete(data uhPair) (deleted bool) {
	for i, el := range d.Data {
		if data.eq(uhPair{el.User, el.Host}) {
			d.Data = append(d.Data[:i], d.Data[i+1:]...)
			d.rewriteFile()
			return true
		}
	}
	return false
}
func (d *dataS) log() {
	if len(tokensData.Data) == 0 {
		fmt.Println(text.FgRed.Sprint("No data found"))
		return
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.AppendHeader(table.Row{"data", "user", text.FgHiWhite.Sprint(delitmer), "host"})
	for _, data := range tokensData.Data {
		t.AppendRow(table.Row{data.Data, data.User, text.FgHiWhite.Sprint(delitmer), data.Host})
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, AlignFooter: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Number: 2, Align: text.AlignRight, AlignFooter: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 3, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 4, Align: text.AlignLeft, AlignFooter: text.AlignLeft, AlignHeader: text.AlignLeft},
	})
	fmt.Println(t.Render())
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
		spaceIndex := strings.Index(os.Args[2], " ")
		atIndex := strings.Index(os.Args[2], "@")
		tokensData.add(record{
			Data: os.Args[2][:spaceIndex],
			User: os.Args[2][spaceIndex+1 : atIndex],
			Host: os.Args[2][atIndex+1:],
		})
		fmt.Println(text.FgGreen.Sprint("Added"))
		return
	} else if os.Args[1] == "--del" || os.Args[1] == "-d" {
		atIndex := strings.Index(os.Args[2], "@")
		data := uhPair{
			User: os.Args[2][:atIndex],
			Host: os.Args[2][atIndex+1:],
		}
		if tokensData.delete(data) {
			fmt.Println(text.FgGreen.Sprint("Deleted"))
		} else {
			fmt.Println(text.FgRed.Sprint("Not found record with user: ") + text.BgHiBlue.Sprint(data.User) + text.FgRed.Sprint(" for host: ") + text.BgHiBlue.Sprint(data.Host))
		}
		return
	} else if os.Args[1] == "--help" || os.Args[1] == "-h" {

	}
}
