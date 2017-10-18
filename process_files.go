package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/coveo/terragrunt/util"
)

var basicFunctions = []string{
	"and", "call", "call", "html", "index", "js", "len", "not", "or", "print", "printf", "println", "urlquery",
	"eq", "ge", "gt", "le", "lt", "ne",
}

const templateExt = ".template"

func processFiles(context interface{}) {
	// Create the evaluation context
	mainTemplate := template.New("context").Funcs(sprig.GenericFuncMap())
	AddExtraFuncs(mainTemplate)

	*includePatterns = append(*includePatterns, "*"+templateExt)

	if *listFunc {
		printFunctions(mainTemplate)
		os.Exit(0)
	}

	pwd := PanicOnError(os.Getwd()).(string)
	var files []string
	if *recursive {
		files = PanicOnError(findFiles(pwd, *includePatterns...)).([]string)
	} else {
		files = PanicOnError(util.FindFiles(pwd, *includePatterns...)).([]string)
	}

	if len(files) == 0 {
		// There is nothing to process
		return
	}

	// Parse the template files
	var validTemplates []string
	templateFiles, files := isolateTemplateFiles(files)
	for _, filename := range templateFiles {
		// We do not use ParseFiles because it names the template with the base name of the file
		// which result in overriding templates with the same base name in different folders.
		if _, err := mainTemplate.New(filename).Parse(string(PanicOnError(ioutil.ReadFile(filename)).([]byte))); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			validTemplates = append(validTemplates, filename)
		}
	}

	var terraformFiles bool

	// Process the files and generate resulting file if there is an output from the template
	for _, file := range validTemplates {
		if *dryRun {
			fmt.Println("Processing file", file)
		}

		out, err := executeTemplate(mainTemplate.Lookup(file), file, context)
		if err != nil {
			if *dryRun {
				fmt.Fprintln(os.Stderr, err, "while executing", file)
				continue
			} else {
				switch err := err.(type) {
				case template.ExecError:
					fmt.Fprintln(os.Stderr, err)
				default:
					PanicOnError(err)
				}
			}
		}

		fileName := strings.TrimSuffix(file, templateExt)
		ext := path.Ext(fileName)
		fileName = fmt.Sprint(strings.TrimSuffix(fileName, ext), ".generated", ext)
		if *dryRun {
			fmt.Printf("  Would create file %s with content:\n%s\n", fileName, out.String())
		} else {
			info, _ := os.Stat(file)
			if err := ioutil.WriteFile(fileName, []byte(out.String()), info.Mode()); err != nil {
				PanicOnError(err)
			}
		}

		terraformFiles = terraformFiles || ext == ".tf" || ext == ".tf.json"
	}

	for _, file := range files {
		if *dryRun {
			fmt.Println("Processing file", file)
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err, "while reading", file)
			continue
		}

		if strings.TrimSpace(string(content)) == "" {
			// For a weird reason, empty file generates the previous template content
			continue
		}

		template, err := mainTemplate.Parse(string(content))
		if err != nil {
			fmt.Fprintln(os.Stderr, err, "while parsing", file)
			continue
		}
		out, err := executeTemplate(template, file, context)
		if err != nil {
			fmt.Fprintln(os.Stderr, err, "while executing", file)
			continue
		}

		if out.String() != string(content) {
			if *dryRun {
				fmt.Printf("  Would modify file %s with content:\n%s\n", file, out.String())
			} else {
				info, _ := os.Stat(file)
				if err := os.Rename(file, file+".original"); err != nil {
					PanicOnError(err)
				}
				ioutil.WriteFile(file, out.Bytes(), info.Mode())
				ext := path.Ext(file)
				terraformFiles = terraformFiles || ext == ".tf" || ext == ".tf.json"
			}
		}
	}

	// If terraform is present, we apply a terraform fmt to the resulting templates
	// to ensure that the resulting files are valids.
	if _, err := exec.LookPath("terraform"); err == nil && terraformFiles {
		cmd := exec.Command("terraform", "fmt")
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// RunningTemplate represents the current running template (this avoid concurrent processing)
var RunningTemplate *template.Template

func executeTemplate(template *template.Template, fileName string, context interface{}) (bytes.Buffer, error) {
	var out bytes.Buffer

	RunningTemplate = template
	RunningTemplate.ParseName = fileName

	err := RunningTemplate.Execute(&out, context)
	return out, err
}

func findFiles(folder string, patterns ...string) ([]string, error) {
	visited := map[string]bool{}
	var walker func(folder string) ([]string, error)
	walker = func(folder string) ([]string, error) {
		var results []string
		folder, _ = filepath.Abs(folder)

		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				visited[path] = true
				files, err := util.FindFiles(path, *includePatterns...)
				if err != nil {
					return err
				}
				results = append(results, files...)
				return nil
			}

			if info.Mode()&os.ModeSymlink != 0 && *followSymLinks {
				link, err := os.Readlink(path)
				if err != nil {
					fmt.Fprintln(os.Stderr, err, "while trying to follow link", path)
					return nil
				}

				if !filepath.IsAbs(link) {
					link = filepath.Join(filepath.Dir(path), link)
				}
				link, _ = filepath.Abs(link)
				if !visited[link] {
					// Check if we already visited that link to avoid recursive loop
					linkFiles, err := walker(link)
					if err != nil {
						return err
					}
					results = append(results, linkFiles...)
				}
			}
			return nil
		})
		return results, nil
	}
	return walker(folder)
}

func isolateTemplateFiles(files []string) (templates, others []string) {
	for _, file := range files {
		if filepath.Ext(file) == templateExt {
			templates = append(templates, file)
		} else {
			others = append(others, file)
		}
	}
	return
}

func printFunctions(template *template.Template) {
	common := reflect.ValueOf(*template).FieldByName("common").Elem().FieldByName("parseFuncs")
	keys := common.MapKeys()
	functions := make([]string, len(keys))
	for i, k := range keys {
		functions[i] = k.String()
	}

	functions = append(basicFunctions, functions...)
	sort.Strings(functions)

	const nbColumn = 5
	colLength := int(math.Ceil(float64(len(functions)) / float64(nbColumn)))
	maxLength := 0

	// Initialize the columns to sort function per column
	var list [nbColumn][]string
	for i := range list {
		list[i] = make([]string, colLength)
	}

	// Place functions into columns
	for i, function := range functions {
		column := list[i/colLength]
		column[i%colLength] = function
		maxLength = int(math.Max(float64(len(function)), float64(maxLength)))
	}

	// Print the columns
	for i := range list[0] {
		for _, column := range list {
			fmt.Printf("%-[1]*[2]s", maxLength+2, column[i])
		}
		fmt.Println()
	}
}
