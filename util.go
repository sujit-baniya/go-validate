package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/gookit/filter"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/reflects"
	"github.com/gookit/goutil/strutil"
)

// NilObject represent nil value for calling functions and should be reflected at custom filters as nil variable.
type NilObject struct{}

var nilObj = NilObject{}

// init a reflect nil value
var nilRVal = reflect.ValueOf(nilObj)

// NilValue TODO a reflect nil value, use for instead of nilRVal
var NilValue = reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem())

// IsNilObj check value is internal NilObject
func IsNilObj(val interface{}) bool {
	_, ok := val.(NilObject)
	return ok
}

// CallByValue call func by reflect.Value
func CallByValue(fv reflect.Value, args ...interface{}) []reflect.Value {
	if fv.Kind() != reflect.Func {
		panicf("parameter must be an func type")
	}

	in := make([]reflect.Value, len(args))
	for k, v := range args {
		// NOTICE: reflect.Call emit panic if kind is Invalid
		if in[k] = reflect.ValueOf(v); in[k].Kind() == reflect.Invalid {
			in[k] = nilRVal
		}
	}

	// NOTICE: CallSlice()与Call() 不一样的是，参数的最后一个会被展开
	// f.CallSlice()
	return fv.Call(in)
}

func parseArgString(argStr string) (ss []string) {
	if argStr == "" { // no arg
		return
	}

	if len(argStr) == 1 { // one char
		return []string{argStr}
	}
	return stringSplit(argStr, ",")
}

// TODO strutil.Split()
func stringSplit(str, sep string) (ss []string) {
	str = strings.TrimSpace(str)
	if str == "" {
		return
	}

	for _, val := range strings.Split(str, sep) {
		if val = strings.TrimSpace(val); val != "" {
			ss = append(ss, val)
		}
	}
	return
}

// TODO use arrutil.StringsToSlice()
func strings2Args(strings []string) []interface{} {
	args := make([]interface{}, len(strings))
	for i, s := range strings {
		args[i] = s
	}
	return args
}

// TODO use arrutil.SliceToStrings()
func args2strings(args []interface{}) []string {
	strSlice := make([]string, len(args))
	for i, a := range args {
		strSlice[i] = strutil.QuietString(a)
	}
	return strSlice
}

func buildArgs(val interface{}, args []interface{}) []interface{} {
	newArgs := make([]interface{}, len(args)+1)
	newArgs[0] = val
	// as[1:] = args // error
	copy(newArgs[1:], args)

	return newArgs
}

// ValueIsEmpty check. TODO use reflects.IsEmpty()
func ValueIsEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return false // this means a bool value is not empty, whether true or false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// ValueLen get value length.
// Deprecated: please use reflects.Len()
func ValueLen(v reflect.Value) int {
	return reflects.Len(v)
}

// ErrConvertFail error
var ErrConvertFail = errors.New("convert value is failure")

func valueToInt64(v interface{}, strict bool) (i64 int64, err error) {
	switch tVal := v.(type) {
	case string:
		if strict {
			return 0, ErrConvertFail
		}
		i64, err = strconv.ParseInt(filter.Trim(tVal), 10, 0)
	case int:
		i64 = int64(tVal)
	case int8:
		i64 = int64(tVal)
	case int16:
		i64 = int64(tVal)
	case int32:
		i64 = int64(tVal)
	case int64:
		i64 = tVal
	case uint:
		i64 = int64(tVal)
	case uint8:
		i64 = int64(tVal)
	case uint16:
		i64 = int64(tVal)
	case uint32:
		i64 = int64(tVal)
	case uint64:
		i64 = int64(tVal)
	case float32:
		if strict {
			return 0, ErrConvertFail
		}
		i64 = int64(tVal)
	case float64:
		if strict {
			return 0, ErrConvertFail
		}
		i64 = int64(tVal)
	default:
		err = ErrConvertFail
	}
	return
}

// CalcLength for input value
func CalcLength(val interface{}) int {
	if val == nil {
		return -1
	}

	// return ValueLen(reflect.ValueOf(val))
	return reflects.Len(reflect.ValueOf(val))
}

