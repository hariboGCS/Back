package dblayer

import (
	"github.com/ajtwoddltka/GCS/src/model"
)

type DBLayer interface {
	GetScore(model.Student) model.Score
	GetRanking(model.Student) model.Score
}
