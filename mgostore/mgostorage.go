package mgostore

import (
	"github.com/RangelReale/osin"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// collection names for the entities
const (
	CLIENT_COL    = "clients"
	AUTHORIZE_COL = "authorizations"
	ACCESS_COL    = "accesses"
)

const REFRESHTOKEN = "refreshtoken"

type MongoStorage struct {
	database *mgo.Database
}

func NewMongoStorage(database *mgo.Database) *MongoStorage {
	storage := &MongoStorage{database}
	index := mgo.Index{
		Key:        []string{REFRESHTOKEN},
		Unique:     false, // refreshtoken is sometimes empty
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	accesses := storage.database.C(ACCESS_COL)
	err := accesses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return storage
}

func (oauth *MongoStorage) GetClient(id string) (*osin.Client, error) {
	clients := oauth.database.C(CLIENT_COL)
	client := new(osin.Client)
	err := clients.FindId(id).One(client)
	return client, err
}

func (oauth *MongoStorage) SetClient(id string, client *osin.Client) error {
	clients := oauth.database.C(CLIENT_COL)
	_, err := clients.UpsertId(id, client)
	return err
}

func (oauth *MongoStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	authorizations := oauth.database.C(AUTHORIZE_COL)
	_, err := authorizations.UpsertId(data.Code, data)
	return err
}

func (oauth *MongoStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	authorizations := oauth.database.C(AUTHORIZE_COL)
	authData := new(osin.AuthorizeData)
	err := authorizations.FindId(code).One(authData)
	return authData, err
}

func (oauth *MongoStorage) RemoveAuthorize(code string) error {
	authorizations := oauth.database.C(AUTHORIZE_COL)
	return authorizations.RemoveId(code)
}

func (oauth *MongoStorage) SaveAccess(data *osin.AccessData) error {
	accesses := oauth.database.C(ACCESS_COL)
	_, err := accesses.UpsertId(data.AccessToken, data)
	return err
}

func (oauth *MongoStorage) LoadAccess(token string) (*osin.AccessData, error) {
	accesses := oauth.database.C(ACCESS_COL)
	accData := new(osin.AccessData)
	err := accesses.FindId(token).One(accData)
	return accData, err
}

func (oauth *MongoStorage) RemoveAccess(token string) error {
	accesses := oauth.database.C(ACCESS_COL)
	return accesses.RemoveId(token)
}

func (oauth *MongoStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	accesses := oauth.database.C(ACCESS_COL)
	accData := new(osin.AccessData)
	err := accesses.Find(bson.M{REFRESHTOKEN: token}).One(accData)
	return accData, err
}

func (oauth *MongoStorage) RemoveRefresh(token string) error {
	accesses := oauth.database.C(ACCESS_COL)
	return accesses.Update(bson.M{REFRESHTOKEN: token}, bson.M{
		"$unset": bson.M{
			REFRESHTOKEN: 1,
		}})
}
