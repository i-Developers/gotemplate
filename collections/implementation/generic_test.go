package implementation

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/coveo/gotemplate/v3/errors"
	"github.com/stretchr/testify/assert"
)

var strFixture = baseList(baseListHelper.NewStringList(strings.Split("Hello World, I'm Foo Bar!", " ")...).AsArray())

func Test_list_Append(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		l      baseIList
		values []interface{}
		want   baseIList
	}{
		{"Empty", baseList{}, []interface{}{1, 2, 3}, baseList{1, 2, 3}},
		{"List of int", baseList{1, 2, 3}, []interface{}{4, 5}, baseList{1, 2, 3, 4, 5}},
		{"List of string", strFixture, []interface{}{"That's all folks!"}, baseList{"Hello", "World,", "I'm", "Foo", "Bar!", "That's all folks!"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Append(tt.values...))
		})
	}
}

func Test_list_Prepend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		l      baseIList
		values []interface{}
		want   baseIList
	}{
		{"Empty", baseList{}, []interface{}{1, 2, 3}, baseList{1, 2, 3}},
		{"List of int", baseList{1, 2, 3}, []interface{}{4, 5}, baseList{4, 5, 1, 2, 3}},
		{"List of string", strFixture, []interface{}{"That's all folks!"}, baseList{"That's all folks!", "Hello", "World,", "I'm", "Foo", "Bar!"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Prepend(tt.values...))
		})
	}
}

func Test_list_AsArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		want []interface{}
	}{
		{"Empty List", baseList{}, []interface{}{}},
		{"List of int", baseList{1, 2, 3}, []interface{}{1, 2, 3}},
		{"List of string", strFixture, []interface{}{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.AsArray())
		})
	}
}

func Test_baseList_Strings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		want []string
	}{
		{"Empty List", baseList{}, []string{}},
		{"List of int", baseList{1, 2, 3}, []string{"1", "2", "3"}},
		{"List of string", strFixture, []string{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Strings())
		})
	}
}

func Test_list_Capacity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseIList
		want int
	}{
		{"Empty List with 100 spaces", baseListHelper.CreateList(0, 100), 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Capacity())
			assert.Equal(t, tt.l.Cap(), tt.l.Capacity(), "Cap and Capacity return different values")
		})
	}
}

func Test_list_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		want baseIList
	}{
		{"Empty List", baseList{}, baseList{}},
		{"List of int", baseList{1, 2, 3}, baseList{1, 2, 3}},
		{"List of string", strFixture, baseList{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Clone())
		})
	}
}

func Test_list_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		l       baseList
		indexes []int
		want    interface{}
	}{
		{"Empty List", baseList{}, []int{0}, nil},
		{"Negative index", baseList{}, []int{-1}, nil},
		{"List of int", baseList{1, 2, 3}, []int{0}, 1},
		{"List of string", strFixture, []int{1}, "World,"},
		{"Get last", strFixture, []int{-1}, "Bar!"},
		{"Get before last", strFixture, []int{-2}, "Foo"},
		{"A way to before last", strFixture, []int{-12}, nil},
		{"Get nothing", strFixture, nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Get(tt.indexes...))
		})
	}
}

func Test_list_Len(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		want int
	}{
		{"Empty List", baseList{}, 0},
		{"List of int", baseList{1, 2, 3}, 3},
		{"List of string", strFixture, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Len())
			assert.Equal(t, tt.l.Len(), tt.l.Count(), "Len and Count return different values")
		})
	}
}

func Test_CreateList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []int
		want    baseIList
		wantErr error
	}{
		{"Empty", nil, baseList{}, nil},
		{"With nil elements", []int{10}, make(baseList, 10), nil},
		{"With capacity", []int{0, 10}, make(baseList, 0, 10), nil},
		{"Too much args", []int{0, 10, 1}, nil, fmt.Errorf("CreateList only accept 2 arguments, size and capacity")},
	}
	for _, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			defer func() { err = errors.Trap(err, recover()) }()
			got := baseListHelper.CreateList(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want.Cap(), got.Capacity())
		})
		if tt.wantErr == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tt.wantErr.Error())
		}
	}
}

