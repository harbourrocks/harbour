package context

import (
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/sirupsen/logrus"
)

type HRock struct {
	L        *logrus.Entry
	CtxIdent uuid.UUID
	IdToken  *auth.IdToken
}
