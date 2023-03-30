package util_test

import (
	"fmt"
	"log"
	"testing"

	"workflow/util"
)

func TestIsDoubleStr(t *testing.T) {
	str := "12.020"
	yes, err := util.IsDoubleStr(str)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("is double: %v\n", yes)
}
