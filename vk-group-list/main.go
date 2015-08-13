package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	. "github.com/cydev/vk"
	"strings"
	"os"
)

type BlankResponse struct {
	Error `json:"error"`
}

func main() {
	token := os.Args[1]
	group := os.Args[2]

	var requests []string
	for i := 0; i < 20; i++ {
		args := url.Values{}
		args.Add("count", "1000")
		args.Add("offset", strconv.Itoa(i * 1000))
		args.Add("group_id", group)

		r := Request{
			Method: "groups.getMembers",
			Values: args,
		}
		requests = append(requests, r.JS())
	}

	code := fmt.Sprintf("return [%s];", strings.Join(requests, ","))

	fmt.Println(code)
	args := url.Values{}
	args.Add("code", code)

	req := Request{
		Token:  token,
		Method: "execute",
		Values: args,
	}

	res, err := DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Error)
	log.Println(res.Response)
}
