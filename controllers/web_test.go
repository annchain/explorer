package controllers

import (
	"testing"
	"regexp"
	"fmt"
)

func TestWebController_Search(t *testing.T) {
	hash := "1"

	reg := regexp.MustCompile(`^[1-9][0-9]*$`)

	fmt.Println(reg.MatchString(hash))
}

