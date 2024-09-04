package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func checkAndCreateGithubRelease(repo, refName, token, releaseNotesFile string) (string, error) {
	// Check if the release exists
	checkCmd := exec.Command("gh", "api", "--method", "GET", fmt.Sprintf("/repos/%s/releases/tags/%s", repo, refName), "--jq", ".body")
	output, err := checkCmd.Output()
	if err == nil {
		log.Printf("Release notes found for tag %s", refName)
		err = os.WriteFile("release_notes.md", output, 0644)
		if err != nil {
			return "", fmt.Errorf("error writing release notes to file: %v", err)
		}
		return string(output), nil
	}

	// If the release does not exist, create it
	createCmd := exec.Command("gh", "api", "--method", "POST", fmt.Sprintf("/repos/%s/releases", repo),
		"-f", fmt.Sprintf("tag_name='%s'", refName),
		"-f", fmt.Sprintf("name='%s'", refName),
		"-F", "generate_release_notes=true",
		"--jq", ".body")
	output, err = createCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error creating GitHub release: %v", err)
	}

	err = os.WriteFile(releaseNotesFile, output, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing release notes to file: %v", err)
	}
	fmt.Println("Release notes written to release_notes.md")

	log.Printf("GitHub release created successfully for tag %s", refName)
	return string(output), nil
}