// value compare.
//
// only check for: int(X), uint(X), float(X), string.
func valueCompare(srcVal, dstVal interface{}, op string) (ok bool) {
	// string compare
	if str1, ok := srcVal.(string); ok {
		str2, err := strutil.ToString(dstVal)
		if err != nil {
			return false
		}

		return strutil.VersionCompare(str1, str2, op)
	}

	// as int or float to compare
	return mathutil.Compare(srcVal, dstVal, op)
}

// getVariadicKind name.
//
// usage:
//
//	getVariadicKind(reflect.TypeOf(v))
func getVariadicKind(typ reflect.Type) reflect.Kind {
	if typ.Kind() == reflect.Slice {
		return typ.Elem().Kind()
	}
	return reflect.Invalid
}

// convTypeByBaseKind convert value type by base kind
//
//nolint:forcetypeassert
func convTypeByBaseKind(srcVal interface{}, srcKind kind, dstType reflect.Kind) (interface{}, error) {
	switch srcKind {
	case stringKind:
		switch dstType {
		case reflect.Int:
			return mathutil.ToInt(srcVal)
		case reflect.Int64:
			return mathutil.ToInt64(srcVal)
		case reflect.Bool:
			return strutil.Bool(srcVal.(string))
		case reflect.String:
			return srcVal.(string), nil
		}
	case intKind, uintKind:
		i64 := mathutil.MustInt64(srcVal)
		switch dstType {
		case reflect.Int64:
			return i64, nil
		case reflect.String:
			return strutil.ToString(srcVal)
		}
	default:
		switch dstType {
		case reflect.String:
			return strutil.ToString(srcVal)
		}
	}
	return nil, ErrConvertFail
}

// convert custom type to generic basic int, string, unit.
// returns string, int64 or error
func convToBasicType(val interface{}) (value interface{}, err error) {
	v := reflect.Indirect(reflect.ValueOf(val))

	switch v.Kind() {
	case reflect.String:
		value = v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = int64(v.Uint()) // always return int64
	default:
		err = ErrConvertFail
	}
	return
}

// From package "text/template" -> text/template/funcs.go
var (
	emptyValue = reflect.Value{}
	errorType  = reflect.TypeOf((*error)(nil)).Elem()
	// fmtStringerType  = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	// reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()
)

func panicf(format string, args ...interface{}) {
	panic("validate: " + fmt.Sprintf(format, args...))
}

func checkValidatorFunc(name string, fn interface{}) reflect.Value {
	if !goodName(name) {
		panicf("validate name %s is not a valid identifier", name)
	}

	fv := reflect.ValueOf(fn)
	if fn == nil || fv.Kind() != reflect.Func { // is nil or not is func
		panicf("validator '%s'. 2th parameter is invalid, it must be an func", name)
	}

	ft := fv.Type()
	if ft.NumIn() == 0 {
		panicf("validator '%s' func at least one parameter position", name)
	}

	if ft.NumOut() != 1 || ft.Out(0).Kind() != reflect.Bool {
		panicf("validator '%s' func must be return a bool value", name)
	}

	return fv
}

func checkFilterFunc(name string, fn interface{}) reflect.Value {
	if !goodName(name) {
		panicf("filter name %s is not a valid identifier", name)
	}

	fv := reflect.ValueOf(fn)
	if fn == nil || fv.Kind() != reflect.Func { // is nil or not is func
		panicf("filter '%s'. 2th parameter is invalid, it must be an func", name)
	}

	ft := fv.Type()
	if ft.NumIn() == 0 {
		panicf("filter '%s' func at least one parameter position", name)
	}

	if !goodFunc(ft) {
		panicf("can't install method/function %q with %d results", name, ft.NumOut())
	}

	return fv
}

// goodFunc reports whether the function or method has the right result signature.
func goodFunc(typ reflect.Type) bool {
	// We allow functions with 1 result or 2 results where the second is an error.
	switch {
	case typ.NumOut() == 1:
		return true
	case typ.NumOut() == 2 && typ.Out(1) == errorType:
		return true
	}
	return false
}

