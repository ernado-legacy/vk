package main

import (
	"flag"

	"fmt"
	"github.com/cydev/vk"
	"os"
)

var (
	userID int
	token  string
)

func init() {
	flag.StringVar(&token, "token", "", "vk api token")
	flag.IntVar(&userID, "id", 0, "user id")
}

func main() {
	flag.Parse()
	api := vk.NewWithToken(token)
	groups, err := api.Groups.GetForUser(userID)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(2)
	}
	for _, group := range groups {
		fmt.Printf("%+v\n", group)
	}
}
