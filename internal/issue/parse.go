package issue

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var rangeRegex = regexp.MustCompile(`#?(\d+)\s*-\s*#?(\d+)`)

func Parse(args []string) ([]int, error) {
	raw := strings.Join(args, " ")

	raw = strings.ReplaceAll(raw, ",", " ")
	raw = strings.ReplaceAll(raw, " and ", " ")

	var issues []int

	ranges := rangeRegex.FindAllStringSubmatch(raw, -1)
	for _, match := range ranges {
		start, _ := strconv.Atoi(match[1])
		end, _ := strconv.Atoi(match[2])
		if end < start {
			start, end = end, start
		}
		for i := start; i <= end; i++ {
			issues = append(issues, i)
		}
		raw = strings.Replace(raw, match[0], "", 1)
	}

	numRegex := regexp.MustCompile(`#?(\d+)`)
	singles := numRegex.FindAllStringSubmatch(raw, -1)
	for _, match := range singles {
		num, _ := strconv.Atoi(match[1])
		issues = append(issues, num)
	}

	if len(issues) == 0 {
		return nil, fmt.Errorf("no valid issue numbers found in: %s", strings.Join(args, " "))
	}

	return dedupe(issues), nil
}

func dedupe(nums []int) []int {
	seen := make(map[int]bool)
	var result []int
	for _, n := range nums {
		if !seen[n] {
			seen[n] = true
			result = append(result, n)
		}
	}
	return result
}