func Test_list_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		args []int
		want baseIList
	}{
		{"Empty", nil, nil, baseList{}},
		{"Existing List", baseList{1, 2}, nil, baseList{}},
		{"With Empty spaces", baseList{1, 2}, []int{5}, baseList{nil, nil, nil, nil, nil}},
		{"With Capacity", baseList{1, 2}, []int{0, 5}, baseListHelper.CreateList(0, 5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.l.Create(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want.Cap(), got.Capacity())
		})
	}
}

func Test_list_New(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		args []interface{}
		want baseIList
	}{
		{"Empty", nil, nil, baseList{}},
		{"Existing List", baseList{1, 2}, nil, baseList{}},
		{"With elements", baseList{1, 2}, []interface{}{3, 4, 5}, baseList{3, 4, 5}},
		{"With strings", baseList{1, 2}, []interface{}{"Hello", "World"}, baseList{"Hello", "World"}},
		{"With nothing", baseList{1, 2}, []interface{}{}, baseList{}},
		{"With nil", baseList{1, 2}, nil, baseList{}},
		{"Adding array", baseList{1, 2}, []interface{}{baseList{3, 4}}, baseList{3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.New(tt.args...))
		})
	}
}

func Test_list_CreateDict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		l       baseList
		args    []int
		want    baseIDict
		wantErr error
	}{
		{"Empty", nil, nil, baseDict{}, nil},
		{"With capacity", nil, []int{10}, baseDict{}, nil},
		{"With too much parameter", nil, []int{10, 1}, nil, fmt.Errorf("CreateList only accept 1 argument for size")},
	}
	for _, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			defer func() { err = errors.Trap(err, recover()) }()
			assert.Equal(t, tt.want, tt.l.CreateDict(tt.args...))
		})
		if tt.wantErr == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tt.wantErr.Error())
		}
	}
}

func Test_list_Contains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		args []interface{}
		want bool
	}{
		{"Empty List", nil, []interface{}{}, false},
		{"Search nothing", baseList{1}, nil, true},
		{"Search nothing 2", baseList{1}, []interface{}{}, true},
		{"Not there", baseList{1}, []interface{}{2}, false},
		{"Included", baseList{1, 2}, []interface{}{2}, true},
		{"Partially there", baseList{1, 2}, []interface{}{2, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Contains(tt.args...))
			assert.Equal(t, tt.want, tt.l.Has(tt.args...))
		})
	}
}

func Test_list_First_Last(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		l         baseList
		wantFirst interface{}
		wantLast  interface{}
	}{
		{"Nil", nil, nil, nil},
		{"Empty", baseList{}, nil, nil},
		{"One element", baseList{1}, 1, 1},
		{"Many element ", baseList{1, "two", 3.1415, "four"}, 1, "four"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFirst, tt.l.First())
			assert.Equal(t, tt.wantLast, tt.l.Last())
		})
	}
}

func Test_list_Pop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		l        baseList
		args     []int
		want     interface{}
		wantList baseList
	}{
		{"Nil", nil, nil, nil, baseList{}},
		{"Empty", baseList{}, nil, nil, baseList{}},
		{"Non existent", baseList{}, []int{1}, nil, baseList{}},
		{"Empty with args", baseList{}, []int{1, 3}, baseList{nil, nil}, baseList{}},
		{"List with bad index", baseList{0, 1, 2, 3, 4, 5}, []int{1, 3, 8}, baseList{1, 3, nil}, baseList{0, 2, 4, 5}},
		{"Pop last element", baseList{0, 1, 2, 3, 4, 5}, nil, 5, baseList{0, 1, 2, 3, 4}},
		{"Pop before last", baseList{0, 1, 2, 3, 4, 5}, []int{-2}, 4, baseList{0, 1, 2, 3, 5}},
		{"Pop first element", baseList{0, 1, 2, 3, 4, 5}, []int{0}, 0, baseList{1, 2, 3, 4, 5}},
		{"Pop all", baseList{0, 1, 2, 3}, []int{0, 1, 2, 3}, baseList{0, 1, 2, 3}, baseList{}},
		{"Pop same element many time", baseList{0, 1, 2, 3}, []int{1, 1, 2, 2}, baseList{1, 1, 2, 2}, baseList{0, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotL := tt.l.Pop(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantList, gotL)
		})
	}
}

