// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package yaml

import "github.com/coveo/gotemplate/collections"

// List implementation of IGenericList for yamlList
type List = yamlList
type yamlIList = collections.IGenericList
type yamlList []interface{}

func (l yamlList) AsArray() []interface{} { return []interface{}(l) }
func (l yamlList) Cap() int               { return cap(l) }
func (l yamlList) Capacity() int          { return cap(l) }
func (l yamlList) Clone() yamlIList       { return yamlListHelper.Clone(l) }
func (l yamlList) Contains(values ...interface{}) bool {
	return yamlListHelper.Contains(l, values...)
}
func (l yamlList) Count() int                   { return len(l) }
func (l yamlList) Create(args ...int) yamlIList { return yamlListHelper.CreateList(args...) }
func (l yamlList) CreateDict(args ...int) yamlIDict {
	return yamlListHelper.CreateDictionary(args...)
}
func (l yamlList) First() interface{} { return yamlListHelper.GetIndexes(l, 0) }
func (l yamlList) Get(indexes ...int) interface{} {
	return yamlListHelper.GetIndexes(l, indexes...)
}
func (l yamlList) Has(values ...interface{}) bool    { return l.Contains(values...) }
func (l yamlList) Last() interface{}                 { return yamlListHelper.GetIndexes(l, len(l)-1) }
func (l yamlList) Len() int                          { return len(l) }
func (l yamlList) New(args ...interface{}) yamlIList { return yamlListHelper.NewList(args...) }
func (l yamlList) Reverse() yamlIList                { return yamlListHelper.Reverse(l) }
func (l yamlList) Strings() []string                 { return yamlListHelper.GetStrings(l) }
func (l yamlList) StringArray() strArray             { return yamlListHelper.GetStringArray(l) }
func (l yamlList) TypeName() str                     { return "Yaml" }
func (l yamlList) Join(sep interface{}) str          { return l.StringArray().Join(sep) }
func (l yamlList) Unique() yamlIList                 { return yamlListHelper.Unique(l) }

func (l yamlList) GetHelpers() (collections.IDictionaryHelper, collections.IListHelper) {
	return yamlDictHelper, yamlListHelper
}

func (l yamlList) Append(values ...interface{}) yamlIList {
	return yamlListHelper.Add(l, false, values...)
}

func (l yamlList) Intersect(values ...interface{}) yamlIList {
	return yamlListHelper.Intersect(l, values...)
}

func (l yamlList) Pop(indexes ...int) (interface{}, yamlIList) {
	if len(indexes) == 0 {
		indexes = []int{len(l) - 1}
	}
	return l.Get(indexes...), l.Remove(indexes...)
}

func (l yamlList) Prepend(values ...interface{}) yamlIList {
	return yamlListHelper.Add(l, true, values...)
}

func (l yamlList) Remove(indexes ...int) yamlIList {
	return yamlListHelper.Remove(l, indexes...)
}

func (l yamlList) Set(i int, v interface{}) (yamlIList, error) {
	return yamlListHelper.SetIndex(l, i, v)
}

func (l yamlList) Union(values ...interface{}) yamlIList {
	return yamlListHelper.Add(l, false, values...).Unique()
}

func (l yamlList) Without(values ...interface{}) yamlIList {
	return yamlListHelper.Without(l, values...)
}

// Dictionary implementation of IDictionary for yamlDict
type Dictionary = yamlDict
type yamlIDict = collections.IDictionary
type yamlDict map[string]interface{}

func (d yamlDict) Add(key, v interface{}) yamlIDict    { return yamlDictHelper.Add(d, key, v) }
func (d yamlDict) AsMap() map[string]interface{}       { return (map[string]interface{})(d) }
func (d yamlDict) Native() interface{}                 { return collections.ToNativeRepresentation(d) }
func (d yamlDict) Count() int                          { return len(d) }
func (d yamlDict) Len() int                            { return len(d) }
func (d yamlDict) Clone(keys ...interface{}) yamlIDict { return yamlDictHelper.Clone(d, keys) }
func (d yamlDict) Create(args ...int) yamlIDict        { return yamlListHelper.CreateDictionary(args...) }
func (d yamlDict) CreateList(args ...int) yamlIList    { return yamlHelper.CreateList(args...) }
func (d yamlDict) Flush(keys ...interface{}) yamlIDict { return yamlDictHelper.Flush(d, keys) }
func (d yamlDict) Get(keys ...interface{}) interface{} { return yamlDictHelper.Get(d, keys) }
func (d yamlDict) Has(keys ...interface{}) bool        { return yamlDictHelper.Has(d, keys) }
func (d yamlDict) GetKeys() yamlIList                  { return yamlDictHelper.GetKeys(d) }
func (d yamlDict) KeysAsString() strArray              { return yamlDictHelper.KeysAsString(d) }
func (d yamlDict) Pop(keys ...interface{}) interface{} { return yamlDictHelper.Pop(d, keys) }
func (d yamlDict) GetValues() yamlIList                { return yamlDictHelper.GetValues(d) }
func (d yamlDict) Set(key, v interface{}) yamlIDict    { return yamlDictHelper.Set(d, key, v) }
func (d yamlDict) Transpose() yamlIDict                { return yamlDictHelper.Transpose(d) }
func (d yamlDict) TypeName() str                       { return "Yaml" }

func (d yamlDict) GetHelpers() (collections.IDictionaryHelper, collections.IListHelper) {
	return yamlDictHelper, yamlListHelper
}

func (d yamlDict) Default(key, defVal interface{}) interface{} {
	return yamlDictHelper.Default(d, key, defVal)
}

func (d yamlDict) Delete(key interface{}, otherKeys ...interface{}) (yamlIDict, error) {
	return yamlDictHelper.Delete(d, append([]interface{}{key}, otherKeys...))
}

func (d yamlDict) Merge(dict yamlIDict, otherDicts ...yamlIDict) yamlIDict {
	return yamlDictHelper.Merge(d, append([]yamlIDict{dict}, otherDicts...))
}

func (d yamlDict) Omit(key interface{}, otherKeys ...interface{}) yamlIDict {
	return yamlDictHelper.Omit(d, append([]interface{}{key}, otherKeys...))
}

// Generic helpers to simplify physical implementation
func yamlListConvert(list yamlIList) yamlIList { return yamlList(list.AsArray()) }
func yamlDictConvert(dict yamlIDict) yamlIDict { return yamlDict(dict.AsMap()) }
func needConversion(object interface{}, strict bool) bool {
	return needConversionImpl(object, strict, "Yaml")
}

var yamlHelper = helperBase{ConvertList: yamlListConvert, ConvertDict: yamlDictConvert, NeedConversion: needConversion}
var yamlListHelper = helperList{BaseHelper: yamlHelper}
var yamlDictHelper = helperDict{BaseHelper: yamlHelper}

// DictionaryHelper gives public access to the basic dictionary functions
var DictionaryHelper collections.IDictionaryHelper = yamlDictHelper

// GenericListHelper gives public access to the basic list functions
var GenericListHelper collections.IListHelper = yamlListHelper

type (
	str      = collections.String
	strArray = collections.StringArray
)

var iif = collections.IIf
