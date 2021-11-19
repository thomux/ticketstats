package ticketstats

import (
	"os"
	"testing"
)

func TestSaveConfig(t *testing.T) {
	saveConfig()

	info, err := os.Stat("config.json")
	if err != nil {
		t.Fail()
	}
	if info.Size() == 0 || info.IsDir() {
		t.Fail()
	}
}
