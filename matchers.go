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

// The base of this package. Other Assert* functions could be added
// in the future, likely wrappers for this one.
func AssertThat(t Errorable, expected interface{}, m Matcher) {
    if ok, msg := m(expected); !ok {
        t.Errorf(msg)
    }
}

// Invert the behavior of a matcher.
func Not(m Matcher) Matcher {
    return func(actual interface{}) (bool, string) {
        ok, msg := m(actual)
        return !ok, "not " + msg
    }
}

func isBool(expected bool, actual interface{}) (bool, string) {
    switch a := actual.(type) {
        case bool:
            return a == expected, fmt.Sprintf("'%v' was expected, but got %v", expected, a)
    }
    return false, fmt.Sprintf("'%v' was expected, but got non-boolean of type %s", expected, reflect.TypeOf(actual).Name())
}

// Basic truthiness
var IsTrue Matcher = func (actual interface{}) (bool, string) {
    return isBool(true, actual)
}

// Basic falsiness
var IsFalse Matcher = func (actual interface{}) (bool, string) {
    return isBool(false, actual)
}

// Test the output of a Matcher. Used in tests, but possibly in your code, too.
func HasMessage(expected string, m Matcher) Matcher {
    return func (underTest interface{}) (bool, string) {
        _, actual := m(underTest)
        return expected == actual, fmt.Sprintf("'%v' expected, but got '%v'", expected, actual)
    }
}

func equalsMsg(expected, actual interface{}) string {
    var expectedType, actualType string
    if expected != nil {
        expectedType = reflect.TypeOf(expected).Name()
    }
    if actual != nil {
        actualType = reflect.TypeOf(actual).Name()
    }
    return fmt.Sprintf("'%v%v' expected, but got '%v%v'", 
            expectedType, expected,
            actualType, actual)
}

type Any struct {
}
func (a Any) Equals(other interface{}) (bool,string) {
    return true, ""
}
var __ = Any{}

// A deep equals function. Uses either the Equals method on your type, or the reflect.DeepEqual() function.
func Equals(expectedI interface{}) Matcher {

    switch expected := expectedI.(type) {
        case Equalable:
            return func (actual interface{}) (bool, string) {
                actual, _, _ = actualAndValueAndKind(actual)
                return expected.Equals(actual)
            }
    }
    return func (actual interface{}) (bool, string) {
        actual, actualValue, actualKind := actualAndValueAndKind(actual)
        if actualKind == reflect.Ptr && expectedI == nil {
            return actualValue.IsNil(), equalsMsg(expectedI, actual)
        }
        return reflect.DeepEqual(expectedI, actual), equalsMsg(expectedI, actual)
    }
}

// Matcher for testing emptiness of containers. Supports IsEmptyable{}, arrays, slices, and maps.
var IsEmpty Matcher = func (actualI interface{}) (bool, string) {
    switch actual := actualI.(type) {
        case IsEmptyable:
            return actual.IsEmpty(), fmt.Sprintf("expected empty, found %v", actual)
    }
    v := reflect.ValueOf(actualI)
    l := v.Len()
    return l == 0, fmt.Sprintf("expected empty, but had %d items", l)
}

// Returns Matcher for testing presence of an item within a container
// Supports arrays, slices, and maps
func Contains(expected interface{}) Matcher {
    return func(actualI interface{}) (bool, string) {
        v := reflect.ValueOf(actualI)
        l := v.Len()
        for i := 0; i < l; i++ {
            if eq, _ := Equals(expected)(v.Index(i).Interface()); eq == true {
                return true, ""
            }
        }
        return false, fmt.Sprintf("Unable to find %v within %v", expected, actualI)
    }
}

func (m Matcher) And(other Matcher) Matcher {
    return func(actual interface{}) (bool,string) {
        passed, msg1 := m(actual)
        if !passed {
            return passed, msg1
        }
        return other(actual)
    }
}

func (m Matcher) Or(other Matcher) Matcher {
    return func(actual interface{}) (bool,string) {
        if passed, msg1 := m(actual); passed {
            return passed, msg1
        }
        return other(actual)
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

func HasExactly(items ...interface{}) Matcher {
    return func(actualI interface{}) (bool,string) {
        if actualI == nil {
            return false, fmt.Sprintf("Got 'nil' instead of expected %v", items)
        }
        actualI, valueOfActual, kindOfActual := actualAndValueAndKind(actualI)
        if kindOfActual == reflect.Slice || kindOfActual == reflect.Array {
            lenOfActual := valueOfActual.Len()
            lenOfItems := len(items)
            hasSameLen := (lenOfItems == lenOfActual)
            if !hasSameLen {
                return false, fmt.Sprintf("expected collection of size %d, but got size %d -- expected [%v] vs actual [%v]", lenOfItems, lenOfActual, items, actualI)
            }
            for i := 0; i < lenOfActual; i++ {
                switch t := items[i].(type) {
                case Matcher:
                    switch s := actualI.(type) {
                    case []interface{}:
                        if result, msg := t(s[i]); !result {
                            return result, msg
                        }
                    default:
                        if result, msg := t(valueOfActual.Index(i)); !result {
                            return result, msg
                        }
                    }
                case Equalable:
                    if result, msg := t.Equals(valueOfActual.Index(i)); !result {
                        return result, msg
                    }
                default:
                    if items[i] != valueOfActual.Index(i).Interface() {
                        return false, fmt.Sprintf("discrepancy at index %d - %s", i, equalsMsg(items[i], valueOfActual.Index(i)))
                    }
                }
            }
        } else {
            return false, fmt.Sprintf("HasExactly() matcher requires a collection, found %T (size %v - value %v)", actualI, valueOfActual.Type().Size(), actualI)
        }
        return true, ""
    }
}

func Rtns(returnValues ...interface{}) []interface{} {
    return returnValues
}

type MethodMatcherGenerator struct {
    methodName string
}

type Call struct {
    name string
    args []interface{}
}

type Mock interface {
    Calls() []Call
}

func (mmg MethodMatcherGenerator) WasCalledWith(args ...interface{}) Matcher {
    return func(actualI interface{}) (bool,string) {
        switch actualI.(type) {
            case Mock:
                return true, ""
            default:
                return false, fmt.Sprintf("expected method '%s' was called with args %v", mmg.methodName, args)
        }
        panic("switch failed - should have found a valid case in preceeding switch.")
    }
}

func Method(name string) MethodMatcherGenerator {
    return MethodMatcherGenerator{name}
}
