// Implementation of Hamcrest-like matchers for use in Go tests
package matchers

import (
    "fmt"
    "reflect"
)

// Uses Errorf() from testing.T, so feel free to replace it with your own
type Errorable interface {
    Errorf(string, ...interface{})
}

type Matcher func(interface{}) (bool, string)

// The base of this package. Other Assert* functions could be added
// in the future, likely wrappers for this one.
func AssertThat(t Errorable, expected interface{}, m Matcher) {
    if ok, msg := m(expected); !ok {
        t.Errorf(msg)
    }
}

// Invert the behavior of a matcher
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
func IsTrue(actual interface{}) (bool, string) {
    return isBool(true, actual)
}

// Basic falsiness
func IsFalse(actual interface{}) (bool, string) {
    return isBool(false, actual)
}

// Test the output of a Matcher. Used in tests, but possibly in your code, too.
func HasMessage(expected string, m Matcher) Matcher {
    return func (underTest interface{}) (bool, string) {
        _, actual := m(underTest)
        return expected == actual, fmt.Sprintf("'%v' expected, but got '%v'", expected, actual)
    }
}

// Interface which supports the Equals() method.
type Equalable interface {
    Equals(interface{}) (bool, string)
}

func equalsMsg(expected, actual interface{}) string {
    return fmt.Sprintf("'%v%v' expected, but got '%v%v'", 
            reflect.TypeOf(expected).Name(), expected,
            reflect.TypeOf(actual).Name(), actual)
}

// A deep equals function. Uses either the Equals method on your type, or the reflect.DeepEqual() function.
func Equals(expectedI interface{}) Matcher {

    switch expected := expectedI.(type) {
        case Equalable:
            return func (actual interface{}) (bool, string) {
                return expected.Equals(actual)
            }
    }
    return func (actual interface{}) (bool, string) {
        return reflect.DeepEqual(expectedI, actual), equalsMsg(expectedI, actual)
    }
}

type IsEmptyable interface {
    IsEmpty() bool
}

// Matcher for testing emptiness of containers
func IsEmpty(actualI interface{}) (bool, string) {
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
