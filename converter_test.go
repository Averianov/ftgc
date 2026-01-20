package ftgs

import (
	"os"
	"testing"

	sl "github.com/Averianov/cisystemlog"
)

// TestLogging ###############################################
func TestLogging(t *testing.T) {
	sl.CreateLogs("", "", 4, 0)
	sl.L.Debug("TestLogging start agent\n")

	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	_, err = f.Write([]byte("some text"))
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = ConvertDirectory("./", "./out", "txt")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	t.Errorf("%s", "check result")

	err = os.Remove("test.txt")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = os.Remove("./out/test.txt.go")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}
