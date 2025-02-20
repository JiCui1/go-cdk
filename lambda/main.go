package main

import (
	"fmt"
	"lambda-func/app"
	"net/http"
	"lambda-func/middleware"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
  Username string `json:"username"`
}

func HandleRequest(event MyEvent) (string, error) {
  if (event.Username == "") {
    return "", fmt.Errorf("username cannot be empty")
  }

  return fmt.Sprintf("succssfully called by - %s", event.Username), nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  return events.APIGatewayProxyResponse{
    Body: "This is a protected path",
    StatusCode: http.StatusOK,
  }, nil
}

func main() {
  myApp := app.NewApp()

  lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    switch request.Path {
      // case "/register":
      //   return myApp.ApiHandler.RegisterUserHandler(request)
      case "/login":
        return myApp.ApiHandler.LoginUser(request)
      case "/protected":
        // this syntax is chaining functions, this is how next function is called in the chain
        return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)

      default:
        return events.APIGatewayProxyResponse{
          Body: "Not Found",
          StatusCode: http.StatusNotFound,
        }, nil
    }
  })
}
