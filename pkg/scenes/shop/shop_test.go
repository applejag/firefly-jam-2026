package shop

import "testing"

func TestLongestName(t *testing.T) {
	longest := 0
	for _, kind := range AllItemKinds {
		longest = max(longest, len(kind.String()))
	}
	if longest != LongestItemKind {
		t.Errorf("LongestItemKind should be %d, but was %d", longest, LongestItemKind)
	}
}
