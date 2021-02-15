package helper

import (
	"context"
	"eltropy/model"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/twinj/uuid"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoOnce   sync.Once
	ctx         context.Context
	clientErr   error
	client      *mongo.Client
	redisClient *redis.Client
)

const (
	// URI - uri for mongo db
	URI = "mongodb://localhost:27017"
	// DB - database name
	DB = "eltropy_db"
	// AdminCollection - admin collection
	AdminCollection = "admin"
	// EmployeeCollection - employee collection
	EmployeeCollection = "employee"
)

// TokenDetails - detail for token
type TokenDetails struct {
	Token  string
	UUID   string
	Expiry int64
}

type AccessDetails struct {
	UUID     string
	Username string
}

// GetMongoClient - get mongodb client
func GetMongoClient() (*mongo.Client, error) {

	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(URI)
		client, clientErr = mongo.Connect(context.TODO(), clientOptions)
	})
	return client, clientErr
}

// InistalizeRedis - initialize redis
func InistalizeRedis() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}

// CreateToken - create token using jwt
func CreateToken(username string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.Expiry = time.Now().Add(time.Minute * 10).Unix()
	td.UUID = uuid.NewV4().String()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   username,
		"exp":        time.Now().Add(time.Minute * 10).Unix(),
		"uuid":       td.UUID,
		"authorized": true,
	})
	var err error
	td.Token, err = at.SignedString([]byte("eltropy_secret_key"))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// ExtractToken - extract token
func ExtractToken(r *http.Request) (*AccessDetails, error) {

	bearerToken := r.Header.Get("Authorization")
	token, err := ValidateToken(bearerToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid, ok := claims["uuid"].(string)
		if !ok {
			return nil, err
		}
		username, ok := claims["username"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			UUID:     uuid,
			Username: username,
		}, nil
	}
	return nil, fmt.Errorf("Invalid token")
}

// CreateAuth - Create auth in redis
func CreateAuth(username string, td *TokenDetails) error {
	//converting Unix to UTC(to Time object)
	at := time.Unix(td.Expiry, 0)
	errAccess := redisClient.Set(td.UUID, username, at.Sub(time.Now())).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}

// DeleteAuth - delete auth in redis
func DeleteAuth(uuid string) (int64, error) {
	deleted, err := redisClient.Del(uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// ValidateToken - check whether token is valid or not
func ValidateToken(bearerToken string) (*jwt.Token, error) {
	if bearerToken == "" {
		return nil, fmt.Errorf("Token required")
	}
	token, _ := jwt.Parse(strings.Split(bearerToken, " ")[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("eltropy_secret_key"), nil
	})
	return token, nil
}

// FetchAuth - fetch auth from redis
func FetchAuth(ad *AccessDetails) (string, error) {
	username, err := redisClient.Get(ad.UUID).Result()
	if err != nil {
		return "", err
	}
	return username, nil
}

// Signout - Signout implementation for all user
func Signout(w http.ResponseWriter, r *http.Request) {
	accessDetails, err := ExtractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	deleted, err := DeleteAuth(accessDetails.UUID)
	if err != nil || deleted == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Code: http.StatusOK,
		Msg:  "Successfully logged out",
	})
}
