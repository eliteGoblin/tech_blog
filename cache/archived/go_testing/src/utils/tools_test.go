package utils

import (
	"testing"
	"bytes"
	"text/template"
	"fmt"
)

func ExampleDiv() {
	fmt.Println(Div(1, 2))
	// Output: 0.5
}

func ExampleRand() {
	arr := []int{2, 3, 4, 5, 1, 0}
	for i := 0; i < len(arr); i ++ {
		fmt.Println(i)
	}
	// Unordered output:
	// 5
	// 4
	// 3
	// 2
	// 1
	// 0
}

func TestDivFraction(t *testing.T) {
	res := Div(1, 2)
	if res != 0.5 {
		t.Errorf("Div was incorrect, got: %d, want: %f.", res, 0.1)
	}
}



func TestDivZero(t *testing.T) {
	Div(1, 0)
}

func BenchmarkDiv(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		Div(1, 2)
	}
}

func BenchmarkTemplateParallel(b *testing.B) {
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}