package former

import (
	"strings"
)

/**

将字符串 "PAYPALISHIRING" 以Z字形排列成给定的行数：

P   A   H   N
A P L S I I G
Y   I   R
之后从左往右，逐行读取字符："PAHNAPLSIIGYIR"

实现一个将字符串进行指定行数变换的函数:

string convert(string s, int numRows);
示例 1:

输入: s = "PAYPALISHIRING", numRows = 3
输出: "PAHNAPLSIIGYIR"
示例 2:

输入: s = "PAYPALISHIRING", numRows = 4
输出: "PINALSIGYAHRPI"
解释:

P     I    N
A   L S  I G
Y A   H R
P     I
*/

func convert(s string, numRows int) string {
	if numRows == 1 {
		return s
	}
	rows := make([]string, numRows)
	idx := 0
	dirct := 1
	for _, c := range s {
		ch := string(c)
		rows[idx] += ch
		//log.Println(idx, dirct)
		if dirct == 1 && idx == numRows-1 {
			dirct = -1
		} else if dirct == -1 && idx == 0 {
			dirct = 1
		}
		idx += dirct
		/*for _, s := range rows {
			log.Println(s)
		}
		log.Println("-----")*/
	}
	return strings.Join(rows, "")
}
