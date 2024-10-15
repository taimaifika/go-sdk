package common

import "fmt"

// Plugin name
const (
	PluginEnvConf    = "envconf"
	PluginSimple     = "simple"
	PluginDbMysql    = "mysql"
	PluginDbPostgres = "pg"
	PluginPubSub     = "pubsub"
	PluginItemAPI    = "item-api"
)

// Topic name
const (
	TopicUserLikedItem   = "TopicUserLikedItem"
	TopicUserUnlikedItem = "TopicUserUnlikedIt"
)

// Recovery is a helper function to recover from panic
func Recovery() {
	if r := recover(); r != nil {
		fmt.Println("Recovered:", r)
	}
}

// Database
type DbType int

const (
	DbTypeItem DbType = 1
	DbTypeUser DbType = 2
)

const (
	CurrentUser = "current_user"
)

type Requester interface {
	GetUserID() string
	GetDeviceUUID() string
	GetRole() string
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}
