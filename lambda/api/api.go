package api

import (
	"encoding/json"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
  "fmt"
  "time"
)

type UserHandler struct {
  userStore database.UserStore
}

type BlogHandler struct {
  blogStore database.BlogStore
}

func NewUserHandler(userStore database.UserStore) UserHandler {
  return UserHandler {
    userStore:  userStore,
  }
}

func NewBlogHandler(blogStore database.BlogStore) BlogHandler {
  return BlogHandler {
    blogStore:  blogStore,
  }
}

func (api BlogHandler) CreateBlogHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  var newBlog types.Blog

  err := json.Unmarshal([]byte(request.Body), &newBlog)

  if err != nil {
    return events.APIGatewayProxyResponse{
      Body: "Invalid Request",
      StatusCode: http.StatusBadRequest,
    }, err
  }

  if newBlog.Title == "" || newBlog.Description == "" || newBlog.Content == "" {
    return events.APIGatewayProxyResponse{
      Body: "Invalid Request - fields empty",
      StatusCode: http.StatusBadRequest,
    }, err
  }

  newBlog.Slug = types.Slugify(newBlog.Title)

  date := time.Now()
  newBlog.CreatedAt = date.Format("Jan 2, 2025")


  err = api.blogStore.InsertBlog(newBlog)

  if err != nil {
    return events.APIGatewayProxyResponse{
      Body: "Internal Server Error",
      StatusCode: http.StatusInternalServerError,
    }, err
  }

  return events.APIGatewayProxyResponse{
    Body: "Successfully Created Blog",
    StatusCode: http.StatusOK,
  }, nil
}

func (api UserHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

  userExists, err := api.userStore.DoesUserExist(registerUser.Username) 
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

  err = api.userStore.InsertUser(user)
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

func (api UserHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  type LoginRequest struct {
    Username string
    Password string
  }

  var loginRequest LoginRequest

  err := json.Unmarshal([]byte(request.Body), &loginRequest)

  if err != nil {
    return events.APIGatewayProxyResponse{
      Body: "Invalid Request",
      StatusCode: http.StatusBadRequest,
    }, err
  }

  user, err := api.userStore.GetUser(loginRequest.Username)

  if err != nil {
    return events.APIGatewayProxyResponse{
      Body: "Internal Server Error",
      StatusCode: http.StatusInternalServerError,
    }, err
  }

  if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
    return events.APIGatewayProxyResponse{
      Body: "Invalid Credentials",
      StatusCode: http.StatusBadRequest,
    }, err
  }

  accessToken := types.CreateToken(user)
  successMsg := fmt.Sprintf(`{"access-token": "%s"}`, accessToken)

  return events.APIGatewayProxyResponse{
    Body: successMsg,
    StatusCode: http.StatusOK,
  }, nil
}
