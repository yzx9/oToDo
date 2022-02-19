package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const input = `model`
const output = `tmp`

type customStruct struct {
	File          string
	Name          string
	Fields        []customStructField
	InheritEntity bool
}

type customStructField struct {
	Name    string
	Type    string
	Comment string
}

func main() {
	files, err := getFilePaths()
	if err != nil {
		log.Fatal(err)
	}

	structs := make([]customStruct, 0)
	for _, file := range files {
		s, err := parseFile(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		structs = append(structs, s...)
	}

	for i := range structs {
		if err := generator(structs[i]); err != nil {
			fmt.Println(err)
		}
	}
}

var scanFiles = regexp.MustCompile(`.go$`)
var ignoreFiles = regexp.MustCompile(`README.md$`)

func getFilePaths() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("fails to get work dir: %w", err)
	}

	root := filepath.Join(wd, input)
	fmt.Printf("Scan: %v\n", root)

	var files []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() &&
			scanFiles.MatchString(info.Name()) &&
			!ignoreFiles.MatchString(info.Name()) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

var structBeginRegex = regexp.MustCompile(`type (?P<Name>\S+) struct {`)
var structInheritEntityRegex = regexp.MustCompile(`\s+Entity`)
var structFieldRegex = regexp.MustCompile(`\s+(?P<Name>\S+)\s+(?P<Type>\S+)(\s+(?P<Comment>\S+))?`)
var structEndRegex = regexp.MustCompile(`}`)

func parseFile(filepath string) ([]customStruct, error) {
	lines, err := readLines(filepath)
	if err != nil {
		return nil, fmt.Errorf("fails to read file: %w", err)
	}

	s := customStruct{File: filepath}
	arr := make([]customStruct, 0)

	processField := func(line string) {
		matches := structFieldRegex.FindStringSubmatch(line)
		f := customStructField{
			Name: matches[1],
			Type: matches[2],
		}

		if len(matches) > 3 {
			f.Comment = matches[3]
		}

		s.Fields = append(s.Fields, f)
	}

	state := 0
	for i := range lines {
		switch state {
		case 0:
			if structBeginRegex.MatchString(lines[i]) {
				state = 1

				matches := structBeginRegex.FindStringSubmatch(lines[i])
				s.Name = matches[1]
			}

		case 1:
			if structFieldRegex.MatchString(lines[i]) {
				processField(lines[i])

			} else if structInheritEntityRegex.MatchString(lines[i]) {
				state = 2

				s.InheritEntity = true
			} else if structEndRegex.MatchString(lines[i]) {
				state = 3
			}

		case 2:
			if structFieldRegex.MatchString(lines[i]) {
				processField(lines[i])

			} else if structEndRegex.MatchString(lines[i]) {
				state = 3
			}

		case 3:
			arr = append(arr, s)
			s = customStruct{File: filepath}
			state = 0
		}
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("not found any struct in file: %v", filepath)
	}

	return arr, nil
}

func readLines(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("fails to open file: %v", filepath)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("fails to read file: %v", filepath)
	}

	return strings.Split(string(content), "\n"), nil
}

const injectEntity = `
  id: number
  createdAt: string
  updatedAt: string`

var jsonComment = regexp.MustCompile("json:\"(?P<name>[\\w-]+)(;|\")")

func generator(s customStruct) error {
	content := "interface " + s.Name + " {"
	if s.InheritEntity {
		content += injectEntity
	}
	content += "\n"

	for i := range s.Fields {
		if c := s.Fields[i].Comment; c != "" && jsonComment.MatchString(c) {
			matches := jsonComment.FindStringSubmatch(c)
			if matches[1] == "-" {
				continue
			}

			s.Fields[i].Name = matches[1]
		}

		tsType, err := transformType(s.Fields[i].Type)
		if err != nil {
			return fmt.Errorf("fails to transform type in file `%v`: %w", s.File, err)
		}

		content += fmt.Sprintf("  %v: %v\n", s.Fields[i].Name, tsType)
	}
	content += "}\n"

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("fails to get work dir: %w", err)
	}

	dest := filepath.Join(wd, output, s.Name+".ts")
	if err := os.WriteFile(dest, []byte(content), os.ModePerm); err != nil {
		return fmt.Errorf("fails to save file: %w", err)
	}

	fmt.Printf("generate %v\n", dest)
	return nil
}

func transformType(goType string) (string, error) {
	goType = strings.TrimPrefix(goType, "*")

	switch goType {
	case "bool":
		return "boolean", nil

	case "string":
		return "string", nil

	case "time.Time":
		return "string", nil

	case "[]byte":
		return "string", nil

	case "int":
		return "number", nil

	case "int8":
		return "number", nil

	case "int16":
		return "number", nil

	case "int32":
		return "number", nil

	case "int64":
		return "number", nil

	case "float32":
		return "number", nil

	case "float64":
		return "number", nil
	}

	return "", fmt.Errorf("unknown go type: %v", goType)
}
