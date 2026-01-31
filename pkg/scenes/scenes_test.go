package scenes

import "testing"

func TestLongestName(t *testing.T) {
	longest := 0
	for _, scene := range AllScenes {
		longest = max(longest, len(scene.String()))
	}
	if longest != LongestSceneName {
		t.Errorf("LongestSceneName should be %d, but was %d", longest, LongestSceneName)
	}
}
