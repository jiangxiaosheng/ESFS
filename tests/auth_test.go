package tests

import (
	"ESFS2.0/client/common"
	"ESFS2.0/message/protos"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	c, conn, err := common.GetAuthenticationClient()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &protos.RegisterRequest{
		Username:         "memeshe",
		Password:         "111",
		DefaultSecondKey: "000",
	}
	response, err := c.Register(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.ErrorMessage)
}

func TestLogin(t *testing.T) {
	c, conn, err := common.GetAuthenticationClient()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &protos.LoginRequest{
		Username: "memeshe",
		Password: "111",
	}
	response, err := c.Login(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.ErrorMessage)
}

func TestGetSecondKey(t *testing.T) {
	//key, err := dataserver.GetDefaultSecondKey("memeshe")
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(key)
	//}
}
