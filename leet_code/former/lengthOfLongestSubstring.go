package former

/**
给定一个字符串，找出不含有重复字符的最长子串的长度。

示例 1:

输入: "abcabcbb"
输出: 3
解释: 无重复字符的最长子串是 "abc"，其长度为 3。
示例 2:

输入: "bbbbb"
输出: 1
解释: 无重复字符的最长子串是 "b"，其长度为 1。
示例 3:

输入: "pwwkew"
输出: 3
解释: 无重复字符的最长子串是 "wke"，其长度为 3。
     请注意，答案必须是一个子串，"pwke" 是一个子序列 而不是子串。
*/

func lengthOfLongestSubstring(s string) int {
	// 出现过的字符map
	maxLength := 0
	mp := make(map[string]int)
	for pos, v := range s {
		c := string(v)
		if posPre, ok := mp[c]; !ok { // 没有重复的
			mp[c] = pos
		} else { // 重复了
			if len(mp) > maxLength {
				maxLength = len(mp)
			}
			for k, v := range mp {
				if v <= posPre {
					delete(mp, k)
				}
			}
			mp[c] = pos
		}
	}
	if maxLength < len(mp) {
		maxLength = len(mp)
	}
	return maxLength
}
