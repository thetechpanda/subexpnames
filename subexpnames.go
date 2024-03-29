// Package subexpnames provides functions for matching regular expressions against strings and extracting match values in a hierarchical manner.
//
// subexpnames.Match first creates a tree-like structure of matches by matching a regular expression against a subject string.
//
// The way the package splits the regular expression into a tree-like structure is by using the regular expression's FindAllStringSubmatchIndex and FindAllStringSubmatch methods.
// It then finds all matches containing the start and end indexes of the match and the corresponding key and value.
//
// Calls to Matches' functions recursively descend into the nested matchValues to find the appropriate match, for this reason using this package on large regular expressions can be slow.
package subexpnames

import (
	"regexp"
	"slices"
)

// MatchValue represents a single match found in the subject string that corresponds to the regular expression.
// It contains the following information about the match:
//   - Key: A string that identifies the match, often corresponding to a named capture group in the regular expression.
//   - Value: The substring from the subject string that was matched.
//   - start: The starting index of the match in the subject string.
//   - end: The ending index of the match in the subject string.
//   - Nested: A slice of pointers to MatchValue structs representing any nested matches.
//     Nested matches occur when the regular expression contains capture groups within other capture groups.
//     This allows for representing the hierarchical structure of matches in a tree-like form.
type MatchValue struct {
	Key    string
	Value  string
	Nested []*MatchValue
	// start and end represent the indexes of the match in the subject string, they are not exported.
	start, end int
}

// Matches represents a collection of MatchValue pointers.
// It is used to store multiple matches found in a subject string that match a regular expression.
type Matches []*MatchValue

// find searches for the first MatchValue within the specified start and end indexes in the given bound.
// It recursively descends into the nested matchValues to find the appropriate match.
// Returns the found MatchValue and a boolean indicating whether the match is contained within the bounds.
// If no match is found within the bounds, it returns nil and false.
func find(bound *MatchValue, start, end int) (mv *MatchValue, contained bool) {
	for _, inner := range bound.Nested {
		if start >= inner.start && end <= inner.end {
			if v, contained := find(inner, start, end); contained {
				return v, contained
			}
			return inner, true
		}
	}
	return nil, false
}

// tree takes a regular expression and a subject string, and returns a hierarchical structure of matches.
// The function constructs a tree-like structure where each node represents a match found in the subject string.
// The tree is built using the regular expression's submatches and their corresponding start and end indexes.
// This function is useful for organizing matches in a way that reflects their nested nature in the regular expression.
func tree(re *regexp.Regexp, subject string) *Matches {
	matches := make([]*MatchValue, 0)
	indexes := re.FindAllStringSubmatchIndex(subject, -1)
	names := re.SubexpNames()
	match := re.FindAllStringSubmatch(subject, -1)
	for i := 0; i < len(indexes); i++ {
		matches = append(matches, &MatchValue{
			Key:    names[i],
			Value:  match[i][0],
			start:  indexes[i][0],
			end:    indexes[i][1],
			Nested: make([]*MatchValue, 0),
		})

		bound := matches[i]
		for j := 2; j < len(indexes[i]); j += 2 {
			matchStart, matchEnd, key, value := indexes[i][j], indexes[i][j+1], names[j/2], match[i][j/2]
			inner, contained := find(bound, matchStart, matchEnd)
			if contained {
				inner.Nested = append(inner.Nested, &MatchValue{
					Key:    key,
					Value:  value,
					start:  matchStart,
					end:    matchEnd,
					Nested: make([]*MatchValue, 0),
				})
				continue
			}
			bound.Nested = append(bound.Nested, &MatchValue{
				Key:    key,
				Value:  value,
				start:  matchStart,
				end:    matchEnd,
				Nested: make([]*MatchValue, 0),
			})
		}
	}
	return (*Matches)(&matches)
}

// descend is a helper function that recursively descends into the nested matchValues to retrieve the values based on the provided keys.
func descend(bound *MatchValue, keys ...string) (values []string) {
	if len(keys) == 0 {
		return []string{bound.Value}
	}
	key := keys[0]
	for _, mv := range bound.Nested {
		if mv.Key == key {
			if len(keys) == 1 {
				values = append(values, mv.Value)
			} else {
				values = append(values, descend(mv, keys[1:]...)...)
			}
		}
	}
	return values
}

// Match checks if the subject string matches the provided regular expression.
// If a match is found, it returns a regMatch object containing the tree-like structure of matchValues.
// Otherwise, it returns nil and false.
func Match(regexp *regexp.Regexp, subject string) (*Matches, bool) {
	if !regexp.MatchString(subject) {
		return nil, false
	}
	return tree(regexp, subject), true
}

// GetAll retrieves all the values that match the provided keys from the specified match.
// If the match or keys are not found, it returns nil and false.
// Otherwise, it returns a slice of strings containing the matching values and true.
func (rm *Matches) GetAll(group int, keys ...string) ([]string, bool) {
	if group < 0 || group >= len(*rm) {
		return nil, false
	}
	var values []string = descend((*rm)[group], keys...)
	return values, len(values) > 0
}

// Get retrieves the value at the specified index from the specified match.
// If the match, value, or keys are not found, it returns an empty string and false.
// The function allows for accessing specific values within a group of matches based on their keys.
func (rm *Matches) Get(group int, value int, keys ...string) (string, bool) {
	values, ok := rm.GetAll(group, keys...)
	if !ok {
		return "", false
	}

	if value < 0 || value >= len(values) {
		return "", false
	}

	return values[value], true
}

// GetFirstValueOfGroup retrieves the first value of the group that matches the provided keys.
// If the keys sequence is not found, it returns an empty string and false.
// This function is a convenience method for quickly accessing the first value in a group of matches.
func (rm *Matches) GetFirstValueOfGroup(group int, keys ...string) (string, bool) {
	return rm.Get(group, 0, keys...)
}

// GetGroup retrieves the match at the specified index from the Matches object.
// If the index is out of bounds, it returns nil and false.
// This function provides access to individual match groups within the collection of matches.
func (rm *Matches) GetGroup(group int) (*MatchValue, bool) {
	if group < 0 || group >= len(*rm) {
		return nil, false
	}
	return (*rm)[group], true
}

// Len returns the number of groups in the Matches object.
func (rm *Matches) Len() int {
	return len(*rm)
}

// descendKeys is a helper function that recursively descends into the nested matchValues to retrieve the keys.
// It accumulates the keys in the 'values' slice, ensuring that each key pair is added only once.
// The 'parents' parameter is used to keep track of the hierarchy of keys during the recursion.
func descendKeys(bound *MatchValue, values *[][]string, parents ...string) {
	for _, mv := range bound.Nested {
		pk := append([]string{}, parents...)
		pk = append(pk, mv.Key)
		alreadyExists := false
		for _, v := range *values {
			if slices.Compare(v, pk) == 0 {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			*values = append(*values, pk)
		}
		descendKeys(mv, values, pk...)
	}
}

// Keys retrieves all the keys from the specified group.
// It returns a slice of slices of strings containing the keys and the keys of their nested matches.
// If a key pair is repeated, it will only be added once.
func (rm *Matches) Keys(group int) [][]string {
	if group < 0 || group >= len(*rm) {
		return nil
	}
	var keys [][]string
	descendKeys((*rm)[group], &keys)
	return keys
}
