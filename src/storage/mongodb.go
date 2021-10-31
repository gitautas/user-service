package storage

import (
	"context"
	"errors"
	"log"
	"time"
	"user-service/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database interface {
	Connect() error
	Disconnect() error
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error
	GetUser(userID string) (file *models.User, err error)
	GetUserList(limit int, skip int, filter map[string]string) (users []*models.User, err error)
}

type Mongo struct {
	Endpoint       string
	DBName         string
	CollectionName string
	QueryTimeout   time.Duration
	client         *mongo.Client
	collection     *mongo.Collection
}

func (m *Mongo) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Endpoint))
	if err != nil {
		log.Printf("Could not connect to MongoDB instance: %v", err)
		return err
	}

	m.client = client
	err = m.client.Ping(ctx, readpref.Primary()) // Establish and check connection
	if err != nil {
		log.Printf("Could not verify connection to MongoDB instance: %v", err)
		return err
	}

	m.collection = client.Database(m.DBName).Collection(m.CollectionName)
	return nil
}

func (m *Mongo) Disconnect() error {
	if m.client != nil {
		log.Println("Disconnecting from MongoDB")
		return m.client.Disconnect(context.TODO())
	}
	return nil
}

func (m *Mongo) CreateUser(user *models.User) error {
	user.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	user.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	bUser, err := bson.Marshal(user)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeout)
	defer cancel()
	_, err = m.collection.InsertOne(ctx, bUser)
	if err != nil {
		log.Printf("Could not create document: %v", err)
		return err
	}
	return nil
}

func (m *Mongo) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	bUser, err := bson.Marshal(user)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeout)
	defer cancel()
	_, err = m.collection.UpdateOne(ctx, bson.M{"id": user.Id}, bUser)
	if err != nil {
		log.Printf("Could not update document: %v", err)
		return err
	}
	return nil
}

func (m *Mongo) DeleteUser(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.collection.DeleteOne(ctx, bson.M{"id": userID})
	if err != nil {
		log.Printf("Could not delete document: %v", err)
		return err
	}
	return nil
}

func (m *Mongo) GetUser(userID string) (user *models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := m.collection.FindOne(ctx, bson.M{"id": userID})
	err = result.Decode(&user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (m *Mongo) GetUserList(limit int, skip int, filter map[string]string) (users []*models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(skip))

	bFilter := bson.M{}
	for key, value := range filter {
		bFilter[key] = value
	}

	cur, err := m.collection.Find(ctx, bFilter, opts)
	if err != nil {
		log.Printf("Could not find documents: %v", err)
		return nil, err
	}

	for cur.Next(ctx) {
		var user *models.User
		err = cur.Decode(&user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}
