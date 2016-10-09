package librevault

import (
	"fmt"
	"testing"
)

func TestLuhnModValue(t *testing.T) {
	test := func(data string, check byte) {
		t.Run(data+"="+string(check), func(t *testing.T) {
			result := luhnModValue([]byte(data))
			if result != check {
				t.Errorf("expected '%c' but got '%c'", check, result)
			}
		})
	}

	testData := []struct {
		data  string
		check byte
	}{
		{"fFfr3UMHoLqjoXPSaWHRySvijJrKJFPz3X8MtnNAzXT", 'Z'},
		{"ETdSkHLVeNPWfqLTsUDWPCUZqKCzF5qjFJtys8KPT3wdQxgtkxk1WTuvZbZx2WJQ9Pd1DBgs6deoBsTNEgFyXNMh", '1'},
		{"AMcu13VWLTfKZfJNxkm18PeRQfJ3jfp19SirnurWzXfh", 'V'},
	}
	for _, d := range testData {
		test(d.data, d.check)
	}
}

func TestAppendChecksum(t *testing.T) {
	test := func(data string, check byte) {
		expected := data + string(check)
		t.Run(data+"="+expected, func(t *testing.T) {
			result := appendChecksum(data)
			if result != expected {
				t.Errorf("expected '%s' but got '%s'", expected, result)
			}
		})
	}

	testData := []struct {
		data  string
		check byte
	}{
		{"fFfr3UMHoLqjoXPSaWHRySvijJrKJFPz3X8MtnNAzXT", 'Z'},
		{"ETdSkHLVeNPWfqLTsUDWPCUZqKCzF5qjFJtys8KPT3wdQxgtkxk1WTuvZbZx2WJQ9Pd1DBgs6deoBsTNEgFyXNMh", '1'},
		{"AMcu13VWLTfKZfJNxkm18PeRQfJ3jfp19SirnurWzXfh", 'V'},
	}
	for _, d := range testData {
		test(d.data, d.check)
	}
}

func TestValidateChecksum(t *testing.T) {
	test := func(data string, ok bool) {
		t.Run(fmt.Sprintf("%s=%t", data, ok), func(t *testing.T) {
			result := validateChecksum(data)
			if result != ok {
				t.Errorf("expected '%s' to be %t", data, ok)
			}
		})
	}

	testData := []struct {
		data string
		ok   bool
	}{
		{"fFfr3UMHoLqjoXPSaWHRySvijJrKJFPz3X8MtnNAzXTZ", true},
		{"fFfr3UMHoLqjoXPSaWHRySvijJrKJFPz3X8MtnNAzXTf", false},
		{"ETdSkHLVeNPWfqLTsUDWPCUZqKCzF5qjFJtys8KPT3wdQxgtkxk1WTuvZbZx2WJQ9Pd1DBgs6deoBsTNEgFyXNMh1", true},
		{"AMcu13VWLTfKZfJNxkm18PeRQfJ3jfp19SirnurWzXfhV", true},
	}
	for _, d := range testData {
		test(d.data, d.ok)
	}
}
