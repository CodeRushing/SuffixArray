package suffix_array

import (
	"fmt"
	"sort"
	"unicode/utf8"
)

type sortRuneStruct struct {
	Index *[]int
	Data  *[]rune
}

type SuffixWord struct {
	Offset int
	Len    int
}

func (s sortRuneStruct) Len() int {
	return len(*(s.Data))
}

func (s sortRuneStruct) Swap(i, j int) {
	(*s.Index)[i], (*s.Index)[j] = (*s.Index)[j], (*s.Index)[i]
}

func (s sortRuneStruct) Less(i, j int) bool {
	return (*s.Data)[(*s.Index)[i]] < (*s.Data)[(*s.Index)[j]]
}

func sortByFirstRune(data []rune) []int {
	data_len := len(data)
	sa := make([]int, data_len)
	for i := 0; i < data_len; i++ {
		sa[i] = i
	}

	sort.Stable(sortRuneStruct{Index: &sa, Data: &data})
	return sa
}

func BuildSuffixByDoublingAlgorithm(data []rune, max_word_len int) []int {
	data_len := len(data)
	sa := sortByFirstRune(data)
	rank_x := make([]int, data_len)
	rank_y := make([]int, data_len)
	count := make([]int, data_len)
	wv := make([]int, data_len)

	rank_x[sa[0]] = 0
	rank_x_len := 1
	for i := 1; i < data_len; i++ {
		if data[sa[i-1]] == data[sa[i]] {
			rank_x[sa[i]] = rank_x_len - 1
		} else {
			rank_x[sa[i]] = rank_x_len
			rank_x_len++
		}
	}

	base_sort_range := rank_x_len
	for i := 1; i < max_word_len; i *= 2 {
		rank_y_len := 0
		for j := data_len - i; j < data_len; j++ {
			rank_y[rank_y_len] = j
			rank_y_len++
		}
		for j := 0; j < data_len; j++ {
			if sa[j] >= i {
				rank_y[rank_y_len] = sa[j] - i
				rank_y_len++
			}
		}
		for j := 0; j < data_len; j++ {
			wv[j] = rank_x[rank_y[j]]
		}

		for j := 0; j < base_sort_range; j++ {
			count[j] = 0
		}
		for j := 0; j < data_len; j++ {
			count[wv[j]]++
		}
		for j := 1; j < base_sort_range; j++ {
			count[j] += count[j-1]
		}
		for j := data_len - 1; j >= 0; j-- {
			count[wv[j]]--
			sa[count[wv[j]]] = rank_y[j]
		}
		rank_x, rank_y = rank_y, rank_x
		rank_x[sa[0]] = 0
		rank_x_len = 1
		is_sorted := true
		for j := 1; j < data_len; j++ {
			if rank_y[sa[j-1]] == rank_y[sa[j]] && (sa[j-1]+i) < data_len &&
				(sa[j]+i) < data_len && rank_y[sa[j-1]+i] == rank_y[sa[j]+i] {
				rank_x[sa[j]] = rank_x_len - 1
				is_sorted = false
			} else {
				rank_x[sa[j]] = rank_x_len
				rank_x_len++
			}
		}
		if is_sorted {
			break
		}
		base_sort_range = rank_x_len
	}
	return sa
}

func printRune(rune_data rune) string {
	buf := make([]byte, 3)
	utf8.EncodeRune(buf, rune_data)
	return string(buf)
}

func outputSa(sa []int, data_rune []rune, max_word_len int) {
	buf := make([]byte, 3)
	for i := range sa {
		var suffix string = ""
		for j := sa[i]; j < len(data_rune) && j < sa[i]+max_word_len; j++ {
			utf8.EncodeRune(buf, data_rune[j])
			suffix += string(buf)
		}

		fmt.Println(i, suffix, sa[i])
	}

}

func (word SuffixWord) String(data []rune) string {
	buf := make([]byte, 3)
	var word_string string = ""
	for i := 0; i < word.Len; i++ {
		n := utf8.EncodeRune(buf, data[i+word.Offset])
		word_string += string(buf[:n])
	}
	return word_string
}

func GetHeight(data []rune, sa []int) []int {
	data_len := len(data)
	rank := make([]int, data_len, data_len)
	height := make([]int, data_len, data_len)
	for i := 0; i < data_len; i++ {
		rank[sa[i]] = i
	}

	prefix_len := 0
	for i := 0; i < data_len; i++ {
		if rank[i] == 0 {
			height[rank[i]] = 0
			continue
		}
		pre_index := sa[rank[i]-1]
		for (i+prefix_len) < data_len && (pre_index+prefix_len) < data_len && data[i+prefix_len] == data[pre_index+prefix_len] {
			prefix_len++
		}
		fmt.Println(prefix_len)
		height[rank[i]] = prefix_len
		if prefix_len--; prefix_len < 0 {
			prefix_len = 0
		}
	}
	return height
}
