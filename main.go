package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func main() {
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%q\n", r.URL.String())
		for _, h := range r.Header {
			log.Println(h)
		}
		rawBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s\n", rawBody)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", srv))
	}()

	config := &aws.Config{
		Region: aws.String("eu-west-1"),
		//Endpoint: aws.String("http://localhost:8080"),
	}

	lambdaClient := lambda.New(session.Must(session.NewSession(config)))
	invokeOutput, err := lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String("linedetails_generator_dev"),
		Payload:      []byte(`{ "test" : "hello!"}`),
	})
	if err != nil {
		log.Fatalf("error invoking lambda: %q", err)
	}
	log.Printf("%q %s\n", invokeOutput, invokeOutput.Payload)
}
