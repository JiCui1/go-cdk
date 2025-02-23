package types

import (
  "github.com/golang-jwt/jwt/v5"
  "golang.org/x/crypto/bcrypt"
  "time"
  "strings"
  "regexp"
)

type RegisterUser struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

type User struct {
  Username string `json:"username"`
  PasswordHash string `json:"password"`
}

type Blog struct {
  Slug string `json:"slug"`
  Title string `json:"title"`
  Description string `json:"description"`
  Content string `json:"content"`
  CreatedAt string  `json:"created_at"`
}

func Slugify(title string) string {
	slug := strings.ToLower(title)

	// Remove all non-alphanumeric characters except spaces
	reg, _ := regexp.Compile(`[^a-z0-9\s]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace multiple spaces with a single space
	regSpace, _ := regexp.Compile(`\s+`)
	slug = regSpace.ReplaceAllString(slug, " ")

	// Trim spaces from the beginning and end, then replace spaces with hyphens
	slug = strings.TrimSpace(slug)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

func NewUser(registerUser RegisterUser) (User, error) {
  //transforms from string to byes
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10) 

  if err != nil {
    return User{}, err
  }

  return User {
    Username: registerUser.Username,
    PasswordHash: string(hashedPassword),
  }, nil
}

func ValidatePassword(hashedPassword, plainTextPassword string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))

  return err == nil
}

func CreateToken(user User) string {
  now := time.Now()
  validUntil := now.Add(time.Hour * 1).Unix()

  claims := jwt.MapClaims{
    "user": user.Username,
    "expires": validUntil,
  }


  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
  secret := "jiahua-zheshiyigemimi-stuff"

  tokenString, err := token.SignedString([]byte(secret))
  if err != nil {
    return ""
  }

  return tokenString
}
