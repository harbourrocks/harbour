package context

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type HRock struct {
	L        *logrus.Entry
	CtxIdent uuid.UUID
}
