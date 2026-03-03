package backend

import (
	"os"
	"testing"

	"github.com/kohirens/stdlib/test"
)

const (
	fixtureDir = "testdata"
	tmpDir     = "tmp"
)

func TestMain(m *testing.M) {
	test.ResetDir(tmpDir, 0777)
	test.ResetDir(tmpDir+"/accounts", 0777)

	os.Exit(m.Run())
}
