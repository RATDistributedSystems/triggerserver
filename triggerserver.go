package main

import (
	"fmt"
	"github.com/RATDistributedSystems/utilities"
	"github.com/RATDistributedSystems/utilities/ratdatabase"
	"github.com/gocql/gocql"
)

var sessionGlobalTS *gocql.Session
var sessionGlobalTR *gocql.Session
var transactionNumGlobal = 0
var configurationServer = utilities.Load()
var auditPool = initializePool(5, 80, "audit")
var serverName = "trigger"

func main() {
	fmt.Println("Connecting to TS database")
	initCassandraTS()
	fmt.Println("Connecting to TR database")
	initCassandraTR()
	fmt.Println("Checking Triggers")
	initCheckTriggers()
}

func initCassandraTS() {
	//connect to database for transaction server databases
	hostname := configurationServer.GetValue("transdb_ip")
	keyspace := configurationServer.GetValue("transdb_keyspace")
	protocol := configurationServer.GetValue("transdb_proto")
	ratdatabase.InitCassandraConnection(hostname, keyspace, protocol)
	sessionGlobalTS = ratdatabase.CassandraConnection
}

func initCassandraTR() {
	//connect to database for trigger server databases
	hostname := configurationServer.GetValue("triggerdb_ip")
	keyspace := configurationServer.GetValue("triggerdb_keyspace")
	protocol := configurationServer.GetValue("triggerdb_proto")
	ratdatabase.InitCassandraConnection(hostname, keyspace, protocol)
	sessionGlobalTR = ratdatabase.CassandraConnection
}

func initCheckTriggers() {
	fmt.Println("Starting Check Triggers")
	go checkBuyTriggers()
	//go checkSellTriggers()

	for{

	}
}
