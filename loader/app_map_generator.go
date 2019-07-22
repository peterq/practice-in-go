//+build ignore

package main

import (
	"io/ioutil"
	"log"
	"os"
	"html/template"
)

type initFunc func()

func main() {

	file, err := os.Open(tplFile)

	if err != nil {
		log.Println(os.Getwd())
		log.Fatalf("Failed to open %s: %q", tplFile, err)
	}

	data, err := ioutil.ReadAll(file)

	tplStr := string(data)

	apps := make([]string, 0)

	apps = append(apps, scanApps()...)

	tpl, err := template.New("tpl").Parse(tplStr)
	if err != nil  {
		log.Fatal(err)
	}

	f, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.Execute(f, apps)
	if err != nil {
		log.Fatal(err)
	}

}

func scanApps() []string {
	apps := make([]string, 0)
	info, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range info {
		if item.IsDir() && (hasFile(item.Name() + "/init.go") ||
			hasFile(item.Name() + "/" + item.Name() + "_init.go")) {
			apps = append(apps, item.Name())
		}
	}

	return apps
}


func hasFile(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
const (
	tplFile      = "loader/apps.go.tpl"
	dstFile      = "loader/apps.go"
)