func Test_list_Intersect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		args []interface{}
		want baseList
	}{
		{"Empty List", nil, []interface{}{}, baseList{}},
		{"Intersect nothing", baseList{1}, nil, baseList{}},
		{"Intersect nothing 2", baseList{1}, []interface{}{}, baseList{}},
		{"Not there", baseList{1}, []interface{}{2}, baseList{}},
		{"Included", baseList{1, 2}, []interface{}{2}, baseList{2}},
		{"Partially there", baseList{1, 2}, []interface{}{2, 3}, baseList{2}},
		{"With duplicates", baseList{1, 2, 3, 4, 5, 4, 3, 2, 1}, []interface{}{3, 4, 5, 6, 7, 8, 7, 6, 5, 5, 4, 3}, baseList{3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Intersect(tt.args...))
		})
	}
}

func Test_list_Union(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		args []interface{}
		want baseList
	}{
		{"Empty List", nil, []interface{}{}, baseList{}},
		{"Intersect nothing", baseList{1}, nil, baseList{1}},
		{"Intersect nothing 2", baseList{1}, []interface{}{}, baseList{1}},
		{"Not there", baseList{1}, []interface{}{2}, baseList{1, 2}},
		{"Included", baseList{1, 2}, []interface{}{2}, baseList{1, 2}},
		{"Partially there", baseList{1, 2}, []interface{}{2, 3}, baseList{1, 2, 3}},
		{"With duplicates", baseList{1, 2, 3, 4, 5, 4, 3, 2, 1}, []interface{}{8, 7, 6, 5, 6, 7, 8}, baseList{1, 2, 3, 4, 5, 8, 7, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Union(tt.args...))
		})
	}
}

func Test_list_Without(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		args []interface{}
		want baseList
	}{
		{"Empty List", nil, []interface{}{}, baseList{}},
		{"Remove nothing", baseList{1}, nil, baseList{1}},
		{"Remove nothing 2", baseList{1}, []interface{}{}, baseList{1}},
		{"Not there", baseList{1}, []interface{}{2}, baseList{1}},
		{"Included", baseList{1, 2}, []interface{}{2}, baseList{1}},
		{"Partially there", baseList{1, 2}, []interface{}{2, 3}, baseList{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Without(tt.args...))
		})
	}
}

func Test_list_Unique(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		want baseList
	}{
		{"Empty List", nil, baseList{}},
		{"Remove nothing", baseList{1}, baseList{1}},
		{"Duplicates following", baseList{1, 1, 2, 3}, baseList{1, 2, 3}},
		{"Duplicates not following", baseList{1, 2, 3, 1, 2, 3, 4}, baseList{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Unique())
		})
	}
}

func Test_list_Reverse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    baseList
		want baseIList
	}{
		{"Empty List", baseList{}, baseList{}},
		{"List of int", baseList{1, 2, 3}, baseList{3, 2, 1}},
		{"List of string", strFixture, baseList{"Bar!", "Foo", "I'm", "World,", "Hello"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Clone().Reverse())
		})
	}
}

