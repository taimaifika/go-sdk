package middleware

import (
	"context"

	"github.com/taimaifika/go-sdk/logger"
	"github.com/taimaifika/go-sdk/sdkcm"
)

type ServiceContext interface {
	Logger(prefix string) logger.Logger
	Get(prefix string) (interface{}, bool)
	MustGet(prefix string) interface{}
}

type CurrentUserProvider interface {
	GetCurrentUser(ctx context.Context, oauthID string) (sdkcm.User, error)
	ServiceContext
}

type Tracker interface {
	TrackApiCall(userId uint32, url string) error
}

type Caching interface {
	GetCurrentUser(ctx context.Context, sig string) (sdkcm.Requester, error)
	WriteCurrentUser(ctx context.Context, sig string, u sdkcm.Requester) error
}
