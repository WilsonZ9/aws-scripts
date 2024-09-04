package main

import (
	"os"
	"regexp"
	"strings"
)

const (
	changesSection = "## What's Changed\n"
)

func getSection(mdContent, sectionTitle string) string {
	parts := strings.Split(mdContent, sectionTitle)
	if len(parts) < 2 {
		return ""
	}
	sectionContent := strings.Split(parts[1], "\n\n")[0]
	return sectionContent
}

func parseChangelist(content string) []map[string]string {
	items := []map[string]string{}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 2 {
			continue
		}
		line = line[2:] // Remove leading "- "
		parts := strings.Split(line, " by @")
		if len(parts) < 2 {
			continue
		}
		prTitle := parts[0]
		authorAndLink := strings.Split(parts[1], " in ")
		if len(authorAndLink) < 2 {
			continue
		}
		author := authorAndLink[0]
		prLink := authorAndLink[1]
		items = append(items, map[string]string{
			"title":  prTitle,
			"author": author,
			"link":   prLink,
		})
	}
	return items
}

func extractChanges(filename string) ([]map[string]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	mdContent := string(content)
	if !strings.Contains(mdContent, changesSection) {
		return nil, nil
	}

	sectionContent := getSection(mdContent, changesSection)
	return parseChangelist(sectionContent), nil
}

func extractIssueID(change, project string) string {
	issuePattern := regexp.MustCompile(`\b` + project + `-\d+\b`)
	matches := issuePattern.FindAllString(change, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}