// goodName reports whether the function name is a valid identifier.
func goodName(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case r == '_':
		case i == 0 && !unicode.IsLetter(r):
			return false
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			return false
		}
	}
	return true
}

/*************************************************************
 * Comparison:
 * From package "text/template" -> text/template/funcs.go
 *************************************************************/

var errBadComparisonType = errors.New("invalid type for operation")

// errBadComparison     = errors.New("incompatible types for comparison")
// errNoComparison      = errors.New("missing argument for comparison")

type kind int

// base kinds
const (
	invalidKind kind = iota
	boolKind
	complexKind
	intKind
	floatKind
	stringKind
	uintKind
)

func basicKindV2(kind reflect.Kind) (kind, error) {
	switch kind {
	case reflect.Bool:
		return boolKind, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intKind, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintKind, nil
	case reflect.Float32, reflect.Float64:
		return floatKind, nil
	case reflect.Complex64, reflect.Complex128:
		return complexKind, nil
	case reflect.String:
		return stringKind, nil
	default:
		// like: slice, array, map ...
		return invalidKind, errBadComparisonType
	}
}

// eq evaluates the comparison a == b
func eq(arg1 reflect.Value, arg2 reflect.Value) (bool, error) {
	v1 := indirectInterface(arg1)
	k1, err := basicKindV2(v1.Kind())
	if err != nil {
		return false, err
	}

	v2 := indirectInterface(arg2)
	k2, err := basicKindV2(v2.Kind())
	if err != nil {
		return false, err
	}

	truth := false
	if k1 != k2 {
		// Special case: Can compare integer values regardless of type's sign.
		switch {
		case k1 == intKind && k2 == uintKind:
			truth = v1.Int() >= 0 && uint64(v1.Int()) == v2.Uint()
		case k1 == uintKind && k2 == intKind:
			truth = v2.Int() >= 0 && v1.Uint() == uint64(v2.Int())
			// default:
			// 	 return false, errBadComparison
		}
		return truth, nil
	}

	switch k1 {
	case boolKind:
		truth = v1.Bool() == v2.Bool()
	case complexKind:
		truth = v1.Complex() == v2.Complex()
	case floatKind:
		truth = v1.Float() == v2.Float()
	case intKind:
		truth = v1.Int() == v2.Int()
	case stringKind:
		truth = v1.String() == v2.String()
	case uintKind:
		truth = v1.Uint() == v2.Uint()
		// default:
		// 	panic("invalid kind")
	}

	return truth, nil
}

// from package: github.com/stretchr/testify/assert/assertions.go
func includeElement(list, element interface{}) (ok, found bool) {
	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	listKind := listValue.Type().Kind()

	// string contains check
	if listKind == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	defer func() {
		if e := recover(); e != nil {
			ok = false // call Value.Len() panic.
			found = false
		}
	}()

	if listKind == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if IsEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if IsEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}

	return true, false
}

/*************************************************************
 * Reflection:
 * From package(go 1.13) "reflect" -> reflect/value.go
 *************************************************************/

// IsZero reports whether v is the zero value for its type.
// It panics if the argument is invalid.
//
// NOTICE: this built-in method in reflect/value.go since go 1.13
func IsZero(v reflect.Value) bool {
	return v.IsZero()
}

// Remove type multiple pointer
func removeTypePtr(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// Remove value multiple pointer
func removeValuePtr(t reflect.Value) reflect.Value {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// ---- From package "text/template" -> text/template/exec.go

// indirect returns the item at the end of indirection, and a bool to indicate if it's nil.
// func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
// 	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
// 		if v.IsNil() {
// 			return v, true
// 		}
// 	}
// 	return v, false
// }

// indirectInterface returns the concrete value in an interface value,
// or else the zero reflect.Value.
// That is, if v represents the interface value x, the result is the same as reflect.ValueOf(x):
// the fact that x was an interface value is forgotten.
func indirectInterface(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Interface {
		return v
	}

	if v.IsNil() {
		return emptyValue
	}
	return v.Elem()
}
