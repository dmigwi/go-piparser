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

	handler := func() {}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			p, err := NewParser(val.repoOwner, val.repo, val.dir, handler)
			if err != nil {
				t.Fatalf("expected no error but found: %v", err)
			}

			var actualRepo, actualOwner = types.DefaultRepo, types.DefaultRepoOwner
			if val.repo != "" {
				actualRepo = testRepo
			}

			if p.repoName != actualRepo {
				t.Fatalf("expected to find the %s repo but found %s ",
					actualRepo, p.repoName)
			}

			if val.repoOwner != "" {
				actualOwner = testROwner
			}

			if p.repoOwner != actualOwner {
				t.Fatalf("expected to find %s repo owner but found %s ",
					actualOwner, p.repoOwner)
			}

			if val.dir != "" && p.cloneDir != testDir && val.dir != invalidPath {
				t.Fatalf("expected to find %s path but found %s path",
					testDir, p.cloneDir)
			}

			if val.dir == invalidPath && !strings.Contains(p.cloneDir, DirPrefix) {
				t.Fatalf("expected a temporary folder to have been created but it wasn't")
			}

			if val.dir == "" && !strings.Contains(p.cloneDir, DirPrefix) {
				t.Fatalf("expected a temporary folder to have been created but it wasn't")
			}

			// clean up
			// drop the temporary folder created.
			if strings.Contains(p.cloneDir, DirPrefix) {
				os.RemoveAll(p.cloneDir)
			}
		})
	}
}
