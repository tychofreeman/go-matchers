package matchers

import (
	"fmt"
	"reflect"
)

// Uses Errorf() from testing.T, so feel free to replace it with your own.
type Errorable interface {
	Errorf(string, ...interface{})
}

// Interface which supports the Equals() method.
type Equalable interface {
	Equals(interface{}) (bool, string)
}

// Short-hand type for function which is used constantly.
type Matcher func(interface{}) (bool, string)

// Define this for your type if you want to assert emptyness.
type IsEmptyable interface {
	IsEmpty() bool
}

// helper is what testing.T has to mark a function as a test helper, so
// that a failure is reported at the line that called it.
type helper interface {
	Helper()
}

// The base of this package. Other Assert* functions could be added
// in the future, likely wrappers for this one.
func AssertThat(t Errorable, actual interface{}, m Matcher) {
	if h, isHelper := t.(helper); isHelper {
		h.Helper()
	}
	if ok, msg := m(actual); !ok {
		t.Errorf("%s", msg)
	}
}

// Invert the behavior of a matcher.
func Not(m Matcher) Matcher {
	return func(actual interface{}) (bool, string) {
		ok, msg := m(actual)
		if ok && msg == "" {
			return false, "expected no match, but " + describe(actual) + " matched"
		}
		return !ok, "not " + msg
	}
}

func isBool(expected bool, actual interface{}) (bool, string) {
	switch a := actual.(type) {
	case bool:
		return a == expected, fmt.Sprintf("'%v' was expected, but got %v", expected, a)
	}
	return false, fmt.Sprintf("'%v' was expected, but got %s", expected, describe(actual))
}

// Basic truthiness
var IsTrue Matcher = func(actual interface{}) (bool, string) {
	return isBool(true, actual)
}

// Basic falsiness
var IsFalse Matcher = func(actual interface{}) (bool, string) {
	return isBool(false, actual)
}

// Test the output of a Matcher. Used in tests, but possibly in your code, too.
func HasMessage(expected string, m Matcher) Matcher {
	return func(underTest interface{}) (bool, string) {
		_, actual := m(underTest)
		return expected == actual, fmt.Sprintf("'%v' expected, but got '%v'", expected, actual)
	}
}

func equalsMsg(expected, actual interface{}) string {
	expected, _, _ = actualAndValueAndKind(expected)
	actual, _, _ = actualAndValueAndKind(actual)
	return fmt.Sprintf("expected %s, but got %s", describe(expected), describe(actual))
}

// describe is v with its type, as a message shows it: int64 42, string "a".
func describe(v interface{}) string {
	if v == nil {
		return "nil"
	}
	if s, isString := v.(string); isString {
		return fmt.Sprintf("%T %q", v, s)
	}
	return fmt.Sprintf("%T %v", v, v)
}

// Any Equals every value. Use Anything for an element whose value a test
// does not care about.
type Any struct {
}

func (a Any) Equals(other interface{}) (bool, string) {
	return true, ""
}

var __ = Any{}

// A deep equals function. Uses either the Equals method on your type, or the reflect.DeepEqual() function.
func Equals(expectedI interface{}) Matcher {

	switch expected := expectedI.(type) {
	case Equalable:
		return func(actual interface{}) (bool, string) {
			actual, _, _ = actualAndValueAndKind(actual)
			return expected.Equals(actual)
		}
	}
	return func(actual interface{}) (bool, string) {
		actual, actualValue, actualKind := actualAndValueAndKind(actual)
		if expectedI == nil && canBeNil(actualKind) {
			return actualValue.IsNil(), equalsMsg(expectedI, actual)
		}
		return reflect.DeepEqual(expectedI, actual), equalsMsg(expectedI, actual)
	}
}

// Matcher for testing emptiness of containers. Supports IsEmptyable{}, nil, arrays, channels, maps,
// slices and strings; anything else fails.
var IsEmpty Matcher = func(actualI interface{}) (bool, string) {
	switch actual := actualI.(type) {
	case nil:
		return true, ""
	case IsEmptyable:
		return actual.IsEmpty(), fmt.Sprintf("expected empty, found %v", actual)
	}
	v := reflect.ValueOf(actualI)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		return l == 0, fmt.Sprintf("expected empty, but had %d items", l)
	}
	return false, "expected a collection, but got " + describe(actualI)
}

