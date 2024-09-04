package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func createJiraRelease(client *http.Client, server, project, user, token, name string) error {
	url := fmt.Sprintf("%s/rest/api/2/version", server)
	release := map[string]interface{}{
		"name":        name,
		"project":     project,
		"released":    false,
		"startDate":   time.Now().Format("2006-01-02"),
		"releaseDate": time.Now().Format("2006-01-02"),
	}
	body, err := json.Marshal(release)
	if err != nil {
		log.Printf("Error marshaling release data: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return err
	}
	req.SetBasicAuth(user, token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to create Jira release: %s", resp.Status)
		return fmt.Errorf("failed to create Jira release: %s", resp.Status)
	}

	log.Printf("Jira release created successfully: %s", name)
	return nil
}

func addReleaseToIssue(client *http.Client, server, user, token, release, issueID string) error {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", server, issueID)
	payload := map[string]interface{}{
		"update": map[string]interface{}{
			"fixVersions": []map[string]interface{}{
				{"add": map[string]string{"name": release}},
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+token)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update Jira issue: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}
