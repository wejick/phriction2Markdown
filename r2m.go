package main

import (
	"bufio"
	"fmt"
	"strings"
)

func mdifize(text string) (revisedText string) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		line = convertOrderedList(line)
		line = convertHeading(line)
		revisedText += fmt.Sprintln(line)
	}
	return
}

func convertHeading(line string) (revisedLine string) {
	revisedLine = strings.TrimSpace(line)

	// H3
	isH3 := strings.HasPrefix(revisedLine, "===") && strings.HasSuffix(revisedLine, "===")

	if isH3 {
		revisedLine = strings.Replace(revisedLine, "===", "### ", 1)
		revisedLine = strings.TrimRight(revisedLine, "===")

		return
	}

	// H2
	isH2 := strings.HasPrefix(revisedLine, "==") && strings.HasSuffix(revisedLine, "==")

	if isH2 {
		revisedLine = strings.Replace(revisedLine, "==", "## ", 1)
		revisedLine = strings.TrimRight(revisedLine, "==")

		return
	}

	// H1
	isH1 := strings.HasPrefix(revisedLine, "=") && strings.HasSuffix(revisedLine, "=")

	if isH1 {
		revisedLine = strings.Replace(revisedLine, "=", "# ", 1)
		revisedLine = strings.TrimRight(revisedLine, "=")

		return
	}

	return
}

func convertOrderedList(line string) (revisedLine string) {
	revisedLine = line
	noSpace := strings.TrimSpace(line)
	isNumList := strings.HasPrefix(noSpace, "#") && !strings.HasSuffix(noSpace, "#")

	if isNumList {
		revisedLine = strings.Replace(line, "#", "1.", 1)
	}

	return
}

func convertTable() {}
