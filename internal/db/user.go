package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user document.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Roles    []string           `bson:"roles"`
}

// ValidatePassword validates the given plaintext password for the user
func (u *User) ValidatePassword(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	if err != nil {
		return err
	}
	return nil
}

// FindUserByName searches for a user with the given username.
func (MGO) FindUserByName(username string) (*User, error) {
	data := &User{}
	filter := bson.M{
		"username": username,
	}
	res := col.FindOne(context.Background(), filter)
	err := res.Decode(data)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("unable to find user with username: %v", username)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

// FindUserByID searches for a user with the given id.
func (MGO) FindUserByID(id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("cannot parse id: %v", id)
	}
	data := &User{}
	filter := bson.M{"_id": oid}
	res := col.FindOne(context.Background(), filter)
	err = res.Decode(data)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("unable to find user with id: %v", id)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}
