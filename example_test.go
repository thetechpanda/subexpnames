package subexpnames_test

import (
	"fmt"
	"regexp"

	"github.com/thetechpanda/subexpnames"
)

func ExampleMatch() {
	// a regular expression to match a date
	re := regexp.MustCompile(`(?P<overlap>(?P<year>(?P<thousands>\d)(?P<hundreds>\d)(?P<tens>\d)(?P<ones>\d))-(?P<month>(?P<tens>\d)(?P<ones>\d)))-(?P<overlap>(?P<day>(?P<tens>\d)(?P<ones>\d)))`)
	// a subject to match the regular expression
	subject := "this is a test subject to see if we can parse 2016-01-02 and 1234-56-78 using the Match() function."

	match, ok := subexpnames.Match(re, subject)
	if !ok {
		return
	}

	keys := []string{"overlap", "year"}
	if v, ok := match.Get(0, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year] = 2016
	}
	if v, ok := match.Get(1, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year] = 1234
	}

	keys = []string{"overlap", "year", "thousands"}
	if v, ok := match.Get(0, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year thousands] = 2
	}
	if v, ok := match.Get(1, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year thousands] = 1
	}
	keys = []string{"overlap", "year", "hundreds"}
	if v, ok := match.Get(0, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year hundreds] = 0
	}
	if v, ok := match.Get(1, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year hundreds] = 2
	}
	keys = []string{"overlap", "year", "tens"}
	if v, ok := match.Get(0, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year tens] = 1
	}
	if v, ok := match.Get(1, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year tens] = 3
	}
	keys = []string{"overlap", "year", "ones"}
	if v, ok := match.Get(0, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year ones] = 6
	}
	if v, ok := match.Get(1, 0, keys...); ok {
		fmt.Println(keys, "=", v) // prints [overlap year ones] = 4
	}
}
