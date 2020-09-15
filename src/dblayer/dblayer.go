package dblayer

import (
	"github.com/ajtwoddltka/GCS/src/model"
)

type DBLayer interface {
	GetScore(model.User)
	GetRanking(model.User)
}
