package gittool

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dmigwi/go-piparser/proposals/types"
)

const (
	// gitCmd defines the prefix string for all commands issued to the git
	// command line interface.
	gitCmd = "git"

	// listCommitsArg defines the git command line argument that lists the repo
	// commit history in a reverse chronoligical order. The oldest commit is
	// listed as the last.
	listCommitsArg = "log"

	// commitPatchArg is an optional argument that is added to show the commit
	// history with a patch (changes made) field included.
	commitPatchArg = "-p"

	// cloneArg is the argument added between the git prefix command and the
	// repository download URL. It is used to download the repository into the
	// underlying platform in the set working directory.
	cloneArg = "clone"

	// versionArg lists the version of the active git installation. It is
	// primarily used to check if git is installed on the underlying platform.
	versionArg = "--version"

	// pullChangesArg is an argument that helps pull latest changes from the
	// remote repository set. git pull = git fetch + git merge.
	pullChangesArg = "pull"

	// remoteURLRef references the remote URL used to clone the repository.
	// 'origin' is the default set.
	remoteURLRef = "origin"

	// sinceArg with syntax "--since <date>" returns commits more recent than a
	// specific date.
	sinceArg = "--since"

	// By default, the commits are shown in reverse chronological order. Using
	// the reverse order argument ensures that all commits returned are listed
	// in chronological order.
	reverseOrder = "--reverse"

	// trackedRemoteURL is part of the arguments that help retrieve the clone url
	// of the tracked repository.
	trackedRemoteURL = "get-url"

	// remoteDef defines the remote argument
	remoteDef = "remote"
)

// Tool defines the fields required by the git tool data source.
type Tool struct {
	cloneDir  string
	repoName  string
	repoOwner string
}

// NewDataSource returns a data source the uses the git commandline tool to
// query raw proposals data and process it into an array of proposal history.
func NewDataSource(owner, repo, cloneDir string) *Tool {
	return &Tool{cloneDir, repo, owner}
}

// SetUpEnv pulls changes from github if they exists or otherwise it
// clones the repository. It also ensures that a working git commandline tool
// is installed in the underlying platform and has the minimum version required.
// If the required repo was cloned earlier, only the latest changes are pulled
// otherwise a fresh clone is made. Should an error occurs while pulling
// updates, the old repo is dropped and a fresh clone made.
func (t *Tool) SetUpEnv() error {
	// check if git exists by checking the git installation version.
	versionStr, err := t.readCommandOutput(gitCmd, versionArg)
	if err != nil {
		return fmt.Errorf("checking git version(%s %s) failed: %v",
			gitCmd, versionArg, err)
	}

	// Check if a valid git version exists. A minimum of v1.5.1 is required.
	if err = types.IsGitVersionSupported(versionStr); err != nil {
		// invalid/unsupported git version must have been found.
		return err
	}

	// full clone directory: includes the expected repository name.
	workingDir := filepath.Join(t.cloneDir, types.CloneRepoAlias)
	_, err = os.Stat(workingDir)

	completeRemoteURL := fmt.Sprintf(types.RemoteURL, t.repoOwner, t.repoName)

	switch {
	case !os.IsNotExist(err):
		// The working directory was found thus check if the tracked repo is the
		// same as the required one.
		trackedRepo, err := t.readCommandOutput(gitCmd, remoteDef, trackedRemoteURL, remoteURLRef)

		// If the required tracked repo was found initiate the updates fetch process
		if err == nil && types.IsMatching(trackedRepo, completeRemoteURL) {
			if err = t.execCommand(gitCmd, pullChangesArg, remoteURLRef); err == nil {
				return nil
			}
		}

		// git pull command failed or the required tracked repo wasn't found.
		// Drop the old repo and proceed to clone the repo a fresh.
		if err = os.RemoveAll(workingDir); err != nil {
			return err
		}

		fallthrough

	default:
		// The required working directory could not be found or the repo update
		// process failed. Clone the remote repository into the clone directory.
		err = t.execCommand(gitCmd, cloneArg, completeRemoteURL, types.CloneRepoAlias)
		if err != nil {
			return fmt.Errorf("failed to clone %s : %v", completeRemoteURL, err)
		}
	}

	return nil
}

