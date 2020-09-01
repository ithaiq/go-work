package singleton

import (
	"fmt"
	"testing"
)

func TestGetInstance(t *testing.T) {
	for i := 0; i < 10; i++ {
		go GetInstance()
	}
	fmt.Scanln()
}
func TestGetInstance2(t *testing.T) {
	for i := 0; i < 10; i++ {
		go GetInstance2()
	}
	fmt.Scanln()
}