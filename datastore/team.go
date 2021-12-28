package datastore

import (
	"context"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Team struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	City             string `json:"city" bson:"city"`
	Season           string `json:"season" bson:"season"`
	League           string `json:"league" bson:"league"`
 }
 
 func NewTeam(name string) (*Team, error) {
	team := &Team{
		Name: name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mgm.Coll(team).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"Name": 1},
		Options: options.Index().SetUnique(true),	
	})
	err := mgm.Coll(team).Create(team)
	return team, err
}

func ListTeams() (*[]Team, error) {
	result := []Team{}
	err := mgm.Coll(&Team{}).SimpleFind(&result, bson.M{})
	return &result, err
}

func GetTeam(teamId string) (*Team, error) {
	team := &Team{}
	err := mgm.Coll(team).FindByID(teamId, team)
	return team, err
}
