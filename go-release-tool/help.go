package main

import "fmt"

func displayHelp() {
	fmt.Println("Usage: jira-release-tool [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help       Show this help message and exit")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  JIRA_SERVER      The Jira server URL")
	fmt.Println("  JIRA_PROJECT     The Jira project key")
	fmt.Println("  JIRA_USER        The Jira username")
	fmt.Println("  JIRA_TOKEN       The Jira API token")
	fmt.Println("  JIRA_RELEASE_NAME The base name for the Jira release")
	fmt.Println("  GITHUB_REF_NAME  The GitHub reference name for the release version")
}
