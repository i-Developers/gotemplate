# gotemplate

## Description

Apply go template over files ending with .template in the current directory.

For more information on Go Template functionality, check this [link](https://golang.org/pkg/text/template).

This little utility just scan the current folder for `*.template` files and apply the [go template](https://golang.org/pkg/text/template) over them.

Every matching `*.ext.template` file will render a file named `*.generated.ext`. Other matched file (if --pattern is supplied) will replace the file with the rendered content and rename the original file `*.ext.original`.

## Usage

```text
usage: gotemplate [<flags>]

A template processor for go.

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
  -i, --import=file ...       Import variables files (could be any of YAML, JSON or HCL format)
  -p, --patterns=pattern ...  Additional patterns that should be processed by gotemplate
  -o, --overwrite             Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)
  -s, --substitute=exp ...    Substitute text in the processed files by applying the regex substitute expression (format: regex/substitution, see: Go regexp)
  -r, --recursive             Process all template files recursively
      --source=SOURCE         Specify a source folder (default to the current folder)
      --target=TARGET         Specify a target folder (default to source folder)
  -f, --follow-symlinks       Follow the symbolic links while using the recursive option
  -d, --dry-run               Do not actually overwrite files, just show the result
  -l, --list-functions        List the available functions
  -v, --version               Get the current version of gotemplate
```

## Available functions

### Base functions

* All functions included in [Go Template](https://golang.org/pkg/text/template) 
* All function included in [Sprig](http://masterminds.github.io/sprig)

## Structured import functions

All following function allows you to import data in your template. The source could be either a template, a file path or
a direct data definition.

Function name | Argument(s() |Description
--- | --- | ---
yaml | source [context] | Convert the data referred by `source` with `context` as argument into a go interface (the resulting YAML must be valid)
json | source [context] | Convert the data referred by `source` with `context` as argument into a go interface (the resulting JSON must be valid)
hcl | source [context] | Convert the data referred by `source` with `context` as argument into a go interface (the resulting HCL must be valid)

_Usage:_

```go
// JSON data definition
{{ $json_data := `"number": 10, "text": "Hello World!, "data": { "a": 1 }, "list": ["value 1, "value 2"]` | json }}

// YAML data definition
{{ $yaml_data := yaml `
number: 10
text: Hello world!
data:
    a: 1
    b: 2
list:
    - value 1
    - value 2
` }}

// HCL data definition from template
{{ define "source" }}
number = 10
text = "Hello world!"

data {
    a = 1
    b = 2
}

list = [
    "value 1",
    "value 2",
]
{{ end }}

{{ $hcl_data := hcl "source" }}

// YAML data defintion for file
{{ $ymal_file := yaml "source.yml" }}
```

# Other functions

Function name | Argument(s() |Description
--- | --- | ---
bool | string | Convert the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`).
concat | objects ... | Returns the concatenation (without separator) of the string representation of objects.
exec | command string, args ... | Returns the result of the shell command as string.
formatList | format string, list | Return a list of strings by applying the format to each element of the supplied list.
mergeLists | lists | Return a single list containing all elements from the lists supplied.
glob | args ... | Returns the expanded list of supplied arguments (expand *[]? on filename).
lorem | type string, min, max int | Returns a random string. Valid types are be `word`, `words`, `sentence`, `para`, `paragraph`, `host`, `email`, `url`.
toYaml | interface | Returns the YAML string representation of the supplied object.

```go
// Random text
{{ $text := lorem "paragraph" 5 10 }}

// Random email
{{ $email := lorem "email" }}

// Concatenation of various objects
{{ $string := concat "Hello" 1 " pi = " 3.14 " 3 x 5 = " (mul 3 5)}}

// Boolean conversion
{{ $bool := bool "True" }}

// Quote all elements of a list
{{ $quoted := list 1 2 3 | formatList "\"%v\"" }}
```
