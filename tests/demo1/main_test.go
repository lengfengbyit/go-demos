package main

import (
	"fmt"
	"os"
	"testing"
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
