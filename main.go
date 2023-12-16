package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"text/template"
	"time"
)

const VERSION = `0.1`

var (
	ErrorLog = log.New(os.Stderr, `error#`, log.Lshortfile)
	DebugLog = log.New(os.Stdout, `debug#`, log.Lshortfile)
)

type TemplateData struct {
	FileName string
	Name     string
	Index    int
	ModTime  time.Time
	Size     int64
}

func helpText() {
	fmt.Println(`https://github.com/vvampirius/batch-file-rename`)
	flag.PrintDefaults()
}

func GetDstName(n int, srcFileName string, nameRegexp *regexp.Regexp, dstTemplate *template.Template) (string, error) {
	fileInfo, err := os.Stat(srcFileName)
	if err != nil {
		return "", err
	}
	if fileInfo.IsDir() {
		return "", errors.New(`is directory`)
	}
	matchName := nameRegexp.FindStringSubmatch(srcFileName)
	if len(matchName) != 2 {
		return "", errors.New(`can't parse name`)
	}
	buffer := bytes.NewBuffer(nil)
	data := TemplateData{
		FileName: srcFileName,
		Name:     matchName[1],
		Index:    n,
		ModTime:  fileInfo.ModTime(),
		Size:     fileInfo.Size(),
	}
	if err := dstTemplate.Execute(buffer, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func main() {
	help := flag.Bool("h", false, "print this help")
	ver := flag.Bool("v", false, "Show version")
	nameRegexpFlag := flag.String("name", `(.*)`, "Name regexp")
	templateFlag := flag.String("template", `{{.Name}}`, "Dst filename template")
	testFlag := flag.Bool("test", false, "Don't really rename! Only print new filename.")
	flag.Parse()

	if *help {
		helpText()
		os.Exit(0)
	}

	if *ver {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	srcFileNames := flag.Args()
	if len(srcFileNames) == 0 {
		helpText()
		os.Exit(1)
	}

	nameRegexp, err := regexp.Compile(*nameRegexpFlag)
	if err != nil {
		ErrorLog.Fatalln(err.Error())
	}

	dstTemplate, err := template.New(`dst`).Parse(*templateFlag)
	if err != nil {
		ErrorLog.Fatalln(err.Error())
	}

	errsCount := 0
	n := 0
	for _, srcFileName := range srcFileNames {
		n++
		dstFileName, err := GetDstName(n, srcFileName, nameRegexp, dstTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR '%s': %s\n", srcFileName, err.Error())
			errsCount++
			continue
		}
		fmt.Printf("'%s' -> '%s'\n", srcFileName, dstFileName)
		if srcFileName == dstFileName {
			fmt.Fprintf(os.Stderr, "ERROR: '%s' = '%s'\n", srcFileName, dstFileName)
			errsCount++
			continue
		}
		if *testFlag {
			continue
		}
		if err := os.Rename(srcFileName, dstFileName); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR '%s': %s\n", srcFileName, err.Error())
			errsCount++
			continue
		}
	}
	if errsCount > 0 {
		os.Exit(1)
	}
}
