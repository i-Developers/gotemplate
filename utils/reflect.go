package utils

import (
	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/stringclass"
)

type str = stringclass.String

var isEmpty = collections.IsEmptyValue

// MergeDictionaries merges multiple dictionaries into a single one prioritizing the first ones.
func MergeDictionaries(args ...map[string]interface{}) (map[string]interface{}, error) {
	if len(args) == 0 {
		return make(map[string]interface{}), nil
	}
	dicts := make([]collections.IDictionary, len(args))
	for i := range dicts {
		var err error
		dicts[i], err = collections.TryAsDictionary(args[i])
		if err != nil {
			return nil, err
		}
	}

	result := collections.CreateDictionary()
	return result.Merge(dicts[0], dicts[1:]...).Native().(map[string]interface{}), nil
}
