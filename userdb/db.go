package userdb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"grpc-gateway-example/config"
	"grpc-gateway-example/model"
	"net"
	"net/http"
	"time"
)

type DB struct {
	*dynamo.DB
	User dynamo.Table
}

func New(dbConfig config.DatabaseConfig) (*DB, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
		TLSHandshakeTimeout:   3 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
	}

	if err := http2.ConfigureTransport(transport); err != nil {
		log.Err(err)
	}

	awsConfig := &aws.Config{
		Region:   aws.String(dbConfig.Region),
		Endpoint: aws.String(dbConfig.Endpoint),
		HTTPClient: &http.Client{
			Transport: transport,
		},
	}

	newSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	database := dynamo.New(newSession)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:   database,
		User: database.Table(dbConfig.User),
	}, nil
}

func (db *DB) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := db.PutUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) PutUser(ctx context.Context, user *model.User) error {
	if err := db.User.Put(user).RunWithContext(ctx); err != nil {
		return errors.Wrapf(err, "PutUser() failed, user_id:%s, nickname: %s, name: %s", user.ID, user.Nickname, user.Name)
	}
	return nil
}

func (db *DB) GetUserByID(userId string) (*model.User, error) {
	var user model.User
	if err := db.User.Get("user_id", userId).One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
