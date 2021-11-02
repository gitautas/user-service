package storage

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"user-service/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database interface {
	Connect() *models.Status
	Disconnect() *models.Status
	CreateUser(user *models.User) *models.Status
	UpdateUser(user *models.User) *models.Status
	DeleteUser(userID string) *models.Status
	GetUser(userID string) (user *models.User, err *models.Status)
	GetUserList(limit int, skip int, filter map[string]string) (users []*models.User, err *models.Status)
}

type Mongo struct {
	endpoint       string
	dbName         string
	collectionName string
	queryTimeout   time.Duration
	client         *mongo.Client
	collection     *mongo.Collection
}

func NewMongo(endpoint string, dbname string, collection string, timeout time.Duration) *Mongo {
	return &Mongo{
		endpoint:       endpoint,
		dbName:         dbname,
		collectionName: collection,
		queryTimeout:   timeout,
	}
}

func (m *Mongo) Connect() *models.Status {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.endpoint))
	if err != nil {
		return models.NewStatus(http.StatusInternalServerError, err.Error())
	}

	m.client = client
	fmt.Println("Pinging MongoDB")
	err = m.client.Ping(ctx, readpref.Primary()) // Establish and check connection
	if err != nil {
		return models.NewStatus(http.StatusInternalServerError, err.Error())
	}

	m.collection = client.Database(m.dbName).Collection(m.collectionName)
	return nil
}

func (m *Mongo) Disconnect() *models.Status {
	if m.client != nil {
		err := m.client.Disconnect(context.TODO())
		if err != nil {
			return models.NewStatus(http.StatusInternalServerError, "could not disconnect from database")
		}
	}
	return nil
}

func (m *Mongo) CreateUser(user *models.User) *models.Status {
	bUser, err := bson.Marshal(user)
	if err != nil {
		return models.NewStatus(http.StatusBadRequest, "could not parse user")
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.queryTimeout)
	defer cancel()
	_, err = m.collection.InsertOne(ctx, bUser)
	if err != nil {
		return models.NewStatus(http.StatusInternalServerError, "could not create user")
	}
	return nil
}

func (m *Mongo) UpdateUser(user *models.User) *models.Status {
	ctx, cancel := context.WithTimeout(context.Background(), m.queryTimeout)
	defer cancel()
	_, err := m.collection.UpdateOne(ctx, bson.M{"id": user.Id}, bson.M{"$set": user})
	if err != nil {
		return models.NewStatus(http.StatusInternalServerError, "could not update user")
	}
	return nil
}

func (m *Mongo) DeleteUser(userID string) *models.Status {
	ctx, cancel := context.WithTimeout(context.Background(), m.queryTimeout)
	defer cancel()
	_, err := m.collection.DeleteOne(ctx, bson.M{"id": userID})
	if err != nil {
		return models.NewStatus(http.StatusInternalServerError, "could not delete user")
	}
	return nil
}

func (m *Mongo) GetUser(userID string) (user *models.User, status *models.Status) {
	ctx, cancel := context.WithTimeout(context.Background(), m.queryTimeout)
	defer cancel()
	result := m.collection.FindOne(ctx, bson.M{"id": userID})
	err := result.Decode(&user)
	if err != nil {
		return nil, models.NewStatus(http.StatusNotFound, "user not found")
	}

	return user, nil
}

func (m *Mongo) GetUserList(limit int, skip int, filter map[string]string) (users []*models.User, status *models.Status) {
	ctx, cancel := context.WithTimeout(context.Background(), m.queryTimeout)
	defer cancel()

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(skip))

	bFilter := bson.M{}
	for key, value := range filter {
		bFilter[key] = value
	}

	fmt.Println(bFilter)

	cur, err := m.collection.Find(ctx, bFilter, opts)
	if err != nil {
		return nil, models.NewStatus(http.StatusInternalServerError, "could not query database")
	}

	for cur.Next(ctx) {
		var user *models.User
		err = cur.Decode(&user)
		if err != nil {
			return nil, models.NewStatus(http.StatusInternalServerError, "database error")
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return []*models.User{}, nil
	}

	return users, nil
}
