package subexpnames_test

import (
	"regexp"
	"slices"
	"testing"

	"github.com/thetechpanda/subexpnames"
)

func expectValue(t *testing.T, match *subexpnames.Matches, index int, keys []string, value string) {
	if v, ok := match.Get(index, 0, keys...); !ok || v != value {
		t.Fatalf("%v: expected %q, got %q", keys, value, v)
	}
}

func expectValues(t *testing.T, match *subexpnames.Matches, index int, keys []string, values []string) {
	if matches, ok := match.GetAll(index, keys...); ok {
		if slices.Compare(matches, values) != 0 {
			t.Fatalf("%v: expected %q, got %q", keys, values, matches)
		}
		return
	}
	t.Fatalf("%v: not found", keys)
}

func TestSubExpNamed(t *testing.T) {
	// a regular expression to match a date
	re := regexp.MustCompile(`(?P<overlap>(?P<year>(?P<thousands>\d)(?P<hundreds>\d)(?P<tens>\d)(?P<ones>\d))-(?P<month>(?P<tens>\d)(?P<ones>\d)))-(?P<overlap>(?P<day>(?P<tens>\d)(?P<ones>\d)))`)

	_, ok := subexpnames.Match(re, "not a match")
	if ok {
		t.Fatalf("expected not a match")
	}

	// a subject to match the regular expression
	subject := "this is a test subject to see if we can parse 2016-01-02 and 1234-56-78 using the Match() function."

	match, ok := subexpnames.Match(re, subject)
	if !ok {
		t.Fatalf("expected a match")
	}

	if len(*match) != 2 || match.Len() != 2 {
		t.Fatalf("expected 2 match, got %d", len(*match))
	}

	if len((*match)[0].Nested) != 2 {
		t.Fatalf("expected 2 nested values, got %d", len((*match)[0].Nested))
	}
	if len((*match)[1].Nested) != 2 {
		t.Fatalf("expected 2 nested values, got %d", len((*match)[1].Nested))
	}

	k := match.Keys(3)
	if len(k) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(k))
	}

	k = match.Keys(0)
	expected := [][]string{
		{"overlap"},
		{"overlap", "year"},
		{"overlap", "year", "thousands"},
		{"overlap", "year", "hundreds"},
		{"overlap", "year", "tens"},
		{"overlap", "year", "ones"},
		{"overlap", "month"},
		{"overlap", "month", "tens"},
		{"overlap", "month", "ones"},
		{"overlap", "day"},
		{"overlap", "day", "tens"},
		{"overlap", "day", "ones"},
	}
	for i, v := range k {
		if slices.Compare(v, expected[i]) != 0 {
			t.Fatalf("expected %v, got %v", expected[i], v)
		}
	}

	if _, ok := match.GetAll(3); ok {
		t.Fatalf("expected not found")
	}

	if _, ok := match.GetGroup(3); ok {
		t.Fatalf("expected not found")
	}

	if v, ok := match.GetAll(0); !ok {
		t.Fatalf("expected found")
	} else if len(v) != 1 {
		t.Fatalf("expected 1 value, got %d", len(v))
	} else if v[0] != "2016-01-02" {
		t.Fatalf("expected 2016-01-02, got %s", v[0])
	}

	if v, ok := match.GetFirstValueOfGroup(0, "overlap", "year"); !ok {
		t.Fatalf("expected a value")
	} else if v != "2016" {
		t.Fatalf("expected 2016, got %s", v)
	}

	if m, ok := match.GetGroup(0); !ok {
		t.Fatalf("expected a group")
	} else if len(m.Nested) != 2 {
		t.Fatalf("expected 2 nested values, got %d", len(m.Nested))
	} else if m.Key != "" {
		t.Fatalf("expected an empty key, got %s", m.Key)
	}

	if _, ok := match.Get(0, 0, "not-found"); ok {
		t.Fatalf("expected not found")
	}

	if _, ok := match.Get(0, 6, "overlap"); ok {
		t.Fatalf("expected not found")
	}

	expectValues(t, match, 0, []string{"overlap"}, []string{"2016-01", "02"})
	expectValues(t, match, 1, []string{"overlap"}, []string{"1234-56", "78"})

	expectValue(t, match, 0, []string{}, "2016-01-02")
	expectValue(t, match, 0, []string{"overlap"}, "2016-01")
	expectValue(t, match, 0, []string{"overlap", "year"}, "2016")
	expectValue(t, match, 0, []string{"overlap", "year", "thousands"}, "2")
	expectValue(t, match, 0, []string{"overlap", "year", "hundreds"}, "0")
	expectValue(t, match, 0, []string{"overlap", "year", "tens"}, "1")
	expectValue(t, match, 0, []string{"overlap", "year", "ones"}, "6")
	expectValue(t, match, 0, []string{"overlap", "month"}, "01")
	expectValue(t, match, 0, []string{"overlap", "month", "tens"}, "0")
	expectValue(t, match, 0, []string{"overlap", "month", "ones"}, "1")
	expectValue(t, match, 0, []string{"overlap", "day"}, "02")
	expectValue(t, match, 0, []string{"overlap", "day", "tens"}, "0")
	expectValue(t, match, 0, []string{"overlap", "day", "ones"}, "2")

	expectValue(t, match, 1, []string{}, "1234-56-78")
	expectValue(t, match, 1, []string{"overlap"}, "1234-56")
	expectValue(t, match, 1, []string{"overlap", "year"}, "1234")
	expectValue(t, match, 1, []string{"overlap", "year", "thousands"}, "1")
	expectValue(t, match, 1, []string{"overlap", "year", "hundreds"}, "2")
	expectValue(t, match, 1, []string{"overlap", "year", "tens"}, "3")
	expectValue(t, match, 1, []string{"overlap", "year", "ones"}, "4")
	expectValue(t, match, 1, []string{"overlap", "month"}, "56")
	expectValue(t, match, 1, []string{"overlap", "month", "tens"}, "5")
	expectValue(t, match, 1, []string{"overlap", "month", "ones"}, "6")
	expectValue(t, match, 1, []string{"overlap", "day"}, "78")
	expectValue(t, match, 1, []string{"overlap", "day", "tens"}, "7")
	expectValue(t, match, 1, []string{"overlap", "day", "ones"}, "8")
}

