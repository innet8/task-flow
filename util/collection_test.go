package util_test

import (
	"fmt"
	"testing"

	"workflow/util"
)

func TestExistsDuplicateInStringsArr(t *testing.T) {
	users := []string{"张三", "李皿", "王五"}
	res := util.ExistsDuplicateInStringsArr(users)
	fmt.Println(res)
}
