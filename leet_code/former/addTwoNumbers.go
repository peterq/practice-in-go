package former

/**
给定两个非空链表来表示两个非负整数。位数按照逆序方式存储，它们的每个节点只存储单个数字。将两数相加返回一个新的链表。

你可以假设除了数字 0 之外，这两个数字都不会以零开头。

示例：

输入：(2 -> 4 -> 3) + (5 -> 6 -> 4)
输出：7 -> 0 -> 8
原因：342 + 465 = 807
*/

/**
* Definition for singly-linked list.
* type ListNode struct {
*     Val int
*     Next *ListNode
* }
 */
type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	full := 0
	result := new(ListNode)
	l1Node := l1
	l2Node := l2
	resultNode := result
	for {
		sum := l1Node.Val + l2Node.Val + full
		resultNode.Next = &ListNode{
			Val: sum % 10,
		}
		resultNode = resultNode.Next
		if sum > 9 {
			full = 1
		} else {
			full = 0
		}
		l1Node = l1Node.Next
		l2Node = l2Node.Next
		if l1Node == nil && l2Node == nil {
			if full == 1 {
				resultNode.Next = &ListNode{
					full, nil,
				}
			}
			break
		}
		if l1Node == nil {
			l1Node = &ListNode{}
		}
		if l2Node == nil {
			l2Node = &ListNode{}
		}
	}
	return result.Next
}
