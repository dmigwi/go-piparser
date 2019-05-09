package gitlib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dmigwi/go-piparser/proposals/types"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/diff"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

// Lib defines the fields required by the git tool data source.
type Lib struct {
	cloneDir  string
	repoName  string
	repoOwner string
	repo      *git.Repository
}

// NewDataSource returns a data source the uses the git golang library to
// query raw proposals data and process it into an array of proposal history.
func NewDataSource(owner, repo, cloneDir string) *Lib {
	return &Lib{cloneDir: cloneDir, repoName: repo, repoOwner: owner}
}

func (l *Lib) SetUpEnv() (err error) {
	// full clone directory: includes the expected repository name.
	workingDir := filepath.Join(l.cloneDir, types.CloneRepoAlias)
	completeRemoteURL := fmt.Sprintf(types.RemoteURL, l.repoOwner, l.repoName)

	l.repo, err = git.PlainClone(workingDir, false, &git.CloneOptions{
		URL: completeRemoteURL,
	})

	switch err {
	case nil:
		// The git repo clone complete successfull.
		return nil
	case git.ErrRepositoryAlreadyExists:
		// git.ErrRepositoryAlreadyExists error was returned.
		l.repo, err = git.PlainOpen(workingDir)
		if err != nil {
			return fmt.Errorf("git.PlainOpen: %v", err)
		}

		w, err := l.repo.Worktree()
		if err != nil {
			return fmt.Errorf("fetching the work tree failed: %v", err)
		}

		remote, err := l.repo.Remote(types.RemoteURLRef)
		if err != nil {
			return fmt.Errorf("fetching the remote url ref failed: %v", err)
		}

		if strings.Contains(remote.String(), completeRemoteURL) {
			// Pull the latest changes from the origin remote and merge into the current branch
			err = w.Pull(&git.PullOptions{RemoteName: types.RemoteURLRef})

			switch err {
			case nil, git.NoErrAlreadyUpToDate:
				return nil
			default:
				return fmt.Errorf("updates Pull failed: %v", err)
			}
		}
		// Incorrect remote reference URL was found. Drop the previous work space
		// and make a fresh clone.
		if err = os.RemoveAll(workingDir); err != nil {
			return fmt.Errorf("dropping the workingDir failed: %s", err)
		}
		l.repo, err = git.PlainClone(workingDir, false, &git.CloneOptions{
			URL: completeRemoteURL,
		})
		return err

	default:
		return err
	}
}

func (l *Lib) PullData(proposalToken string, since ...time.Time) ([]*types.History, error) {
	options := new(git.LogOptions)

	if len(since) > 0 && !since[0].IsZero() {
		los, ok := l.repo.Storer.(storer.LooseObjectStorer)
		if !ok {
			return nil, git.ErrLooseObjectsNotSupported
		}

		err := los.ForEachObjectHash(func(hash plumbing.Hash) error {
			timestamp, err := los.LooseObjectTime(hash)
			if err != nil {
				return err
			}

			if timestamp.Before(since[0]) {
				options.From = hash
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		// retrieves the branch pointed by HEAD
		ref, err := l.repo.Head()
		if err != nil {
			return nil, fmt.Errorf("retrieving branch pointed by Head failed: %v ", err)
		}
		options.From = ref.Hash()
	}

	// retrieves the commit history
	cIter, err := l.repo.Log(options)
	if err != nil {
		return nil, fmt.Errorf("retrieving commit history failed: %v ", err)
	}

	var items []*types.History

	// just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fromTree, err := c.Tree()
		if err != nil {
			return err
		}

		toTree := &object.Tree{}
		if c.NumParents() != 0 {
			firstParent, err := c.Parents().Next()
			if err != nil {
				return err
			}

			toTree, err = firstParent.Tree()
			if err != nil {
				return err
			}
		}

		patch, err := toTree.Patch(fromTree)
		if err != nil {
			return err
		}

		item := types.History{
			Author:    c.Author.String(),
			CommitSHA: c.Hash.String(),
			Date:      c.Committer.When,
		}

		var changes []*types.File
		for _, filePatches := range patch.FilePatches() {
			// ignore empty patches (binary files, submodule refs updates)
			if len(filePatches.Chunks()) == 0 {
				continue
			}

			if proposalToken != "" {
				_, toFile := filePatches.Files()
				if !strings.Contains(toFile.Path(), proposalToken) {
					continue
				}
			}

			fmt.Printf(" >>>> <<< %+v \n", filePatches)

			for _, chunk := range filePatches.Chunks() {
				if chunk.Type() == diff.Add {
					// Add the square brackets and commas to complete the JSON string array format.
					filePatch := "[" + types.ReplaceAny(chunk.Content(), "}{", "},{") + "]"

					var v types.Votes

					if err = json.Unmarshal([]byte(filePatch), &v); err != nil {
						return fmt.Errorf("Unmarshalling File failed: %v %s", err, filePatch)
					}

					if len(v) > 0 {
						changes = append(changes, &types.File{proposalToken, v})
					}
				}
			}
		}

		item.Patch = changes

		items = append(items, &item)

		return nil
	})

	if err == io.EOF || err == nil {
		return nil, nil
	}

	fmt.Printf(" >>>>>>> %v \n", items)

	return items, fmt.Errorf("retrieving patch failed: %v ", err)
}

// FetchProporties returns the set repo owner, repo name and the clone directory.
func (l *Lib) FetchProporties() (owner, name, cloneDir string) {
	return l.repoOwner, l.repoName, l.cloneDir
}
