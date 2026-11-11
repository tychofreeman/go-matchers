go_matchers
===========

Hamcrest-style matchers for unit testing in Go.

```go
import (
	"testing"

	. "github.com/tychofreeman/go-matchers"
)

func TestTotals(t *testing.T) {
	AssertThat(t, total(cart), Equals(42))
	AssertThat(t, cart.Items, HasExactly("apple", Anything, "pear"))
	AssertThat(t, cart.Coupons, Contains("SALE").Or(IsEmpty))
	AssertThat(t, ratio(cart), CloseTo(0.25, 1e-9))
	AssertThat(t, len(cart.Items), Between(1, 10))
}
```

A failure is reported at the line of the test, with what was expected and
what was got: `expected int 42, but got int 41`.

## Matchers

- `Equals(v)`: `reflect.DeepEqual`, or `v.Equals` for an `Equalable`; `Equals(nil)` matches any nil.
- `IsTrue`, `IsFalse`
- `IsEmpty`: arrays, slices, maps, channels, strings, nil, or an `IsEmptyable`.
- `Contains(v)`: an element of an array or slice, or a value of a map, that equals `v` or matches it if it is a Matcher.
- `HasExactly(items...)`: a slice or array whose elements match the items in order; `Anything` matches any element.
- `GreaterThan(v)`, `LessThan(v)`, `AtLeast(v)`, `AtMost(v)`, `Between(low, high)`: numbers and strings, of the same type as `v`.
- `CloseTo(v, delta)`: floating-point values within `delta` of `v`.
- `Not(m)`, `m.And(other)`, `m.Or(other)`, `HasMessage(msg, m)`

A Matcher is a `func(interface{}) (bool, string)`, so writing your own is a
function that says whether the value matches and, if not, why.