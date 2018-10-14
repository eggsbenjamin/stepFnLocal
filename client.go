package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

func main() {

}

func invokeLambda(port string, out interface{}) {
	c, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatalln(err)
	}

	var result messages.InvokeResponse
	request := &messages.InvokeRequest{
		Payload: []byte(os.Getenv("V")),
	}
	err = c.Call("Function.Invoke", request, &result)

	if err != nil {
		log.Fatal(err)
	}
	if result.Error != nil {
		log.Fatal(result.Error.Message)
	}

	fmt.Println("Function.Invoke =", string(result.Payload))

}