func TestSubExpUnnamed(t *testing.T) {
	// a regular expression to match a date
	re := regexp.MustCompile(`(((\d)(\d)(\d)(\d))-((\d)(\d)))-(((\d)(\d)))`)
	// a subject to match the regular expression
	subject := "this is a test subject to see if we can parse 2016-01-02 and 1234-56-78 using the Match() function."

	match, ok := subexpnames.Match(re, subject)
	if !ok {
		t.Errorf("expected a match")
	}

	if len(*match) != 2 {
		t.Errorf("expected 2 match, got %d", len(*match))
	}

	if len((*match)[0].Nested) != 2 {
		t.Errorf("expected 2 nested values, got %d", len((*match)[0].Nested))
	}
	if len((*match)[1].Nested) != 2 {
		t.Errorf("expected 2 nested values, got %d", len((*match)[1].Nested))
	}

	expectValues(t, match, 0, []string{}, []string{"2016-01-02"})
	expectValues(t, match, 0, []string{""}, []string{"2016-01", "02"})
	expectValues(t, match, 0, []string{"", ""}, []string{"2016", "01", "02"})
	expectValues(t, match, 0, []string{"", "", ""}, []string{"2", "0", "1", "6", "0", "1", "0", "2"})

	expectValues(t, match, 1, []string{}, []string{"1234-56-78"})
	expectValues(t, match, 1, []string{""}, []string{"1234-56", "78"})
	expectValues(t, match, 1, []string{"", ""}, []string{"1234", "56", "78"})
	expectValues(t, match, 1, []string{"", "", ""}, []string{"1", "2", "3", "4", "5", "6", "7", "8"})

}

func TestMultipleNestedGroups(t *testing.T) {
	regex := regexp.MustCompile(`(?P<outer>(?P<inner1>\d+)|(?P<inner2>[a-z]+))`)
	subject := "123abc"
	matches, ok := subexpnames.Match(regex, subject)
	if !ok {
		t.Errorf("Expected a match")
	}
	if matches.Len() != 2 {
		t.Errorf("Expected 2 matches, got %d", matches.Len())
	}

	if (*matches)[0].Nested[0].Key != "outer" || (*matches)[0].Nested[0].Value != "123" {
		t.Errorf("Expected outer match with value '123', got '%s'", (*matches)[0].Value)
	}

	if (*matches)[1].Nested[0].Key != "outer" || (*matches)[1].Nested[0].Value != "abc" {
		t.Errorf("Expected outer match with value 'abc', got '%s'", (*matches)[1].Value)
	}
}

