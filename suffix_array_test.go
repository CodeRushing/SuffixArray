package suffix_array

import (
	"bytes"
	"testing"
)

func TestSuffixArray(t *testing.T) {
	data := bytes.Runes([]byte("aabaaaab"))
	sa := BuildSuffixByDoublingAlgorithm(data, len(data))
	result := []string{"aaaab", "aaab", "aab", "aabaaaab", "ab", "abaaaab", "b", "baaaab"}
	for i, index := range sa {
		suffix_word := SuffixWord{Offset: index, Len: len(data) - index}
		if suffix_word.String(data) != result[i] {
			t.Errorf("suffix build error [%d:%s %s]", i, suffix_word.String(data), result[i])
		}
	}
	height := GetHeight(data, sa)
	result_height := []int{0, 3, 2, 3, 1, 2, 0, 1}
	for i := 0; i < len(height); i++ {
		if height[i] != result_height[i] {
			t.Errorf("height[%d] is wrong![right:%d][wrong:%d]", i, result_height[i], height[i])
		}
	}
}
