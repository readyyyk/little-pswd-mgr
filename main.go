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

//var dataFile os.File

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
func logData(dataSet []record) {
	if len(dataSet) == 0 {
		fmt.Println(text.FgRed.Sprint("No data found"))
		return
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.AppendHeader(table.Row{"data", "user", text.FgHiBlack.Sprint(delitmer), "host"})
	for _, data := range dataSet {
		t.AppendRow(table.Row{data.Data, data.User, text.FgHiBlack.Sprint(delitmer), data.Host})
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, AlignFooter: text.AlignLeft, AlignHeader: text.AlignCenter},
		{Number: 2, Align: text.AlignRight, AlignFooter: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 3, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 4, Align: text.AlignLeft, AlignFooter: text.AlignLeft, AlignHeader: text.AlignLeft},
	})
	fmt.Println(t.Render())
}
func logHelp() {
	tempHelpData := table.NewWriter()
	tempHelpData.Style().Options.DrawBorder = false
	tempHelpData.Style().Options.SeparateColumns = false
	tempHelpData.Style().Options.SeparateHeader = false
	tempHelpData.AppendRows([]table.Row{
		{"\t<no flag>", "shows all your stored data"}, {""},
		{"\t-h, --help", `shows this message`}, {""},
		{"\t-a, --add", "adds new record, after space provide information you want to save\nin format " + `"<data> <user>@<host>"`}, {""},
		{"\t-d, --del", "deletes record with provided credentials, after space provide\ncredentials of record you want to delete in format <user>@<host>"}, {""},
		{"\t-s, --sort", "shows only strings that contain provided data"},
	})
	fmt.Println(tempHelpData.Render())
}

func init() {
	e, _ := os.Executable()
	dataPath = filepath.Dir(e) + dataPath

	tokensData.read()
}
func main() {
	if len(os.Args) == 1 {
		logData(tokensData.Data)
		return
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		logHelp()
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
	}
	if os.Args[1] == "--del" || os.Args[1] == "-d" {
		atIndex := strings.Index(os.Args[2], "@")
		if atIndex == -1 {
			fmt.Println(text.FgRed.Sprint("provide valid data"))
			return
		}
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
	}
	if os.Args[1] == "--sort" || os.Args[1] == "-s" {
		if len(os.Args) < 3 {
			fmt.Println(text.FgRed.Sprint("Provide data to sort (any string)"))
			return
		}

		var dataSet []record
		for _, el := range tokensData.Data {
			if strings.Contains(el.Data, os.Args[2]) || strings.Contains(el.User, os.Args[2]) || strings.Contains(el.Host, os.Args[2]) {
				if strings.Contains(el.Data, os.Args[2]) {
					el.Data = text.FgHiBlue.Sprint(el.Data)
				}
				if strings.Contains(el.User, os.Args[2]) {
					el.User = text.FgHiBlue.Sprint(el.User)
				}
				if strings.Contains(el.Host, os.Args[2]) {
					el.Host = text.FgHiBlue.Sprint(el.Host)
				}
				dataSet = append(dataSet, el)
			}
		}
		logData(dataSet)
		return
	}
}