func Test_list_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		i int
		v interface{}
	}
	tests := []struct {
		name    string
		l       baseIList
		args    args
		want    baseIList
		wantErr error
	}{
		{"Empty", baseList{}, args{2, 1}, baseList{nil, nil, 1}, nil},
		{"List of int", baseList{1, 2, 3}, args{0, 10}, baseList{10, 2, 3}, nil},
		{"List of string", strFixture, args{2, "You're"}, baseList{"Hello", "World,", "You're", "Foo", "Bar!"}, nil},
		{"Negative", baseList{}, args{-1, "negative value"}, nil, fmt.Errorf("index must be positive number")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.Clone().Set(tt.args.i, tt.args.v)
			assert.Equal(t, tt.want, got)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

var mapFixture = map[string]interface{}{
	"int":     123,
	"float":   1.23,
	"string":  "Foo bar",
	"list":    []interface{}{1, "two"},
	"listInt": []int{1, 2, 3},
	"map": map[string]interface{}{
		"sub1": 1,
		"sub2": "two",
	},
	"mapInt": map[int]interface{}{
		1: 1,
		2: "two",
	},
}

var dictFixture = baseDict(baseDictHelper.AsDictionary(mapFixture).AsMap())

func dumpKeys(t *testing.T, d1, d2 baseIDict) {
	for key := range d1.AsMap() {
		v1, v2 := d1.Get(key), d2.Get(key)
		if reflect.DeepEqual(v1, v2) {
			continue
		}
		t.Logf("'%[1]v' = %[2]v (%[2]T) vs %[3]v (%[3]T)", key, v1, v2)
	}
}

func Test_dict_AsMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    baseDict
		want map[string]interface{}
	}{
		{"Nil", nil, nil},
		{"Empty", baseDict{}, map[string]interface{}{}},
		{"Map", dictFixture, map[string]interface{}(dictFixture)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.AsMap())
		})
	}
}

func Test_dict_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    baseDict
		keys []interface{}
		want baseIDict
	}{
		{"Nil", nil, nil, baseDict{}},
		{"Empty", baseDict{}, nil, baseDict{}},
		{"Map", dictFixture, nil, dictFixture},
		{"Map with Fields", dictFixture, []interface{}{"int", "list"}, baseDict(dictFixture).Omit("float", "string", "listInt", "map", "mapInt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.Clone(tt.keys...)
			assert.Equal(t, tt.want, got)

			// Ensure that the copy is distinct from the original
			got.Set("NewField", "Test")
			assert.NotEqual(t, tt.want, got)
			assert.True(t, got.Has("NewField"))
			assert.Equal(t, "Test", got.Get("NewField"))
			assert.Equal(t, tt.want.Count()+1, got.Len())
			assert.False(t, tt.d.Has("NewField"), "Has: Original dictionary has been modified")
			assert.Nil(t, tt.d.Get("NewField"), "Get: Original dictionary has been modified")
		})
	}
}

func Test_baseDict_CreateList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		d            baseDict
		args         []int
		want         baseIList
		wantLen      int
		wantCapacity int
	}{
		{"Nil", nil, nil, baseList{}, 0, 0},
		{"Empty", baseDict{}, nil, baseList{}, 0, 0},
		{"Map", dictFixture, nil, baseList{}, 0, 0},
		{"Map with size", dictFixture, []int{3}, baseList{nil, nil, nil}, 3, 3},
		{"Map with capacity", dictFixture, []int{0, 10}, baseList{}, 0, 10},
		{"Map with size&capacity", dictFixture, []int{3, 10}, baseList{nil, nil, nil}, 3, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.CreateList(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantLen, got.Len())
			assert.Equal(t, tt.wantCapacity, got.Cap())
		})
	}
}

func Test_dict_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		d       baseDict
		args    []int
		want    baseIDict
		wantErr error
	}{
		{"Empty", nil, nil, baseDict{}, nil},
		{"With capacity", nil, []int{10}, baseDict{}, nil},
		{"With too much parameter", nil, []int{10, 1}, nil, fmt.Errorf("CreateList only accept 1 argument for size")},
	}
	for _, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			defer func() { err = errors.Trap(err, recover()) }()
			assert.Equal(t, tt.want, tt.d.Create(tt.args...))
		})
		if tt.wantErr == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tt.wantErr.Error())
		}
	}
}

func Test_dict_Default(t *testing.T) {
	t.Parallel()

	type args struct {
		key    interface{}
		defVal interface{}
	}
	tests := []struct {
		name string
		d    baseDict
		args args
		want interface{}
	}{
		{"Empty", nil, args{"Foo", "Bar"}, "Bar"},
		{"Map int", dictFixture, args{"int", 1}, 123},
		{"Map float", dictFixture, args{"float", 1}, 1.23},
		{"Map Non existant", dictFixture, args{"Foo", "Bar"}, "Bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Default(tt.args.key, tt.args.defVal))
		})
	}
}

