package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
  dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
  return ApiHandler {
    dbStore:  dbStore,
  }
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {

  if event.Username == "" || event.Password == "" {
    return fmt.Errorf("request has empty params")
  }

  userExists, err := api.dbStore.DoesUserExist(event.Username) 
  if err != nil {
    return fmt.Errorf("this is an error: %w", err)
  }

  if userExists {
    return fmt.Errorf("user already exists")
  }

  err = api.dbStore.InsertUser(event)
  if err != nil {
    return fmt.Errorf("failed to insert user: %w", err)
  }

  return nil
}
