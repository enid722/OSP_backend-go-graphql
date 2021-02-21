package database

import (
	"context"
	"log"
	"time"

	"github.com/enid722/OSP_backend-go-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB - DB structure
type DB struct {
	client *mongo.Client
}

//Connect - Connect mongoDB
func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	return &DB{
		client: client,
	}
}

//Save - Save new survey
func (db *DB) Save(input *model.SurveyInput) *model.Survey {
	collection := db.client.Database("OSP").Collection("surveys")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, input)

	if err != nil {
		log.Fatal(err)
	}

	return &model.Survey{
		ID:        res.InsertedID.(primitive.ObjectID).Hex(),
		Title:     input.Title,
		Token:     input.Token,
		IsDeleted: false,
	}
}

//FindByID - Find survey by ID
func (db *DB) FindByID(ID string) *model.Survey {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Fatal(err)
	}
	collection := db.client.Database("OSP").Collection("surveys")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"_id": ObjectID})
	survey := model.Survey{}
	res.Decode(&survey)

	return &survey
}

//All - Find all surveys
func (db *DB) All() []*model.Survey {
	collection := db.client.Database("OSP").Collection("surveys")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var surveys []*model.Survey
	for cur.Next(ctx) {
		var survey *model.Survey
		err := cur.Decode(&survey)
		if err != nil {
			log.Fatal(err)
		}
		surveys = append(surveys, survey)
	}
	return surveys
}
