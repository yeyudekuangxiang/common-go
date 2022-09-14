package open

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestMinCut(t *testing.T) {
	s := "21212"
	n := len(s)
	f := make([][]int8, n)
	for i := range f {
		f[i] = make([]int8, n)
	}
	var isPalindrome func(i, j int) int8
	isPalindrome = func(i, j int) int8 {
		if i >= j {
			return 1
		}
		if f[i][j] != 0 {
			return f[i][j]
		}
		f[i][j] = -1
		if s[i] == s[j] {
			f[i][j] = isPalindrome(i+1, j-1)
		}
		return f[i][j]
	}

	var splits []string
	var dfs func(int)
	var ret []int
	var res string
	ans := make(map[string]struct{}, n*2)

	dfs = func(i int) {
		if i == n {
			return
		}
		for j := i; j < n; j++ {
			if isPalindrome(i, j) > 0 {
				ans[s[i:j+1]] = struct{}{}
				splits = append(splits, s[i:j+1])
				dfs(j + 1)
				splits = splits[:len(splits)-1]
			}
		}
	}
	dfs(0)
	l := len(ans)
	for i, _ := range ans {
		atoi, _ := strconv.Atoi(i)
		ret = append(ret, atoi)
	}
	sort.Ints(ret)
	for i := l - 1; i > 0; i-- {
		res += strconv.Itoa(ret[i])
	}
	fmt.Println(res)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func isPalindrome(s string) bool {
	var sgood string
	for i := 0; i < len(s); i++ {
		if isalnum(s[i]) {
			sgood += string(s[i])
		}
	}

	n := len(sgood)
	sgood = strings.ToLower(sgood)
	for i := 0; i < n/2; i++ {
		if sgood[i] != sgood[n-1-i] {
			return false
		}
	}
	return true
}

func isalnum(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')
}