// PullData fetches the data from github with the provided arguments. It returns
// the patch data converted into an array of types.History.
func (t *Tool) PullData(proposalToken string, since ...time.Time) ([]*types.History, error) {
	args := []string{listCommitsArg, reverseOrder, commitPatchArg}

	// Append the proposal token limiting argument if it exists.
	if proposalToken != "" {
		args = append(args, proposalToken)
	}

	nilTime := time.Time{}
	var timestamp time.Time

	// Append the time limiting argument if it exists.
	if len(since) > 0 && since[0] != nilTime {
		args = append(args, []string{sinceArg,
			since[0].Format(types.CmdDateFormat)}...)
		timestamp = since[0]
	}

	// Fetch the data via git cmd.
	patchData, err := t.readCommandOutput(gitCmd, args...)
	if err != nil {
		return nil, fmt.Errorf("fetching proposal(s) history failed: %v", err)
	}

	if nilTime == timestamp {
		return constructHistory(patchData)
	}
	return constructHistory(patchData, timestamp)
}

// FetchProporties returns the set repo owner, repo name and the clone directory.
func (t *Tool) FetchProporties() (owner, name, cloneDir string) {
	return t.repoOwner, t.repoName, t.cloneDir
}

// readCommandOutput reads the std output messages of the run command.
func (t *Tool) readCommandOutput(cmdName string, args ...string) (string, error) {
	cmd, err := t.processCommand(cmdName, args...)
	if err != nil {
		return "", err
	}

	stdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", formatError(cmdName, args, err)
	}

	return string(stdOutput), nil
}

// execCommand executes commands that do not return necessary std output messages.
func (t *Tool) execCommand(cmdName string, args ...string) error {
	cmd, err := t.processCommand(cmdName, args...)
	if err != nil {
		return err
	}

	err = cmd.Run()

	return formatError(cmdName, args, err)
}

// processCommand checks if an empty command prefix was provided. It also sets the
// the working directory.
func (t *Tool) processCommand(cmdName string, args ...string) (*exec.Cmd, error) {
	if cmdName == "" {
		return nil, fmt.Errorf("missing command")
	}
	cmd := exec.Command(cmdName, args...)

	// set the working directory.
	cmd.Dir = t.workingDir()
	return cmd, nil
}

// workingDir return (cloneDir + repoName) directory path if the target repo
// exists otherwise returns cloneDir as the working directory.
func (t *Tool) workingDir() string {
	dir := filepath.Join(t.cloneDir, types.CloneRepoAlias)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = t.cloneDir
	}
	return dir
}

// formatError replaces "exit status 128" with "git args... failed to execute."
func formatError(cmd string, args []string, err error) error {
	if err == nil {
		return err
	}
	str := fmt.Sprintf("%s command with %v failed to execute", cmd, args)
	return fmt.Errorf(types.ReplaceAny(err.Error(), "exit status 128", str))
}

func constructHistory(patch string, since ...time.Time) (items []*types.History, err error) {
	data := strings.Split(patch, "commit")

	for _, entry := range data {
		// strings.Split returns some split strings as empty or with just
		// whitespaces and other special charactes. This happens when
		// the seperating argument is the first in the source string or is
		// surrounded by whitespaces and other special characters.
		if len(strings.TrimSpace(entry)) == 0 {
			continue
		}

		var h types.History

		// entry string is not a valid JSON string format thus the use of a
		// customized unmarshaller.
		if err = types.CustomUnmashaller(&h, entry, since...); err != nil {
			return nil, fmt.Errorf("CustomUnmashaller failed: %v", err)
		}

		// Do not store any empty history data.
		if len(h.Patch) == 0 || h.Author == "" || h.CommitSHA == "" {
			continue
		}

		items = append(items, &h)
	}
	return
}
