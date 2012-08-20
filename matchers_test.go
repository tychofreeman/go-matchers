package matchers

import (
    "testing"
)


// Mock part of testing.T
type MockT struct {
    errArgs []string
}

func (mock *MockT) Errorf(format string, args ...interface{}) {
    mock.errArgs = append(mock.errArgs, format)
}

// Document that we don't care about this parameter.
type Ignored struct {
}

// Create a Matcher function with the given return values
func ConstMatcher(matches bool, msg string) Matcher {
    return func(actual interface{}) (bool, string) {
        return matches, msg
    }
}

func TestAssertThatLogsStringWhenMatchFails(t *testing.T) {
    var mockT MockT
    AssertThat(&mockT, Ignored{}, ConstMatcher(false, "explanation"))

    if mockT.errArgs[0] != "explanation" {
        t.Errorf("Expected that matchers.AssertThat() would pass 'explanation', but it received '%v'", mockT.errArgs)
    }
}

func TestAssertThatLogsNothingWhenMatchSucceeds(t *testing.T) {
    var mockT MockT
    AssertThat(&mockT, Ignored{}, ConstMatcher(true, "explanation"))

    if len(mockT.errArgs) > 0 {
        t.Errorf("Expected that matchers.AssertThat() would not fail when matching passes, but it did.")
    }
}

func TestNotInvertsFalseMatcher(t *testing.T) {
    AssertThat(t, Ignored{}, Not(ConstMatcher(false, "")))
}

func TestNotNotKeepsTrueMatcher(t *testing.T) {
    AssertThat(t, Ignored{}, Not(Not(ConstMatcher(true, ""))))
}

func TestNotAddsNotToMessage(t *testing.T) {
    var mockT MockT
    AssertThat(&mockT, Ignored{}, Not(ConstMatcher(true, "bob")))
    if mockT.errArgs[0] != "not bob" {
        t.Errorf("Expected Not() to add 'not ' to error message, but the actual message was '%v'", mockT.errArgs)
    }
}

func TestFalseIsTrueIsFalse(t *testing.T) {
    AssertThat(t, false, Not(IsTrue))
}

func TestTrusIsTrueIsTrue(t *testing.T) {
    AssertThat(t, true, IsTrue)
}

func TestFalseIsFalseIsTrue(t *testing.T) {
    AssertThat(t, false, IsFalse)
}

func TestTrusIsFalseIsFalse(t *testing.T) {
    AssertThat(t, true, Not(IsFalse))
}

func TestRecognizesConstantMessage(t *testing.T) {
    AssertThat(t, Ignored{}, HasMessage("test message", ConstMatcher(false, "test message")))
}

func TestFailsIfMessageIsNotGiven(t *testing.T) {
    AssertThat(t, Ignored{}, Not(HasMessage("not found", ConstMatcher(false, "this is the real msg"))))
}