func Test_dict_Delete(t *testing.T) {
	t.Parallel()

	type args struct {
		key  interface{}
		keys []interface{}
	}
	tests := []struct {
		name    string
		d       baseDict
		args    args
		want    baseIDict
		wantErr error
	}{
		{"Empty", nil, args{}, baseDict{}, fmt.Errorf("key <nil> not found")},
		{"Map", dictFixture, args{}, dictFixture, fmt.Errorf("key <nil> not found")},
		{"Non existant key", dictFixture, args{"Test", nil}, dictFixture, fmt.Errorf("key Test not found")},
		{"Map with keys", dictFixture, args{"int", []interface{}{"list"}}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt"), nil},
		{"Map with keys + non existant", dictFixture, args{"int", []interface{}{"list", "Test"}}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt"), fmt.Errorf("key Test not found")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Clone().Delete(tt.args.key, tt.args.keys...)
			assert.Equal(t, tt.want, got)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func Test_dict_Flush(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    baseDict
		keys []interface{}
		want baseIDict
	}{
		{"Empty", nil, nil, baseDict{}},
		{"Map", dictFixture, nil, baseDict{}},
		{"Non existant key", dictFixture, []interface{}{"Test"}, dictFixture},
		{"Map with keys", dictFixture, []interface{}{"int", "list"}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt")},
		{"Map with keys + non existant", dictFixture, []interface{}{"int", "list", "Test"}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			got := d.Flush(tt.keys...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, d, got)
		})
	}
}

func Test_dict_Keys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    baseDict
		want baseIList
	}{
		{"Empty", nil, baseList{}},
		{"Map", dictFixture, baseList{str("float"), str("int"), str("list"), str("listInt"), str("map"), str("mapInt"), str("string")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.GetKeys())
		})
	}
}

func Test_dict_KeysAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    baseDict
		want strArray
	}{
		{"Empty", nil, strArray{}},
		{"Map", dictFixture, strArray{"float", "int", "list", "listInt", "map", "mapInt", "string"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.KeysAsString())
		})
	}
}

func Test_dict_Merge(t *testing.T) {
	t.Parallel()

	adding1 := baseDict{
		"int":        1000,
		"Add1Int":    1,
		"Add1String": "string",
	}
	adding2 := baseDict{
		"Add2Int":    1,
		"Add2String": "string",
		"map": map[string]interface{}{
			"sub1":   2,
			"newVal": "NewValue",
		},
	}
	type args struct {
		baseDict baseIDict
		dicts    []baseIDict
	}
	tests := []struct {
		name string
		d    baseDict
		args args
		want baseIDict
	}{
		{"Empty", nil, args{nil, []baseIDict{}}, baseDict{}},
		{"Add map to empty", nil, args{dictFixture, []baseIDict{}}, dictFixture},
		{"Add map to same map", dictFixture, args{dictFixture, []baseIDict{}}, dictFixture},
		{"Add empty to map", dictFixture, args{nil, []baseIDict{}}, dictFixture},
		{"Add new1 to map", dictFixture, args{adding1, []baseIDict{}}, dictFixture.Clone().Merge(adding1)},
		{"Add new2 to map", dictFixture, args{adding2, []baseIDict{}}, dictFixture.Clone().Merge(adding2)},
		{"Add new1 & new2 to map", dictFixture, args{adding1, []baseIDict{adding2}}, dictFixture.Clone().Merge(adding1, adding2)},
		{"Add new1 & new2 to map", dictFixture, args{adding1, []baseIDict{adding2}}, dictFixture.Clone().Merge(adding1).Merge(adding2)},
	}
	for _, tt := range tests {
		go t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Clone().Merge(tt.args.baseDict, tt.args.dicts...))
		})
	}
}

func Test_dict_Values(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    baseDict
		want baseIList
	}{
		{"Empty", nil, baseList{}},
		{"Map", dictFixture, baseList{1.23, 123, baseList{1, "two"}, baseList{1, 2, 3}, baseDict{"sub1": 1, "sub2": "two"}, baseDict{"1": 1, "2": "two"}, "Foo bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.GetValues())
		})
	}
}

