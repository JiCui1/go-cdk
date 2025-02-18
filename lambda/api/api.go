package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
  dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
  return ApiHandler {
    dbStore:  dbStore,
  }
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  var registerUser types.RegisterUser

  //map json to type
  err := json.Unmarshal([]byte(request.Body), &registerUser)

  if err != nil {
    return events.APIGatewayProxyResponse{
      Body: "Invalid Request",
      StatusCode: http.StatusBadRequest,
    }, err
  }

  if registerUser.Username == "" || registerUser.Password == "" {
    return events.APIGatewayProxyResponse{
      Body: "Invalid Request - fields empty",
      StatusCode: http.StatusBadRequest,
    }, err
  }

  userExists, err := api.dbStore.DoesUserExist(registerUser.Username) 
  if err != nil {
    return events.APIGatewayProxyResponse{
      Body: "Internal Server Error",
      StatusCode: http.StatusInternalServerError,
    }, err
  }

  if userExists {
    return events.APIGatewayProxyResponse{
      Body: "User already exists",
      StatusCode: http.StatusConflict,
    }, nil
  }

  user, err := types.NewUser(registerUser)
  if err != nil {
    return events.APIGatewayProxyResponse{
      Body: "Internal Server Error",
      StatusCode: http.StatusInternalServerError,
    }, nil
  }

  err = api.dbStore.InsertUser(user)
  if err != nil {
    return events.APIGatewayProxyResponse{
      Body: "Internal Server Error",
      StatusCode: http.StatusInternalServerError,
    }, err
  }

  return events.APIGatewayProxyResponse{
    Body: "Successfully Registered",
    StatusCode: http.StatusOK,
  }, nil
}
