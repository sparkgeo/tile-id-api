package handler

import "slices"

func SortIdentifiers(allIdentifiers []string, firstValue string) []string {
	ownAllIdentifiers := allIdentifiers[:]
	slices.Sort(ownAllIdentifiers)
	split := -1
	for i, identifier := range ownAllIdentifiers {
		if identifier == firstValue {
			split = i
		}
	}
	if split == -1 {
		return ownAllIdentifiers
	}
	if split == 0 {
		return ownAllIdentifiers
	}
	return append(append(append([]string{}, firstValue), ownAllIdentifiers[0:split]...), ownAllIdentifiers[split+1:]...)
}
