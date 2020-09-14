package main

//给定一个字符串，请你找出其中不含有重复字符的 最长子串 的长度。
//输入: "pwwkew"
//输出: 3
//解释: 因为无重复字符的最长子串是 "wke"，所以其长度为 3。
//     请注意，你的答案必须是 子串 的长度，"pwke" 是一个子序列，不是子串。
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/longest-substring-without-repeating-characters

func lengthOfLongestSubstring(s string) int {
	if len(s) < 2 {
		return len(s)
	}
	j := 0
	maxLen := 0
	tmp := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		if _, ok := tmp[s[i]]; ok {
			for j < i {
				if s[j] == s[i] {
					tmp[s[i]] = i
					j++
					break
				}
				delete(tmp, s[j])
				j++
			}
		} else {
			tmp[s[i]] = i
			if len(tmp) > maxLen {
				maxLen = len(tmp)
			}
		}
	}
	return maxLen
}