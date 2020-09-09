package main

// 最长回文子串
//给定一个字符串 s，找到 s 中最长的回文子串。你可以假设 s 的最大长度为 1000。
//输入: "babad"
//输出: "bab"
//注意: "aba" 也是一个有效答案。
//https://leetcode-cn.com/problems/longest-palindromic-substring/

func longestPalindrome(s string) string {
	if len(s) < 2 {
		return s
	}
	start := 0
	maxLen := 1
	for i := 0; i < len(s); i++ {
		getPalindrome(i-1, i+1, &maxLen, &start, s)
		getPalindrome(i, i+1, &maxLen, &start, s)
	}
	return s[start : start+maxLen]
}

func getPalindrome(left, right int, maxLen, start *int, s string) {
	for left >= 0 && right < len(s) && s[left] == s[right] {
		if right-left+1 > *maxLen {
			*start = left
			*maxLen = right - left + 1
		}
		left--
		right++
	}
}
