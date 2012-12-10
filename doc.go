// Copyright 2012 Christopher Freeman
// Use of this file is governed by a BSD-like license
// which can be found in the LICENSE file.

/*
  Package matchers is an implementation of Hamcrest-like matchers for use in Go tests. You can use it to simplify your Go tests.
  Currently, only Errorf() is supported, so if your tests require a short-circuiting assertion, please fix your tests. Or, you know, submit a pull request.

  This package is used to change this:
     import "my_boolean_pkg"
     import "fmt"
     import "testing"
     func TestTrueIsTrue(t *testing.T) {
        if true != my_boolean_pkg.True {
            t.Errorf(fmt.Sprintf("Help! My True value is not true!\n  Actual: %v\n  Expected: %v\n", true, my_boolean_pkg.True))
        }
     }

  into the more readable:

     import "my_boolean_pkg"
     import "fmt"
     import "testing"
     import . "github.com/tychofreeman/go-matchers"
     func TestTrueIsTrue(t *testing.T) {
        AssertTrue(t, my_boolean_pkg.True, IsTrue)
     }

  The latter will output a (usually) highly readable message, which is consistent across your project. (I find creating error messages to be a bit of a pain.)

  If you wish to have your own text, you may create a Matcher specific to your needs. For example:
     // This could be so much better with structural typing instead of reflection...
     func IsOdd(actual interface{}) (bool,string) {
         actualValue := reflect.ValueOf(actual)
         switch actualValue.Kind() {
             case reflect.Int:   fallthrough
             case reflect.Int8:  fallthrough
             case reflect.Int16: fallthrough
             case reflect.Int32: fallthrough
             case reflect.Int64:
                 return (actualValue.Int() % 2 == 0), fmt.Sprintf("expecting an odd integer, but got %v (signed)", actualValue.Int())
             case reflect.Uint:   fallthrough
             case reflect.Uint8:  fallthrough
             case reflect.Uint16: fallthrough
             case reflect.Uint32: fallthrough
             case reflect.Uint64:
                 return (actualValue.Uint() % 2 == 0), fmt.Sprintf("expecting an odd integer, but got %v (unsigned)", actualValue.Uint())
         }
         return false, fmt.Sprintf("expecting an odd integer, but got non-boolean of type %s", reflect.TypeOf(actual).Name())
     }

  It is common for a Matcher to be implemented as a closure around an expected value, like:
     func HasColor(expected Color) Matcher {
        return func(actual interface{}) (bool,string) {
            switch a := actual.(type) {
                case Pixel:
                    return a.Color() == expected, fmt.Sprintf("expecting color %v but got '%v'", expected, actual)
            }
            return false, fmt.Sprintf("expecting the color %v, but got a non-color thing of type %s", expected, reflect.TypeOf(actual).Name())
        }
     }

  It can be useful to combine Matchers. For example, there is a Not() method which takes a Matcher in and returns a modified Matcher. It is used to invert the meaning of your Matcher, and the output attempts to reflect that by prepending the word 'not ' to the error message.

  Future Extensions:
  
  * It would be nice to be able to enhance the current error messages with test-specific text, such as this GTest snippet:
     ASSERT_EQ(true, false) << "This should never pass!"
  * Other Matcher-modifying methods would be usefule: And, Or, Unless, etc.

  Hopefully, you find your tests easier to write and think about by using this package.

*/
package matchers
