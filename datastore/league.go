package datastore

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type League struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
}
 
func NewLeague(name string) (*League, error) {
	league := &League{
		Name: name,
	}

	mgm.Coll(league).Indexes().CreateOne(mgm.Ctx(), mongo.IndexModel{
		Keys: bson.M{"Name": 1},
		Options: options.Index().SetUnique(true),	
	})
	err := mgm.Coll(league).Create(league)
	return league, err
}

func ListLeagues() (*[]League, error) {
	result := []League{}
	err := mgm.Coll(&League{}).SimpleFind(&result, bson.M{})
	return &result, err
}

func GetLeague(leagueId string) (*League, error) {
	league := &League{}
	err := mgm.Coll(league).FindByID(leagueId, league)
	return league, err
}

func DeleteLeague(leagueId string) error {
	league, err := GetLeague(leagueId)
	if err != nil {
		return err
	}
	return mgm.Coll(league).Delete(league)
}
