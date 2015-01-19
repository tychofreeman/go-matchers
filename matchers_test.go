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

type MockEquals struct {
    UsedEqualsMethod bool
    Expected bool
}

func (e *MockEquals) Equals(other interface{}) (bool, string) {
    e.UsedEqualsMethod = true
    return e.Expected, "formatted equals message"
}

func TestUsesEqualsIfExists(t *testing.T) {
    e := new(MockEquals)
    Equals(e)(false)

    AssertThat(t, e.UsedEqualsMethod, IsTrue)
}

func TestUsesEqualSignIfNotEuqalable(t *testing.T) {
    AssertThat(t, true, Equals(true))
}

func TestCanCompareArbitraryTypes(t *testing.T) {
    AssertThat(t, Ignored{}, Equals(Ignored{}))
}

func TestCatchesDifferencesInArbitraryTypes(t *testing.T) {
    AssertThat(t, MockEquals{Expected: true}, Not(Equals(MockEquals{})))
}

func TestEmptySpliceIsEmpty(t *testing.T) {
    AssertThat(t, make([]int, 0), IsEmpty)
}

func TestNonEmptyArrayIsNotEmpty(t *testing.T) {
    nonEmpty := []bool{ true, false, true }
    AssertThat(t, nonEmpty, Not(IsEmpty))
}

func TestEmptyMapIsEmpty(t *testing.T) {
    empty := make(map[string]string)
    AssertThat(t, empty, IsEmpty)
}

func TestNonEmptyMapIsNotEmpty(t *testing.T) {
    nonEmpty := make(map[string]string)
    nonEmpty["abc"] = "def"
    AssertThat(t, nonEmpty, Not(IsEmpty))
}

type MockIsEmpty bool
func (m MockIsEmpty) IsEmpty() bool {
    return bool(m)
}

func TestIsEmptyIsUsedWhenFound(t *testing.T) {
    empty := MockIsEmpty(true)
    AssertThat(t, empty, IsEmpty)

    notEmpty := MockIsEmpty(false)
    AssertThat(t, notEmpty, Not(IsEmpty))
}

func TestIsEmptyPanicsForArbitreryTypes(t *testing.T) {
    defer func() {
        e := recover()
        if e == nil {
            t.Errorf("Expected to panic when given non-emptable type.")
        }
    }()

    AssertThat(t, 0, IsEmpty)
}

func TestContainsDoesntMatchEmptyIntegerSlice(t *testing.T) {
    list := []int{ }
    AssertThat(t, list, Not(Contains(5)))
    AssertThat(t, list, Not(Contains(10)))
}

func TestContainsMatchesIntegersInAList(t *testing.T) {
    list := []int{ 5, 10 }
    AssertThat(t, list, Contains(5))
    AssertThat(t, list, Contains(10))
}

func TestAndMatchesBothPredicates(t *testing.T) {
    AssertThat(t, "", IsEmpty.And(Equals("")))
}

func TestEqualsHandlesNilExpectations(t *testing.T) {
    result, _ := Equals(nil)("")
    AssertThat(t, result, IsFalse)
}

func TestEqualsHandlesNilActuals(t *testing.T) {
    result, _ := Equals("")(nil)
    AssertThat(t, result, IsFalse)
}

func TestNilEqualsNil(t *testing.T) {
    result, _ := Equals(nil)(nil)
    AssertThat(t, result, IsTrue)
}

func TestEmptyHasExactlyListMatchesEmptySplice(t *testing.T) {
    AssertThat(t, []interface{}{}, HasExactly())
}

func TestHasExactlyTrueDoesNotMatchEmptySlice(t *testing.T) {
    AssertThat(t, []interface{}{}, Not(HasExactly(true)))
}

func TestHasExactlyDoesNotMatchScalarValues(t *testing.T) {
    AssertThat(t, 1, Not(HasExactly(1)))
}

func TestHasExactlyMatchesEachElementInTurn(t *testing.T) {
    AssertThat(t, []interface{}{true}, Not(HasExactly(false)))
}

func TestAnyMatchesAnything(t *testing.T) {
    AssertThat(t, true, Equals(__))
    AssertThat(t, false, Equals(__))
    AssertThat(t, 1, Equals(__))
    AssertThat(t, 0, Equals(__))
    AssertThat(t, []interface{}{}, Equals(__))
}

func TestAnyCanBeUsedInHasExactly(t *testing.T) {
    AssertThat(t, []interface{}{true}, HasExactly(__))
}

func TestAnyCannotBeProductivelyUsedAsActual(t *testing.T) {
    AssertThat(t, []interface{}{__}, Not(HasExactly(true)))
}

var testMatcher Matcher = func(actual interface{}) (bool,string) {
    return true, ""
}

func TestMatchersCanBeUsedInHasExactly(t *testing.T) {
    AssertThat(t, []interface{}{1}, HasExactly(testMatcher))
}

func TestArbitrarySlicesCanBeUsedInHasExactly(t *testing.T) {
    AssertThat(t, []int{1}, HasExactly(1))
}

func TestCanPassVariadicFunctions(t *testing.T) {
    var f = func() (string,bool) { return "", false }
    AssertThat(t, Rtns(f()), HasExactly("", false))
}

type TestMock struct {}
func (tm TestMock) Test(in string) bool {
    return "a" == in
}
/*
func TestMethod_WasCalledWith_MatchesParams(t *testing.T) {
    m := TestMock{}
    m.Test("a")
    AssertThat(t, m, Method("Test").WasCalledWith(__))
    AssertThat(t, m, Not(Method("Test").WasCalledWith("b")))
}

*/

