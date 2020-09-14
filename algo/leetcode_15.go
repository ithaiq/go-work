package main

import (
	"sort"
)

// 三数之和
//给你一个包含 n 个整数的数组 nums，判断 nums 中是否存在三个元素 a，b，c ，使得 a + b + c = 0 ？请你找出所有满足条件且不重复的三元组。
//
//注意：答案中不可以包含重复的三元组。
//
// 
//
//示例：
//
//给定数组 nums = [-1, 0, 1, 2, -1, -4]，
//
//满足要求的三元组集合为：
//[
//  [-1, 0, 1],
//  [-1, -1, 2]
//]
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/3sum

func threeSum(nums []int) [][]int {
	if len(nums) < 3 {
		return nil
	}
	sort.Ints(nums)
	var result [][]int

	for i := 0; i <= len(nums)-3; i++ {
		if i == 0 || nums[i] != nums[i-1] {
			start, end := i+1, len(nums)-1
			for start < end {
				if nums[start]+nums[end]+nums[i] == 0 {
					result = append(result, []int{nums[i], nums[start], nums[end]})
					start++
					end--
					for start < end && nums[start] == nums[start-1] {
						start++
					}
					for start < end && nums[end] == nums[end+1] {
						end--
					}
				} else if nums[start]+nums[end]+nums[i] < 0 {
					start++
				} else {
					end--
				}
			}
		}
	}
	return result
}
