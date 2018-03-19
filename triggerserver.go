package main

import (
	"github.com/RATDistributedSystems/utilities"
	"github.com/RATDistributedSystems/utilities/ratdatabase"
	"github.com/gocql/gocql"
)

var sessionGlobalTS *gocql.Session
var sessionGlobalTR *gocql.Session
var transactionNumGlobal = 0
var configurationServer = utilities.GetConfigurationFile("config.json")
var auditPool = initializePool(150, 190, "audit")
var serverName = "trigger"


func main() {
	initCassandraTS()
	initCassandraTR()
	initCheckTriggers()
}


func initCassandraTS() {
	//connect to database for transaction server databases
	hostname := configurationServer.GetValue("cassandra_ip_ts")
	keyspace := configurationServer.GetValue("cassandra_keyspace_ts")
	protocol := configurationServer.GetValue("cassandra_proto_ts")
	ratdatabase.InitCassandraConnection(hostname, keyspace, protocol)
	sessionGlobalTS = ratdatabase.CassandraConnection
}

func initCassandraTR() {
	//connect to database for trigger server databases
	hostname := configurationServer.GetValue("cassandra_ip_tr")
	keyspace := configurationServer.GetValue("cassandra_keyspace_tr")
	protocol := configurationServer.GetValue("cassandra_proto_tr")
	ratdatabase.InitCassandraConnection(hostname, keyspace, protocol)
	sessionGlobalTR = ratdatabase.CassandraConnection
}


func initCheckTriggers(){
	checkBuyTriggers()
	checkSellTriggers()
}

