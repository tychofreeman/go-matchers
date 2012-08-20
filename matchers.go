package matchers

type Errorable interface {
    Errorf(string, ...interface{})
}

type Matcher func(interface{}) (bool, string)

func AssertThat(t Errorable, expected interface{}, m Matcher) {
    if ok, msg := m(expected); !ok {
        t.Errorf(msg)
    }
}
