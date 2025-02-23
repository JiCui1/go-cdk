package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "lambda-func/types"
  "fmt"
)

const (
  USERS_TABLE="userTable"
  BLOGS_TABLE="blogsTable"
)

type UserStore interface {
  DoesUserExist(username string) (bool, error)
  InsertUser(user types.User) error
  GetUser(username string) (types.User, error)
}

type BlogStore interface {
  // GetBlog(BlogSlug string) (types.Blog, error)
  InsertBlog(blog types.Blog)  error
  // GetAllBlogs() ([]types.Blog, error)
}

type DynamoDBClient struct {
  userStore UserStore
  blogStore BlogStore
}

func (d *DynamoDBClient) UserStore() UserStore {
    return d.userStore
}

func (d *DynamoDBClient) BlogStore() BlogStore {
    return d.blogStore
}

type DynamoUserStore struct {
  databaseStore *dynamodb.DynamoDB
}

type DynamoBlogStore struct {
  databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
  dbSession := session.Must(session.NewSession())
  db := dynamodb.New(dbSession)

  return DynamoDBClient{
    userStore: &DynamoUserStore{databaseStore: db},
    blogStore: &DynamoBlogStore{databaseStore: db},
  }
}

func (u DynamoBlogStore) InsertBlog(blog types.Blog) error {

  item := &dynamodb.PutItemInput{
    TableName: aws.String(BLOGS_TABLE),
    Item: map[string]*dynamodb.AttributeValue{
      "slug": {
        S: aws.String(blog.Slug),
      },
      "title": {
        S: aws.String(blog.Title),
      },
      "description": {
        S: aws.String(blog.Description),
      },
      "content": {
        S: aws.String(blog.Content),
      },
      "created_at": {
        S: aws.String(blog.CreatedAt),
      },
    },
  }

  _, err := u.databaseStore.PutItem(item)
  if err != nil {
    return err
  }

  return nil
}

func (u DynamoUserStore) DoesUserExist(username string) (bool, error) {
  // aws force to pass in reference here, also passing reference is faster than passing copy
  result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
    // checking if there's a record in the dynamodb table where key is username and value is what we pass in
    TableName: aws.String(USERS_TABLE),
    Key: map[string]*dynamodb.AttributeValue{
      "username": {
        // S here means string for aws, similar things applied to boolean, int
        S: aws.String(username),
      },

    },
  })

  // if there is error
  if err != nil {
    return true, err
  }

  // if the user does no exist
  if result.Item == nil {
    return false, nil
  }

  return true, nil
}

func (u DynamoUserStore) InsertUser(user types.User) error {
  // assemble the type that dynamodb understand first
  item := &dynamodb.PutItemInput{
    TableName: aws.String(USERS_TABLE),
    Item: map[string]*dynamodb.AttributeValue{
      "username": {
        S: aws.String(user.Username),
      },
      "password": {
        // this at this point is plain text password
        S: aws.String(user.PasswordHash),
      },
    },
  }

  _, err := u.databaseStore.PutItem(item)
  if err != nil {
    return err
  }

  return nil
}

func (u DynamoUserStore) GetUser(username string) (types.User, error) {
  var user types.User
  result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
    TableName: aws.String(USERS_TABLE),
    Key: map[string]*dynamodb.AttributeValue {
      "username": {
        S: aws.String(username),
      },
    },
  })

  if err != nil {
    return user, err
  }

  if result.Item == nil {
    return user, fmt.Errorf("user not found")
  }

  // map result to user struct
  err = dynamodbattribute.UnmarshalMap(result.Item, &user)
  if err != nil {
    return user, err
  }

  return user, nil
}
