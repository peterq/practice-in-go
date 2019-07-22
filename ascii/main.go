package main

import (
	"log"
	"math"
	"strings"
)

var sourceMap [62]string
var encodeMap [93]string
var encodeMapRevert = make(map[string]uint8)

func init() {
	str := "abcdefghijklmnopqrstuvwxyz"
	str += "ABCDEFGHIJKLNMOPQRSTUVWXYZ"
	str += "1234567890"
	for idx, ch := range str {
		sourceMap[idx] = string(ch)
	}
	str += "~`!@#$%^&*()_+-=<>,.?/:|\\;'{}[]"
	for idx, ch := range str {
		encodeMap[idx] = string(ch)
		encodeMapRevert[string(ch)] = uint8(idx)
	}
}

func encode(src string) string {
	result := ""
	for i := 0; i < len(src)-1; i += 3 {
		s := src[i : i+3]
		temp := decimalToAny(anyToDecimal(s, 62), 93)
		log.Println(s, temp, anyToDecimal(s, 62))
		temp = strings.Repeat(encodeMap[0], 2-len(temp)) + temp
		result += temp
	}
	return result
}

func decode(src string) string {
	result := ""
	for i := 0; i < len(src)-1; i += 2 {
		s := src[i : i+2]
		temp := decimalToAny(anyToDecimal(s, 93), 62)
		temp = strings.Repeat(encodeMap[0], 3-len(temp)) + temp
		result += temp
	}
	return result
}

// 10进制转任意进制
func decimalToAny(num, n int) string {
	new_num_str := ""
	var remainder int
	var remainder_string string
	for num != 0 {
		remainder = num % n
		remainder_string = encodeMap[remainder]
		new_num_str = remainder_string + new_num_str
		num = num / n
	}
	return new_num_str
}

// 任意进制转10进制
func anyToDecimal(num string, n int) int {
	new_num := float64(0.0)
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := float64(encodeMapRevert[value])
		new_num += tmp * math.Pow(float64(n), float64(nNum))
		nNum--
	}
	return int(new_num)
}

func main() {
	log.Println(encode("d41d8cd98f00b204e9800998ecf8427e"))
}
