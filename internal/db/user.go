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
	Mail     string             `bson:"mail"`
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

// FindUserByUsernameOrMail searches for a user which has the given string as username or mail address.
func (MGO) FindUserByUsernameOrMail(s string) (*User, error) {
	data := &User{}
	filter := bson.M{
		"$or": []bson.M{
			{"username": s},
			{"mail": s},
		},
	}
	res := col.FindOne(context.Background(), filter)
	err := res.Decode(data)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("unable to find user '%s'", s)
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
		return nil, fmt.Errorf("unable to find user with id: %s", id)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}
