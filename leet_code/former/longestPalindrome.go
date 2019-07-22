package former

import (
	"strings"
)

/**
给定一个字符串 s，找到 s 中最长的回文子串。你可以假设 s 的最大长度为1000。

示例 1：

输入: "babad"
输出: "bab"
注意: "aba"也是一个有效答案。
示例 2：

输入: "cbbd"
输出: "bb"
*/

func longestPalindrome(s string) string {
	if len(s) <= 1 {
		return s
	}
	ori := strings.Split(s, "") // 字符数组
	lg := len(ori)
	dest := make([]string, 0)
	lastCenter := 2*lg - 1
	destLg := 0
	for ctr := 1; ctr <= lastCenter; ctr++ {
		odd := ctr%2 == 1
		tmp := make([]string, 0)
		// 假设中心点是奇数
		startRight := (ctr + 1) / 2
		startLeft := ctr / 2
		if !odd { // 如果是偶数
			startRight = ctr/2 + 1
			startLeft = ctr/2 - 1
		}
		for idxL, idxR := startLeft, startRight; idxL >= 0 && idxR < lg; idxL, idxR = idxL-1, idxR+1 {
			if ori[idxL] == ori[idxR] {
				tmp = append(tmp, ori[idxR])
			} else {
				break
			}
		}
		if odd {
			if len(tmp)*2 > destLg {
				dest = tmp
				destLg = len(tmp) * 2
			}
		} else { // 偶数需要加上1
			if len(tmp) > 0 && len(tmp)*2+1 > destLg {
				dest = append([]string{ori[ctr/2]}, tmp...)
				destLg = len(tmp)*2 + 1
			}
		}
	}
	add := ""
	if len(dest) > 0 {
		for i := 0; i < len(dest); i++ {
			add = dest[i] + add
		}
		if destLg%2 == 1 {
			add = add[:len(add)-1]
		}
	}
	if len(dest) == 1 && destLg%2 == 1 {
		return ""
	}
	ret := add + strings.Join(dest, "")
	if len(ret) == 0 {
		return ori[0]
	}
	return ret
}
