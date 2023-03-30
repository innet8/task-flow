package util_test

import (
	"fmt"
	"testing"
)

type Phone interface {
	call()
}
type Nokia struct{}
type Apple struct{}

func (nokia Nokia) call() {
	fmt.Println("Nokia")
}
func (apple Apple) call() {
	fmt.Println("Apple")
}
func TestTest(t *testing.T) {
	var phone Phone
	phone = new(Nokia)
	phone.call()
	phone = new(Apple)
	phone.call()

}
