package main

import (
	"fmt"

	. "github.com/cydev/vk"
)

func main() {
	fmt.Println(Auth{Scope: NewScope(PermOffline, "messages"), ID: 3897553}.URL())
}
