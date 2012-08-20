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

func 
