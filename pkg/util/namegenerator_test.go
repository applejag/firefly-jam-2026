package util

import "testing"

func TestLongestName(t *testing.T) {
	longest := 0
	for _, name := range names {
		longest = max(longest, len(name))
	}
	if longest != longestName {
		t.Errorf("longestName should be %d, but was %d", longest, longestName)
	}
}

func TestLongestTitle(t *testing.T) {
	longest := 0
	for _, title := range titles {
		longest = max(longest, len(title))
	}
	if longest != longestTitle {
		t.Errorf("longestTitle should be %d, but was %d", longest, longestTitle)
	}
}
