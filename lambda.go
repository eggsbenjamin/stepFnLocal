package main

import (
	"log"
	"net/rpc"
	"os"
	"time"

	awslambda "github.com/aws/aws-lambda-go/lambda"
	awslambdamsg "github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/pkg/errors"
)

func main() {
	const port = "9898"

	helloWorld := NewLambda(
		func() (string, error) {
			return "hello world", nil
		},
		port,
	)

	go func() {
		helloWorld.Listen()
	}()

	time.Sleep(time.Second * 1)
	res, err := helloWorld.Invoke(&InvokeRequest{
		Payload: []byte(`""`),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", res.Payload)
}

type InvokeRequest awslambdamsg.InvokeResponse

type InvokeResponse awslambdamsg.InvokeResponse

type Lambda interface {
	Listen()
	Invoke(*InvokeRequest) (*InvokeResponse, error)
}

type lambda struct {
	handler interface{}
	port    string
}

func NewLambda(handler interface{}, port string) Lambda {
	return &lambda{
		handler: handler,
		port:    port,
	}
}

func (l *lambda) Listen() {
	os.Setenv("_LAMBDA_SERVER_PORT", l.port)
	awslambda.Start(l.handler)
}

func (l *lambda) Invoke(req *InvokeRequest) (*InvokeResponse, error) {
	c, err := rpc.Dial("tcp", "127.0.0.1:"+l.port)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to rpc endpoint")
	}

	var res InvokeResponse
	if err := c.Call("Function.Invoke", req, &res); err != nil {
		return nil, errors.Wrap(err, "error invoking lambda")
	}

	return &res, nil
}

func reverse(v string) (string, error) {
	var out string
	for i := len(v) - 1; i >= 0; i-- {
		out += string(v[i])
	}

	return out, nil
}
