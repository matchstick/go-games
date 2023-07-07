package gopark

import "testing"

func TestFibonacci(t *testing.T) {
	expected := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55,
		89, 144, 233, 377, 610, 987, 1597, 2584, 4181}

	f := fibonacci()

	for i := 0; i < 20; i++ {
		result := f()
		if result != expected[i] {
			t.Errorf("ERROR: Expected %d got %d at index %d", result, expected[i], i)
		}
	}
}
