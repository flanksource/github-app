/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handler

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/flanksource/build-tools/pkg/junit"
	"github.com/google/go-github/v32/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type CheckSuiteHandler struct {
	githubapp.ClientCreator

	preamble string
}

func (h *CheckSuiteHandler) Handles() []string {
	return []string{"check_suite"}
}

func (h *CheckSuiteHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
 	var event github.CheckSuiteEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse check suite event payload")
	}
	if event.GetAction() != "completed" {
		return nil //we only want to process the results at completion and ignore anything else
	}
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return errors.Wrapf(err, "failed to get github client from installationID %s given in event", installationID)
	}
	repo := event.GetRepo()
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	// the check suite might be associated with multiple PRs, handle each one:
	for _,pr := range event.GetCheckSuite().PullRequests {
		prNum := pr.GetNumber()
		ctxPr, logger := githubapp.PreparePRContext(ctx, installationID, repo, prNum)

		// ☹️ no easy direct way to go from the check run to the workflow run
		// so we get this first page of workflow runs for this owner/repo
		// and look for a matching head SHA commit on the branch for this event
		wfrList, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, repoOwner, repoName, &github.ListWorkflowRunsOptions{
			Branch:      *event.CheckSuite.HeadBranch,
			ListOptions: github.ListOptions{},
		})
		if err != nil {
			return err
		}

		results :=  make([]string, 0, 20)

		for _, wfr := range wfrList.WorkflowRuns {
			if *wfr.HeadSHA != *event.CheckSuite.HeadSHA {
				continue // ignore the run if it isn't for our commit, got to next workflowrun
			}
			// cool now for this workflow run we get the artifacts
			artifactList, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, repoOwner, repoName, *wfr.ID,&github.ListOptions{})
			if err != nil {
				logger.Error().Err(err).Msg("failed to list workflowrun artifacts")
				continue //ignore error try next workflowrun
			}
			for _, artifact := range artifactList.Artifacts {
				if !strings.HasPrefix(*artifact.Name, "test-results") {
					continue // we only care about 'test-results*', skip this artifact
				}

				url, _, err := client.Actions.DownloadArtifact(ctx, repoOwner, repoName,*artifact.ID,true)
				if err != nil {
					logger.Error().Err(err).Msg("failed to get artifact download url")
					continue //ignore error, try next artifact
				}
				tmpfile, err := ioutil.TempFile("/tmp", "downloadedzip")
				defer os.Remove(tmpfile.Name())
				if err != nil {
					continue //ignore error, try next artifact
				}
				err = downloadFile(url.String(),tmpfile)
				contents, err := getUnzippedFileContents(tmpfile.Name(),"results.xml")
				results = append(results, contents)
			}
		}

		if len(results)> 0 {
			tr, err := junit.ParseJunitResultStrings(results...)
			if err != nil {
				return err
			}

			msg := fmt.Sprintf("Commit %s\n%s", *event.CheckSuite.HeadSHA,tr.GenerateMarkdown())
			// yeah, IssueComment. PRComments are review comments and have extra metadata
			prComment := github.IssueComment{
				Body: &msg,
			}

			if _, _, err := client.Issues.CreateComment(ctxPr, repoOwner, repoName, prNum, &prComment); err != nil {
				logger.Error().Err(err).Msg("Failed to comment on pull request")
			}
		}
	}
	return nil
}

// downloadFile will download a url to a given open local file.
// The file will be closed on completion.
func downloadFile(url string, f *os.File) error {
	// Get the data
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("dowload failed - HTTP response status was %v", response.StatusCode)
	}

	defer f.Close()

	// Write the body to file
	_, err = io.Copy(f, response.Body)
	return err
}

// getUnzippedFileContents opens a zip file specified by its file name `zipname`
// and returns the string contents of the file specified by `filename`
func getUnzippedFileContents(zipname string, filename string ) (string, error) {
	z, err := zip.OpenReader(zipname)
	if err != nil {
		return "", err
	}
	defer z.Close()

	for _, f := range z.File {
		fmt.Printf("Contents of %s: %s\n", zipname, f.Name)
		if f.Name == filename {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("error opening file %s in zip %s: %v", filename, zipname, err)
			}
			defer rc.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(rc)
			contents := buf.String()
			return contents, nil
		}
	}
	return "", fmt.Errorf("file %s not found", filename)
}