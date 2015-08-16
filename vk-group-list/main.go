package main

import (
	"fmt"
	"github.com/cydev/vk"
	"github.com/spf13/viper"
	"time"
)

var (
	groupID int
	token   string
)

func init() {
	viper.SetEnvPrefix("vk")
	viper.BindEnv("token")
	viper.BindEnv("id")
}

func getAllUsers(api *vk.Client) (err error) {
	var (
		gotUsers int
		n        int
		users    []vk.User
	)
	fields := vk.GroupGetFields{
		Fields:  vk.UserFields,
		Offset:  0,
		GroupID: groupID,
	}
	for {
		fields.Offset = gotUsers
		users, n, err = api.Groups.GetBatch(fields)
		if err != nil {
			return err
		}
		gotUsers += len(users)
		fmt.Println("got", gotUsers, "of", n)
		if gotUsers >= n {
			// got all users
			fmt.Println("got all", gotUsers)
			return nil
		}
	}
}

func main() {
	viper.AutomaticEnv()
	groupID = viper.GetInt("id")
	token = viper.GetString("token")
	fmt.Println(groupID, token)
	api := vk.NewWithToken(token)
	start := time.Now()
	getAllUsers(api)
	end := time.Now()
	fmt.Println("users loaded", end.Sub(start))
}
