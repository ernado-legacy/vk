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
	fields := vk.GroupGetFields{
		Fields:  "sex",
		Offset:  0,
		GroupID: 26188163,
	}
	users, n, err := api.Groups.GetBatch(fields)
	fmt.Println(n)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(2)
	}
	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}
}
