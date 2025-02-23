package app

import (
  "lambda-func/api"
  "lambda-func/database"
)

type App struct {
  UserHandler api.UserHandler
  BlogHandler api.BlogHandler
}

func NewApp() App {
  db := database.NewDynamoDBClient()
  userHandler := api.NewUserHandler(db.UserStore())
  blogHandler := api.NewBlogHandler(db.BlogStore())

  return App {
    UserHandler: userHandler,
    BlogHandler: blogHandler,
  }
}
