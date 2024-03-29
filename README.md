# SubExpNames

[![Go Report Card](https://goreportcard.com/badge/github.com/thetechpanda/subexpnames)](https://goreportcard.com/report/github.com/thetechpanda/subexpnames)
[![Go Reference](https://pkg.go.dev/badge/github.com/thetechpanda/subexpnames.svg)](https://pkg.go.dev/github.com/thetechpanda/subexpnames)
[![Release](https://img.shields.io/github/release/thetechpanda/subexpnames.svg?style=flat-square)](https://github.com/thetechpanda/subexpnames/releases)

## Documentation

You can find the generated go doc [here](godoc.txt).

## Motivation

Working with regular expressions in Go can be challenging, especially when it comes to handling named capture groups in a hierarchical manner. The standard `regexp` package provides basic support for matching patterns and extracting groups, but it falls short when dealing with complex expressions where named groups are nested or repeated.

The `subexpnames` package is designed to bridge this gap. It offers an intuitive way to work with named groups in regular expressions, allowing developers to retrieve matched values in the same hierarchical structure as defined in the expression. This makes it easier to parse and manipulate data captured by complex regex patterns, simplifying tasks that would otherwise require cumbersome manual processing.

By providing a more user-friendly interface for dealing with named groups, `subexpnames` aims to enhance the usability of regular expressions in Go, enabling developers to focus more on their core logic and less on the intricacies of pattern matching.

## Usage

For more examples look at [example_test.go](example_test.go)

```go
package main

import (
	"fmt"
	"github.com/thetechpanda/subexpnames"
)

func main() {
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
}
```

## Code coverage

```
$ go test -cover ./...
ok      github.com/thetechpanda/subexpnames     0.268s  coverage: 100.0% of statements
```

## Installation

```bash
go get github.com/thetechpanda/subexpnames
```

## Contributing

Contributions are welcome and very much appreciated! 

Feel free to open an issue or submit a pull request.

## License

`SubExpNames` is released under the MIT License. See the [LICENSE](LICENSE) file for details.
