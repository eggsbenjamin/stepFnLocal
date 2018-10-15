//go:generate mockgen -package lambda -source=lambda.go -destination lambda_mock.go

package lambda

import (
	"github.com/aws/aws-sdk-go/service/lambda"
)

// Client defines the lambda client interface
type Client interface {
	Invoke(*lambda.InvokeInput) (*lambda.InvokeOutput, error)
}
