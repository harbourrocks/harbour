package context

import (
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	l "github.com/sirupsen/logrus"
)

// HRockTrait returns HRock
type HRockTrait interface {
	traits.IdTokenTrait
	GetHRock() HRock
	SetHRock(HRock)
}

// HRockModel holds the request
type HRockModel struct {
	hRock HRock
	traits.IdTokenModel
}

func (m HRockModel) GetHRock() HRock {
	return m.hRock
}

func (m *HRockModel) SetHRock(s HRock) {
	m.hRock = s
}

func AddHRock(trait HRockTrait) {
	rqId := uuid.New()

	trait.SetHRock(HRock{
		L:          l.WithField("rqId", rqId.String()),
		CtxIdent:   rqId,
		IdToken:    trait.GetToken(),
		IdTokenStr: trait.GetTokenStr(),
	})
}
