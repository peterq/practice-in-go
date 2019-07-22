package former

/**
实现一个基本的计算器来计算一个简单的字符串表达式的值。

字符串表达式可以包含左括号 ( ，右括号 )，加号 + ，减号 -，非负整数和空格  。

示例 1:

输入: "1 + 1"
输出: 2
示例 2:

输入: " 2-1 + 2 "
输出: 3
示例 3:

输入: "(1+(4+5+2)-3)+(6+8)"
输出: 23
说明：

你可以假设所给定的表达式都是有效的。
请不要使用内置的库函数 eval。
*/

func calculate(s string) int {
	s = "(" + s + ")"
	// stack 符号加数字存储
	stk := []int{1, 0}
	p := 1
	number := 0
	for _, ch := range s {
		//log.Println(ch)
		if ch == '(' {
			stk = append(stk, p, 0)
			p = 1
		} else if ch == ')' {
			stk[len(stk)-1] += number * p
			//log.Println(stk)
			val := stk[len(stk)-1] * stk[len(stk)-2]
			stk = stk[:len(stk)-2]
			stk[len(stk)-1] += val
			number = 0
		} else if ch >= '0' && ch <= '9' {
			number *= 10
			number += int(ch - '0')
		} else if ch == '+' || ch == '-' {
			stk[len(stk)-1] += number * p
			number = 0
			p = 1
			if ch == '-' {
				p = -1
			}
		}
		//log.Println(stk, string(ch))
	}
	return stk[0] * stk[1]
}
