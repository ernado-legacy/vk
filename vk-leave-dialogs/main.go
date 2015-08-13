package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	. "github.com/cydev/vk"
)

type BlankResponse struct {
	Error `json:"error"`
}

func main() {
	ownerID, _ := strconv.Atoi(os.Args[1])
	token := os.Args[2]
	offset := os.Args[3]

	args := url.Values{}
	args.Add("count", "200")
	args.Add("offset", offset)

	req := Request{
		Token:  token,
		Method: "messages.getDialogs",
		Values: args,
	}

	type Dialogs struct {
		Items []struct {
			Message struct {
				ID    int `json:"chat_id"`
				Admin int `json:"admin_id"`
				Out   int `json:"out"`
			} `json:"message"`
		} `json:"items"`
	}

	var dialogs Dialogs

	res, err := DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if err := res.To(&dialogs); err != nil {
		log.Fatal(err)
	}
	fmt.Println(dialogs)
	for _, dialog := range dialogs.Items {
		if ownerID != dialog.Message.Admin {
			continue
		}
		fmt.Println(dialog)

		if dialog.Message.Out != 1 {
			fmt.Println("leaving chat")
			args := url.Values{}
			args.Add("chat_id", strconv.Itoa(dialog.Message.ID))
			args.Add("user_id", "214321467")
			args.Add("count", "10000")
			response := BlankResponse{}
			res, err := DefaultClient.Do(Request{
				Token:  token,
				Method: "messages.removeChatUser",
				Values: args,
			})
			if err != nil {
				log.Fatal(err)
			}
			if err := res.To(&response); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second / 3)
		}
		//
		// fmt.Println("deleting messages")
		// args = url.Values{}
		// args.Add("chat_id", strconv.Itoa(dialog.Message.ID))
		// args.Add("count", "10000")
		// response := BlankResponse{}
		// if err := DefaultClient.Do(Request{
		// 	Token:  token,
		// 	Method: "messages.deleteDialog",
		// 	Values: args,
		// }, &response); err != nil {
		// 	fmt.Println(err)
		// }
		// time.Sleep(time.Second / 3)

	}
}
