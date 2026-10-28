package matchers

import (
	"fmt"
	"testing"
)

// Mock part of testing.T
type MockT struct {
	errArgs []string
}

func (mock *MockT) Errorf(format string, args ...interface{}) {
	mock.errArgs = append(mock.errArgs, fmt.Sprintf(format, args...))
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
	Expected         bool
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
	nonEmpty := []bool{true, false, true}
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

func TestIsEmptyFailsForArbitraryTypes(t *testing.T) {
	AssertThat(t, 0, Not(IsEmpty))
}

func TestContainsDoesntMatchEmptyIntegerSlice(t *testing.T) {
	list := []int{}
	AssertThat(t, list, Not(Contains(5)))
	AssertThat(t, list, Not(Contains(10)))
}

func TestContainsMatchesIntegersInAList(t *testing.T) {
	list := []int{5, 10}
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

var testMatcher Matcher = func(actual interface{}) (bool, string) {
	return true, ""
}

func TestMatchersCanBeUsedInHasExactly(t *testing.T) {
	AssertThat(t, []interface{}{1}, HasExactly(testMatcher))
}

type SomeSlice []interface{}

func TestTypeAliaseCanBeUsedInHasExatly(t *testing.T) {
	AssertThat(t, SomeSlice{1}, HasExactly(1))
}

func TestArbitrarySlicesCanBeUsedInHasExactly(t *testing.T) {
	AssertThat(t, []int{1}, HasExactly(1))
}

func TestCanPassVariadicFunctions(t *testing.T) {
	var f = func() (string, bool) { return "", false }
	AssertThat(t, Rtns(f()), HasExactly("", false))
}

type TestMock struct{}

func (tm TestMock) Test(in string) bool {
	return "a" == in
}

type List []interface{}

func TestAliasesOfSliceCanBeUsedInExactly(t *testing.T) {
	AssertThat(t, List{"a", List{"b"}, "c"}, HasExactly("a", HasExactly("b"), Equals("c")))
}

/*
func TestMethod_WasCalledWith_MatchesParams(t *testing.T) {
    m := TestMock{}
    m.Test("a")
    AssertThat(t, m, Method("Test").WasCalledWith(__))
    AssertThat(t, m, Not(Method("Test").WasCalledWith("b")))
}

*/

// recorder is an Errorable that keeps the messages it is given, and
// whether it was marked as a test helper.
type recorder struct {
	messages []string
	helped   bool
}

func (r *recorder) Errorf(format string, args ...interface{}) {
	r.messages = append(r.messages, fmt.Sprintf(format, args...))
}

func (r *recorder) Helper() { r.helped = true }

// AssertThat marks itself a helper, so that go test reports a failure at
// the line of the test that called it.
func TestAssertThatIsATestHelper(t *testing.T) {
	r := &recorder{}
	AssertThat(r, 1, Equals(2))
	AssertThat(t, r.helped, IsTrue)
}

func TestAssertThatReportsMessagesAsTheyAre(t *testing.T) {
	r := &recorder{}
	AssertThat(r, "50% off", Equals("60% off"))
	AssertThat(t, r.messages, HasExactly(`expected string "60% off", but got string "50% off"`))
}

func TestEqualsNamesTheTypeAndTheValue(t *testing.T) {
	AssertThat(t, int64(42), HasMessage("expected int64 43, but got int64 42", Equals(int64(43))))
	AssertThat(t, 1, HasMessage("expected int64 1, but got int 1", Equals(int64(1))))
	AssertThat(t, nil, HasMessage("expected int 1, but got nil", Equals(1)))
}

func TestHasExactlyNamesTheTypeAndTheValue(t *testing.T) {
	AssertThat(t, []int{1, 2}, HasMessage("discrepancy at index 1 - expected int 3, but got int 2", HasExactly(1, 3)))
}

func TestIsEmptyFailsForWhatIsntACollection(t *testing.T) {
	AssertThat(t, 5, Not(IsEmpty))
	AssertThat(t, 5, HasMessage("expected a collection, but got int 5", IsEmpty))
	AssertThat(t, nil, IsEmpty)
	AssertThat(t, "", IsEmpty)
}

func TestContainsFailsForWhatIsntACollection(t *testing.T) {
	AssertThat(t, 5, Not(Contains(5)))
	AssertThat(t, 5, HasMessage("expected an array, slice or map, but got int 5", Contains(5)))
	AssertThat(t, nil, Not(Contains(5)))
}

func TestContainsLooksAtAMapsValues(t *testing.T) {
	AssertThat(t, map[string]int{"a": 1}, Contains(1))
	AssertThat(t, map[string]int{"a": 1}, Not(Contains("a")))
}
