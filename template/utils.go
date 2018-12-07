package template

import (
	"path/filepath"

	"github.com/coveo/gotemplate/collections"
	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
)

var must = errors.Must
var trapError = errors.Trap
var iif = collections.IIf
var ifUndef = collections.IfUndef
var defval = collections.Default

type (
	iList       = collections.IGenericList
	iDictionary = collections.IDictionary
)

func getTargetFile(targetFile, sourcePath, targetPath string) string {
	if targetPath != "" {
		targetFile = filepath.Join(targetPath, utils.Relative(sourcePath, targetFile))
	}
	return targetFile
}