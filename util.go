package menuscreen

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func cellCnt(rs []rune) int {

	cnt := 0

	for _, r := range rs {
		// some unicode char takes 2 cells, but its strlen is 3
		cnt += min(2, len(string(r)))
	}

	return cnt

}

func cloneRuneSlice(src []rune) []rune {
	dst := make([]rune, len(src))
	copy(dst, src)
	return dst
}
