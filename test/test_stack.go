package test

import (
	"log"
	"strconv"
)

func stackTest() {
	log.Println(parseIns("r2(lr3(ui))y1(z)"))
}

func parseIns(ins string) string {
	stack := make([]int, 0)
	res := ""
	tempNumber := ""
	for _, v := range ins {
		c := string(v)
		if c >= "0" && c <= "9" {
			tempNumber += c
			continue
		}

		if c == "(" {
			repeatNumber, _ := strconv.Atoi(tempNumber)
			stack = append(stack, repeatNumber, len(res))
			tempNumber = ""
			continue
		}

		if c == ")" {
			repeatNumber := stack[len(stack)-2]
			repeatStr := res[stack[len(stack)-1]:]
			stack = stack[:len(stack)-2]
			for i := 0; i < repeatNumber-1; i++ {
				res += repeatStr
			}
			continue
		}

		res += c

	}
	return res
}
