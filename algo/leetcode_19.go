package main

//删除链表的倒数第N个节点
//给定一个链表，删除链表的倒数第 n 个节点，并且返回链表的头结点。
//
//示例：
//
//给定一个链表: 1->2->3->4->5, 和 n = 2.
//
//当删除了倒数第二个节点后，链表变为 1->2->3->5.
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/remove-nth-node-from-end-of-list

//func removeNthFromEnd(head *ListNode, n int) *ListNode {
//	dummy := &ListNode{Next: head}
//	tmp := dummy
//	pre := dummy
//	for i := 0; i < n; i++ {
//		tmp = tmp.Next
//	}
//	for tmp.Next != nil {
//		pre = pre.Next
//		tmp = tmp.Next
//	}
//	pre.Next = pre.Next.Next
//	return head
//}
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	if n == 0 {
		return head
	}
	dummy := &ListNode{Next: head}
	tmp := dummy
	pre := dummy
	for i := 0; i <= n; i++ {
		tmp = tmp.Next
	}
	for tmp != nil {
		pre = pre.Next
		tmp = tmp.Next
	}
	pre.Next = pre.Next.Next
	return dummy.Next
}