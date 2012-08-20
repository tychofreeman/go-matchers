package matchers

import (
    "testing"
)

type MockT struct {
    errArgs []string
}

func (mock *MockT) Errorf(format string, args ...interface{}) {
    mock.errArgs = append(mock.errArgs, format)
}

type Errorable interface {
    Errorf(string, ...interface{})
}

type Matcher func(interface{}) (bool, string)

func AssertThat(t Errorable, expected interface{}, m Matcher) {
    if ok, msg := m(expected); !ok {
        t.Errorf(msg)
    }
}

func TestAssertThatLogsStringWhenMatchFails(t *testing.T) {
    var mockT MockT
    AssertThat(&mockT, true, func(actual interface{}) (bool, string) {
        return false, "explanation"
    })

    if mockT.errArgs[0] != "explanation" {
        t.Errorf("Expected that matchers.AssertThat() would pass 'explanation', but it received '%v'", mockT.errArgs)
    }
}

func TestAssertThatLogsNothingWhenMatchSucceeds(t *testing.T) {
    var mockT MockT
    AssertThat(&mockT, true, func(actual interface{}) (bool, string) {
        return true, "explanation"
    })

    if len(mockT.errArgs) > 0 {
        t.Errorf("Expected that matchers.AssertThat() would not fail when matching passes, but it did.")
    }
}
