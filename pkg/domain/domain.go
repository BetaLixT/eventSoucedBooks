package domain

import (
	"eventSourcedBooks/pkg/domain/auth"
	"eventSourcedBooks/pkg/domain/courtroom"

	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	auth.NewAuthService,
	courtroom.NewCourtroomService,
)
