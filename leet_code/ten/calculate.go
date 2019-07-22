package ten

/**
实现一个基本的计算器来计算一个简单的字符串表达式的值。

字符串表达式仅包含非负整数，+， - ，*，/ 四种运算符和空格  。 整数除法仅保留整数部分。

示例 1:

输入: "3+2*2"
输出: 7
示例 2:

输入: " 3/2 "
输出: 1
示例 3:

输入: " 3+5 / 2 "
输出: 5
说明：

你可以假设所给定的表达式都是有效的。
请不要使用内置的库函数 eval。
*/

func calculate(s string) int {
	s = s + "+0"
	number := 0
	result := 0
	group := 1
	outer := 1
	inner := '*'
	for _, ch := range s {
		if ch == '+' || ch == '-' {
			if inner == '*' {
				group *= number
			} else {
				group /= number
			}
			result += outer * group
			inner = '*'
			group = 1
			number = 0
			if ch == '+' {
				outer = 1
			} else {
				outer = -1
			}
		}
		if ch == '*' || ch == '/' {
			if inner == '*' {
				group *= number
			} else {
				group /= number
			}
			inner = ch
			number = 0
		}
		if ch <= '9' && ch >= '0' {
			number *= 10
			number += int(ch - '0')
		}
	}
	return result
}
