package proposals

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/dmigwi/go-piparser/proposals/types"
)

var testDir string

func TestMain(m *testing.M) {
	var err error
	testDir, err = ioutil.TempDir(os.TempDir(), "test-DIR-")
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.RemoveAll(testDir)
	os.Exit(code)
}

// TestNewParser tests the setting of the default value.
func TestNewParser(t *testing.T) {
	type testData struct {
		repoOwner, repo, dir string
	}

	testRepo := "samaki"
	testROwner := "dmigwi"
	invalidPath := "/invalid/path"

	testROwnerWSpaces := " " + testROwner
	testRepoWSpaces := "   " + testRepo + " "
	testDirWSpaces := "  " + testDir + "  "

	td := []testData{
		{"", "", ""},
		{"", "", testDir},
		{"", testRepo, ""},
		{testROwner, "", testDir},
		{testROwner, testRepoWSpaces, ""},
		{testROwnerWSpaces, testRepo, testDirWSpaces},
		{testROwnerWSpaces, testRepoWSpaces, invalidPath},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			p, err := NewParser(val.repoOwner, val.repo, val.dir, true)
			if err != nil {
				t.Fatalf("expected no error but found: %v", err)
			}

			var actualRepo, actualOwner = types.DefaultRepo, types.DefaultRepoOwner
			if val.repo != "" {
				actualRepo = testRepo
			}

			setRepoOwner, setRepoName, setCloneDir := p.source.FetchProporties()

			if setRepoName != actualRepo {
				t.Fatalf("expected to find the %s repo but found %s ",
					actualRepo, setRepoName)
			}

			if val.repoOwner != "" {
				actualOwner = testROwner
			}

			if setRepoOwner != actualOwner {
				t.Fatalf("expected to find %s repo owner but found %s ",
					actualOwner, setRepoOwner)
			}

			if val.dir != "" && setCloneDir != testDir && val.dir != invalidPath {
				t.Fatalf("expected to find %s path but found %s path",
					testDir, setCloneDir)
			}

			if val.dir == invalidPath && !strings.Contains(setCloneDir, types.DirPrefix) {
				t.Fatalf("expected a temporary folder to have been created but it wasn't")
			}

			if val.dir == "" && !strings.Contains(setCloneDir, types.DirPrefix) {
				t.Fatalf("expected a temporary folder to have been created but it wasn't")
			}

			// clean up
			// drop the temporary folder created.
			if strings.Contains(setCloneDir, types.DirPrefix) {
				os.RemoveAll(setCloneDir)
			}
		})
	}
}
