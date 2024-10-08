package main

import (
	"fmt"
	"os"
	"testing"
	"unicode"
)

func TestMain(m *testing.M) {
	fmt.Println("Before all tests")

	exitCode := m.Run()

	fmt.Println("After all tests")
	os.Exit(exitCode)
}

func TestDemo(t *testing.T) {
	t.Log("test demo")
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{
			a:    0,
			b:    0,
			want: 0,
		},
		{
			a:    -1,
			b:    1,
			want: 0,
		},
		{
			a:    100,
			b:    -99,
			want: 1,
		},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%d add %d", test.a, test.b)
		t.Run(name, func(t *testing.T) {
			r := Add(test.a, test.b)
			if r != test.want {
				t.Errorf("Add(%d, %d) = %d; want %d", test.a, test.b, r, test.want)
			}
		})
	}
}

// countEnglishChars 返回字符串中英文字符的数量
func countEnglishChars(s string) int {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) && (unicode.IsUpper(char) || unicode.IsLower(char)) {
			count++
		}
	}
	return count
}

func TestCountEnglishChars(t *testing.T) {
	testString := "Hello, 世界! This is a test."
	count := countEnglishChars(testString)
	fmt.Printf("Number of English characters in \"%s\": %d\n", testString, count)
}
