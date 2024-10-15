package sdkgocql

import (
	"flag"
	"strings"

	"github.com/gocql/gocql"
)

type GcqlOpt struct {
	prefix   string
	strHosts string
	keyspace string
	username string
	password string
}

type gcqlDB struct {
	name string
	*GcqlOpt
}

func NewGcqlDB(name, prefix string) *gcqlDB {
	return &gcqlDB{
		name: name,
		GcqlOpt: &GcqlOpt{
			prefix: prefix,
		},
	}
}

func (gdb *gcqlDB) GetPrefix() string {
	return gdb.name
}

func (gdb *gcqlDB) Name() string {
	return gdb.name
}

func (gdb *gcqlDB) InitFlags() {
	flag.StringVar(&gdb.strHosts, gdb.prefix+"-hosts", "127.0.0.1,127.0.0.2,127.0.0.3", "Cassandra database connection-string")
	flag.StringVar(&gdb.keyspace, gdb.prefix+"-db-keyspace", "", "Cassandra database keyspace")
	flag.StringVar(&gdb.username, gdb.prefix+"-db-username", "", "Cassandra database username")
	flag.StringVar(&gdb.password, gdb.prefix+"-db-password", "", "Cassandra database password")
}

func (gdb *gcqlDB) Configure() error {
	return nil
}

func (gdb *gcqlDB) Run() error {
	return nil
}

func (gdb *gcqlDB) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

func (gdb *gcqlDB) Get() interface{} {
	cluster := gocql.NewCluster(strings.Split(gdb.strHosts, ",")...)
	cluster.Keyspace = gdb.keyspace
	cluster.ProtoVersion = 4

	if gdb.username != "" && gdb.password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: gdb.username,
			Password: gdb.password,
		}
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	return session
}
