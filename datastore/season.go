package datastore

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Season struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	League           string `json:"league" bson:"league"`
}

func NewSeason(name, leagueId string) (*Season, error) {
	season := &Season{
		Name: name,
		League: leagueId,
	}

	mgm.Coll(season).Indexes().CreateOne(mgm.Ctx(), mongo.IndexModel{
		Keys: bson.M{
			"Name":   1,
			"League": 1,
		},
		Options: options.Index().SetUnique(true),	
	})
	err := mgm.Coll(season).Create(season)
	return season, err	
}

func GetSeason(seasonId string) (*Season, error) {
	season := &Season{}
	err := mgm.Coll(season).FindByID(seasonId, season)
	return season, err
}

func DeleteSeason(seasonId string) error {
	season, err := GetSeason(seasonId)
	if err != nil {
		return err
	}
	return mgm.Coll(season).Delete(season)
}
