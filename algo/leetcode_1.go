package main

//两数之和
//给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那 两个 整数，并返回他们的数组下标。
//
//你可以假设每种输入只会对应一个答案。但是，数组中同一个元素不能使用两遍
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/two-sum

func twoSum(nums []int, target int) []int {
	var ret []int
	tmp := make(map[int]int)

	for i, v := range nums {
		if j, ok := tmp[target-v]; ok {
			return []int{j, i}
		} else {
			tmp[v] = i
		}
	}
	return ret
}
