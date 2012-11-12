package matchers

import (
    "fmt"
    "reflect"
)

type Errorable interface {
    Errorf(string, ...interface{})
}

type Matcher func(interface{}) (bool, string)

func AssertThat(t Errorable, expected interface{}, m Matcher) {
    if ok, msg := m(expected); !ok {
        t.Errorf(msg)
    }
}

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

func IsTrue(actual interface{}) (bool, string) {
    return isBool(true, actual)
}

func IsFalse(actual interface{}) (bool, string) {
    return isBool(false, actual)
}

func HasMessage(expected string, m Matcher) Matcher {
    return func (underTest interface{}) (bool, string) {
        _, actual := m(underTest)
        return expected == actual, fmt.Sprintf("'%v' expected, but got '%v'", expected, actual)
    }
}

type Equalable interface {
    Equals(interface{}) (bool, string)
}

func EqualsMsg(expected, actual interface{}) string {
    return fmt.Sprintf("'%v%v' expected, but got '%v%v'", 
            reflect.TypeOf(expected).Name(), expected,
            reflect.TypeOf(actual).Name(), actual)
}

func Equals(expectedI interface{}) Matcher {

    switch expected := expectedI.(type) {
        case Equalable:
            return func (actual interface{}) (bool, string) {
                return expected.Equals(actual)
            }
    }
    return func (actual interface{}) (bool, string) {
        return expectedI == actual, EqualsMsg(expectedI, actual)
    }
}

type IsEmptyable interface {
    IsEmpty() bool
}

func IsEmpty(actualI interface{}) (bool, string) {
    switch actual := actualI.(type) {
        case IsEmptyable:
            return actual.IsEmpty(), fmt.Sprintf("expected empty, found %v", actual)
    }
    v := reflect.ValueOf(actualI)
    l := v.Len()
    return l == 0, fmt.Sprintf("expected empty, but had %d items", l)
}