func Test_dict_Pop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		d          baseDict
		args       []interface{}
		want       interface{}
		wantObject baseIDict
	}{
		{"Nil", dictFixture, nil, nil, dictFixture},
		{"Pop one element", dictFixture, []interface{}{"float"}, 1.23, dictFixture.Omit("float")},
		{"Pop missing element", dictFixture, []interface{}{"undefined"}, nil, dictFixture},
		{"Pop element twice", dictFixture, []interface{}{"int", "int", "string"}, baseList{123, 123, "Foo bar"}, dictFixture.Omit("int", "string")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			assert.Equal(t, tt.want, d.Pop(tt.args...))
			assert.Equal(t, tt.wantObject, d)
		})
	}
}

func Test_dict_Add(t *testing.T) {
	t.Parallel()

	type args struct {
		key interface{}
		v   interface{}
	}
	tests := []struct {
		name string
		d    baseDict
		args args
		want baseIDict
	}{
		{"Empty", nil, args{"A", 1}, baseDict{"A": 1}},
		{"With element", baseDict{"A": 1}, args{"A", 2}, baseDict{"A": baseList{1, 2}}},
		{"With element, another value", baseDict{"A": 1}, args{"B", 2}, baseDict{"A": 1, "B": 2}},
		{"With list element", baseDict{"A": baseList{1, 2}}, args{"A", 3}, baseDict{"A": baseList{1, 2, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Add(tt.args.key, tt.args.v))
		})
	}
}

func Test_dict_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		key interface{}
		v   interface{}
	}
	tests := []struct {
		name string
		d    baseDict
		args args
		want baseIDict
	}{
		{"Empty", nil, args{"A", 1}, baseDict{"A": 1}},
		{"With element", baseDict{"A": 1}, args{"A", 2}, baseDict{"A": 2}},
		{"With element, another value", baseDict{"A": 1}, args{"B", 2}, baseDict{"A": 1, "B": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Set(tt.args.key, tt.args.v))
		})
	}
}

func Test_dict_Transpose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    baseDict
		want baseIDict
	}{
		{"Empty", nil, baseDict{}},
		{"Base", baseDict{"A": 1}, baseDict{"1": str("A")}},
		{"Multiple", baseDict{"A": 1, "B": 2, "C": 1}, baseDict{"1": baseList{str("A"), str("C")}, "2": str("B")}},
		{"List", baseDict{"A": []int{1, 2, 3}, "B": 2, "C": 3}, baseDict{"1": str("A"), "2": baseList{str("A"), str("B")}, "3": baseList{str("A"), str("C")}}},
		{"Complex", baseDict{"A": baseDict{"1": 1, "2": 2}, "B": 2, "C": 3}, baseDict{"2": str("B"), "3": str("C"), fmt.Sprint(baseDict{"1": 1, "2": 2}): str("A")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Transpose())
		})
	}
}

func Test_baseList_Get(t *testing.T) {
	type args struct {
		indexes []int
	}
	tests := []struct {
		name string
		l    baseList
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Get(tt.args.indexes...))
		})
	}
}

func Test_baseList_TypeName(t *testing.T) {
	tests := []struct {
		name string
		l    baseList
		want str
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.TypeName())
		})
	}
}

func Test_base_TypeName(t *testing.T) {
	t.Run("list", func(t *testing.T) { assert.Equal(t, baseList{}.TypeName(), str("base")) })
	t.Run("dict", func(t *testing.T) { assert.Equal(t, baseDict{}.TypeName(), str("base")) })
}

func Test_base_GetHelper(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		gotD, gotL := baseList{}.GetHelpers()
		assert.Equal(t, gotD.CreateDictionary().TypeName(), baseDictHelper.CreateDictionary().TypeName())
		assert.Equal(t, gotL.CreateList().TypeName(), baseListHelper.CreateList().TypeName())
	})
	t.Run("dict", func(t *testing.T) {
		gotD, gotL := baseDict{}.GetHelpers()
		assert.Equal(t, gotD.CreateDictionary().TypeName(), baseDictHelper.CreateDictionary().TypeName())
		assert.Equal(t, gotL.CreateList().TypeName(), baseListHelper.CreateList().TypeName())
	})
}
