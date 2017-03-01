package mix_test

import (
	"github.com/Andoryuuta/renlib/io/mix"
	"testing"
)

func TestParseFile(t *testing.T) {
	// Try to parse the mix file
	out, err := mix.ParseFile("data/test500b500map.mix")
	if err != nil {
		t.Error(err.Error())
	}

	// Make a slice of anon struct with the expected results
	expected := []struct {
		Name     string
		ID       uint32
		DataSize int
	}{
		{"test500b500map.lsd", 0x23A7EE8F, 4259661},
		{"test500b500map.ldd", 0x26246A19, 1822},
		{"test500b500map.dep", 0x2BF6DE9D, 8},
		{"test500b500map.ddb", 0xC1549E94, 24},
	}

	// Check if the amount of files returned and the amount of expected files match
	if len(out.Files) != len(expected) {
		t.Error("len(out) != len(expected)")
	}

	// Compare the file details to the expected ones
	for i := 0; i < len(out.Files); i++ {
		cOut := out.Files[i]
		cExp := expected[i]
		if cOut.Name != cExp.Name {
			t.Errorf("Output name (%s) does not match expected name (%s)", cOut.Name, cExp.Name)
		}
		if cOut.ID != cExp.ID {
			t.Errorf("Output ID (%s) does not match expected ID (%s)", cOut.ID, cExp.ID)
		}
		if len(cOut.Data) != cExp.DataSize {
			t.Errorf("Output size (%s) does not match expected size (%s)", len(cOut.Data), cExp.DataSize)
		}
	}
}
