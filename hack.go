package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/eggsbenjamin/stepFnLocal/sfn"
	"github.com/eggsbenjamin/stepFnLocal/state"
)

func main() {
	defJSON := `
	{
  "Comment": "Anatwine line-details generator step function",
  "StartAt": "Linedetails_Validator",
  "States": {
    "Linedetails_Validator": {
			"InputPath" : "$",
      "Type": "Task",
      "Resource": "arn:aws:lambda:eu-west-1:459476646026:function:linedetails_validator_dev",
      "Next": "Linedetails_Generator"
    },
    "Linedetails_Generator": {
      "Type": "Task",
			"OutputPath" : "$.result",
      "Resource": "arn:aws:lambda:eu-west-1:459476646026:function:linedetails_generator_dev",
      "End": true
		}
  }
}
	`

	var def state.MachineDefinition
	if err := json.Unmarshal([]byte(defJSON), &def); err != nil {
		log.Fatal(err)
	}

	overrides := map[string]sfn.OverrideFn{
		"arn:aws:lambda:eu-west-1:459476646026:function:linedetails_generator_dev": func([]byte) ([]byte, error) {
			return []byte(`{ "result" : [1,2,3,4,5] }`), nil
		},
	}
	config := &aws.Config{
		Region: aws.String("eu-west-1"),
	}
	lambdaClient := lambda.New(session.Must(session.NewSession(config)))
	stateFactory := sfn.NewStateFactory(overrides, lambdaClient)

	stepFn, err := sfn.New(def, stateFactory)
	if err != nil {
		log.Fatal(err)
	}

	result, err := stepFn.StartExecution([]byte(`"zalandomp_14102018.zip"`))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status: %s\n", result.Status)
	log.Printf("input: %s\n", result.Input)
	log.Printf("output: %s\n", result.Output)
}
