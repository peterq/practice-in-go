package contest123

import (
	"log"
)

func addToArrayForm(A []int, K int) []int {

	var kk []int
	for K > 0 {
		n := K % 10
		K /= 10
		t := []int{n}
		kk = append(t, kk...)
	}
	if len(kk) > len(A) {
		A, kk = kk, A
	}
	A = append([]int{0}, A...)
	kk = append(make([]int, len(A)-len(kk)), kk...)
	more := 0
	for i := 0; i < len(kk); i++ {
		idxK := len(kk) - i - 1
		idxA := len(A) - i - 1
		add := A[idxA] + kk[idxK] + more
		A[idxA] = add % 10
		more = add / 10
	}
	if A[0] == 0 {
		A = A[1:]
	}

	return A
}

func Question1() {
	A := []int{1, 2, 6, 3, 0, 7, 1, 7, 1, 9, 7, 5, 6, 6, 4, 4, 0, 0, 6, 3}
	K := 516
	A = []int{9, 9, 9, 9, 9, 9}
	K = 1

	result := addToArrayForm(A, K)
	log.Println(A, K)
	log.Println(result)
}
