package main

import (
	"log"
	"net/http"
	"os"
)

// Jira/Github Automation Tool

const (
	envJiraServer      = "JIRA_SERVER"
	envJiraProject     = "JIRA_PROJECT"
	envJiraUser        = "JIRA_USER"
	envJiraToken       = "JIRA_TOKEN"
	envJiraReleaseName = "JIRA_RELEASE_NAME"
	envGithubRefName   = "GITHUB_REF_NAME"
	envGithubRepo      = "GITHUB_REPOSITORY"
	envGithubToken     = "GH_TOKEN"
	envPreviousVersion = "PREVIOUS_VERSION"
	envNewVersion      = "NEW_VERSION"
	releaseNotesFile   = "release_notes.md"
)

var (
	checkAndCreateGithubReleaseFunc = checkAndCreateGithubRelease
	createJiraReleaseFunc           = createJiraRelease
	extractChangesFunc              = extractChanges
	addReleaseToIssueFunc           = addReleaseToIssue
	extractIssueIDFunc              = extractIssueID
)

func checkEnvVars(vars ...string) {
	for _, v := range vars {
		if os.Getenv(v) == "" {
			log.Fatalf("Missing required environment variable: %s", v)
		}
	}
}

func loadEnvVars() (map[string]string, error) {
	envVars := []string{
		envJiraServer,
		envJiraProject,
		envJiraUser,
		envJiraToken,
		envJiraReleaseName,
		envGithubRefName,
		envGithubRepo,
		envGithubToken,
		envPreviousVersion,
		envNewVersion,
	}
	checkEnvVars(envVars...)
	envMap := make(map[string]string)
	for _, v := range envVars {
		envMap[v] = os.Getenv(v)
	}
	return envMap, nil
}

func createHttpClient() *http.Client {
	return &http.Client{}
}

func handleRelease(client *http.Client, envMap map[string]string) {
	//Github Release Creation
	log.Println("Creating GitHub release")
	releaseURL, err := checkAndCreateGithubReleaseFunc(envMap[envGithubRepo], envMap[envNewVersion], envMap[envGithubToken], releaseNotesFile)
	if err != nil {
		log.Fatalf("Error creating GitHub release: %v", err)
	}
	log.Printf("GitHub release creation completed: %s", releaseURL)

	// Jira Release Creation
	releaseName := envMap[envJiraReleaseName] + " " + envMap[envNewVersion]

	if err := createJiraReleaseFunc(client, envMap[envJiraServer], envMap[envJiraProject], envMap[envJiraUser], envMap[envJiraToken], releaseName); err != nil {
		log.Fatalf("Error creating Jira release: %v", err)
	} else {
		log.Printf("Jira Release: %s", releaseName)
	}

	log.Println("Extracting changes")
	issueMaps, err := extractChangesFunc(releaseNotesFile)
	if err != nil {
		log.Fatalf("Error extracting changes: %v", err)
	}

	log.Println("Release Issues:")
	for _, issueMap := range issueMaps {
		issueID := extractIssueIDFunc(issueMap["title"], envMap[envJiraProject])
		if issueID == "" {
			log.Printf("No issue id: %s", issueMap["title"])
			continue
		}
		log.Printf("Updating issue %s with release %s", issueID, releaseName)
		err := addReleaseToIssueFunc(client, envMap[envJiraServer], envMap[envJiraUser], envMap[envJiraToken], releaseName, issueID)
		if err != nil {
			log.Printf("Error updating issue %s: %v", issueID, err)
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-h" || arg == "--help" {
			displayHelp()
			return
		}
	}

	envMap, err := loadEnvVars()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	client := createHttpClient()
	handleRelease(client, envMap)
}