func TestContiguousGroups(t *testing.T) {
	regex := regexp.MustCompile(`(?P<digits>\d+)(?P<alpha>\w+)`)
	subject := "123abc"
	matches, ok := subexpnames.Match(regex, subject)
	if !ok {
		t.Errorf("Expected a match")
	}
	if matches.Len() != 1 {
		t.Errorf("Expected 1 matches, got %d", matches.Len())
	}

	if (*matches)[0].Nested[0].Key != "digits" || (*matches)[0].Nested[0].Value != "123" {
		t.Errorf("Expected digits match with value '123', got '%s'", (*matches)[0].Nested[0].Value)
	}

	if (*matches)[0].Nested[1].Key != "alpha" || (*matches)[0].Nested[1].Value != "abc" {
		t.Errorf("Expected alpha match with value '123abc', got '%s'", (*matches)[0].Nested[1].Value)
	}
}

func TestOptionalGroups(t *testing.T) {
	regex := regexp.MustCompile(`(?P<year>\d{4})?(?P<month>\d{2})(?P<day>\d{2})`)
	subject := "20210304"
	matches, ok := subexpnames.Match(regex, subject)
	if !ok {
		t.Errorf("Expected a match")
	}

	if matches.Len() != 1 {
		t.Errorf("Expected 1 matches, got %d", matches.Len())
	}

	match := (*matches)[0]

	if match.Key != "" || match.Value != "20210304" {
		t.Errorf("Expected year match with value '20210304', got '%s'", match.Value)
	}

	if match.Nested[0].Key != "year" || match.Nested[0].Value != "2021" {
		t.Errorf("Expected month match with value '03', got '%s'", match.Nested[0].Value)
	}

	if s, _ := matches.Get(0, 0, "year"); s != "2021" {
		t.Errorf("Expected second year match with value '2021', got '%s'", s)
	}

	if match.Nested[1].Key != "month" || match.Nested[1].Value != "03" {
		t.Errorf("Expected month match with value '03', got '%s'", match.Nested[0].Value)
	}

	if s, _ := matches.Get(0, 0, "month"); s != "03" {
		t.Errorf("Expected second month match with value '03', got '%s'", s)
	}

	if match.Nested[2].Key != "day" || match.Nested[2].Value != "04" {
		t.Errorf("Expected day match with value '04', got '%s'", match.Nested[1].Value)
	}

	if s, _ := matches.Get(0, 0, "day"); s != "04" {
		t.Errorf("Expected second month match with value '04', got '%s'", s)
	}

}

func TestNamedGroupsUsedMultipleTimes(t *testing.T) {
	regex := regexp.MustCompile(`(?P<number>\d+).+?(?P<number>\d+)`)
	subject := "123abc456"
	matches, ok := subexpnames.Match(regex, subject)
	if !ok {
		t.Errorf("Expected a match")
	}

	if matches.Len() != 1 {
		t.Errorf("Expected 1 matches, got %d", matches.Len())
	}

	match := (*matches)[0]

	if match.Nested[0].Key != "number" || match.Nested[0].Value != "123" {
		t.Errorf("Expected first number match with value '123', got '%s'", match.Nested[0].Value)
	}

	if s, _ := matches.Get(0, 0, "number"); s != "123" {
		t.Errorf("Expected second number match with value '456', got '%s'", s)
	}

	if match.Nested[1].Key != "number" || match.Nested[1].Value != "456" {
		t.Errorf("Expected second number match with value '456', got '%s'", match.Nested[1].Value)
	}

	if s, _ := matches.Get(0, 1, "number"); s != "456" {
		t.Errorf("Expected second number match with value '456', got '%s'", s)
	}

}

func TestNonCapturingGroups(t *testing.T) {
	regex := regexp.MustCompile(`(?:\d{4})(?P<month>\d{2})(?P<day>\d{2})`)
	subject := "20210304"
	matches, ok := subexpnames.Match(regex, subject)
	if !ok {
		t.Errorf("Expected a match")
	}

	if matches.Len() != 1 {
		t.Errorf("Expected 1 matches, got %d", matches.Len())
	}

	match := (*matches)[0]

	if match.Nested[0].Key != "month" || match.Nested[0].Value != "03" {
		t.Errorf("Expected month match with value '03', got '%s'", match.Value)
	}

	if s, _ := matches.Get(0, 0, "month"); s != "03" {
		t.Errorf("Expected second number match with value '03', got '%s'", s)
	}

	if match.Nested[1].Key != "day" || match.Nested[1].Value != "04" {
		t.Errorf("Expected day match with value '04', got '%s'", match.Nested[0].Value)
	}

	if s, _ := matches.Get(0, 0, "day"); s != "04" {
		t.Errorf("Expected second number match with value '04', got '%s'", s)
	}

}