// Returns Matcher for testing presence of an item within a container
// Supports arrays, slices, and maps (their values); anything else fails
func Contains(expected interface{}) Matcher {
	return func(actualI interface{}) (bool, string) {
		v := reflect.ValueOf(actualI)
		var items []reflect.Value
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				items = append(items, v.Index(i))
			}
		case reflect.Map:
			for values := v.MapRange(); values.Next(); {
				items = append(items, values.Value())
			}
		default:
			return false, "expected an array, slice or map, but got " + describe(actualI)
		}
		for _, item := range items {
			if eq, _ := matcherFor(expected)(item.Interface()); eq {
				return true, ""
			}
		}
		return false, fmt.Sprintf("Unable to find %v within %v", expected, actualI)
	}
}

func (m Matcher) And(other Matcher) Matcher {
	return func(actual interface{}) (bool, string) {
		passed, msg1 := m(actual)
		if !passed {
			return passed, msg1
		}
		return other(actual)
	}
}

func (m Matcher) Or(other Matcher) Matcher {
	return func(actual interface{}) (bool, string) {
		passed, msg1 := m(actual)
		if passed {
			return passed, msg1
		}
		passed, msg2 := other(actual)
		return passed, msg1 + "; or " + msg2
	}
}

// TODO: Rename this - it actually unwraps Value objects, while also returning the desired Value and Kinds
func actualAndValueAndKind(x interface{}) (interface{}, reflect.Value, reflect.Kind) {
	if x == nil {
		return nil, reflect.ValueOf(nil), reflect.Invalid
	}
	asValue, isReflectValue := x.(reflect.Value)
	if isReflectValue && asValue.CanInterface() {
		return actualAndValueAndKind(asValue.Interface())
	}
	return x, reflect.ValueOf(x), reflect.TypeOf(x).Kind()
}

// HasExactly matches a slice or array with one element for each item, in
// order. An item is a Matcher for its element, or a value its element Equals.
func HasExactly(items ...interface{}) Matcher {
	return func(actualI interface{}) (bool, string) {
		actualI, valueOfActual, kindOfActual := actualAndValueAndKind(actualI)
		if kindOfActual != reflect.Slice && kindOfActual != reflect.Array {
			return false, "expected a slice or array, but got " + describe(actualI)
		}
		if valueOfActual.Len() != len(items) {
			return false, fmt.Sprintf("expected collection of size %d, but got size %d -- expected %v vs actual %v", len(items), valueOfActual.Len(), items, actualI)
		}
		for i, item := range items {
			if ok, msg := matcherFor(item)(valueOfActual.Index(i).Interface()); !ok {
				return false, fmt.Sprintf("discrepancy at index %d - %s", i, msg)
			}
		}
		return true, ""
	}
}

func Rtns(returnValues ...interface{}) []interface{} {
	return returnValues
}

// canBeNil says whether a value of kind k can be nil.
func canBeNil(k reflect.Kind) bool {
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	}
	return false
}

// matcherFor is item if it is a Matcher, and otherwise Equals(item).
func matcherFor(item interface{}) Matcher {
	if m, isMatcher := item.(Matcher); isMatcher {
		return m
	}
	return Equals(item)
}

// Anything matches any value, as an item of HasExactly or with Equals.
var Anything = Any{}

// Ordered is the types whose values the comparison matchers can order.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// GreaterThan matches a value greater than expected.
func GreaterThan[T Ordered](expected T) Matcher {
	return compare("greater than "+describe(expected), func(a T) bool { return a > expected })
}

// LessThan matches a value less than expected.
func LessThan[T Ordered](expected T) Matcher {
	return compare("less than "+describe(expected), func(a T) bool { return a < expected })
}

// AtLeast matches a value greater than or equal to expected.
func AtLeast[T Ordered](expected T) Matcher {
	return compare("at least "+describe(expected), func(a T) bool { return a >= expected })
}

// AtMost matches a value less than or equal to expected.
func AtMost[T Ordered](expected T) Matcher {
	return compare("at most "+describe(expected), func(a T) bool { return a <= expected })
}

// Between matches a value from low to high, both included.
func Between[T Ordered](low, high T) Matcher {
	want := "between " + describe(low) + " and " + describe(high)
	return compare(want, func(a T) bool { return low <= a && a <= high })
}

// CloseTo matches a floating-point value within delta of expected.
func CloseTo[T ~float32 | ~float64](expected, delta T) Matcher {
	want := "within " + describe(delta) + " of " + describe(expected)
	return compare(want, func(a T) bool { return expected-delta <= a && a <= expected+delta })
}

// compare matches a T for which ok holds; want says what was expected. A
// value of another type fails, as it does for Equals.
func compare[T Ordered](want string, ok func(T) bool) Matcher {
	return func(actualI interface{}) (bool, string) {
		actualI, _, _ = actualAndValueAndKind(actualI)
		actual, isT := actualI.(T)
		return isT && ok(actual), "expected " + want + ", but got " + describe(actualI)
	}
}
