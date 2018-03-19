package main

import (
	"log"
	"net"
	"sync"
)

type connectionPool struct {
	serverName      string
	currentSize     int
	stepSize        int
	maxPoolSize     int
	freeConnections []net.Conn
	mux             sync.Mutex
}

func initializePool(poolSize int, maxPoolSize int, serverName string) *connectionPool {
	var connPool connectionPool
	connPool.maxPoolSize = maxPoolSize
	connPool.freeConnections = make([]net.Conn, 0)
	connPool.serverName = serverName
	connPool.currentSize = poolSize
	connPool.stepSize = 10
	connPool.addConnections(connPool.currentSize)
	return &connPool
}

func (c *connectionPool) getConnection() net.Conn {
	c.mux.Lock()

	// if none are free add more
	if len(c.freeConnections) < 1 {
		c.addConnections(c.stepSize)
	}

	conn := c.freeConnections[0]
	c.freeConnections = c.freeConnections[1:]
	c.mux.Unlock()
	return conn
}

func (c *connectionPool) returnConnection(conn net.Conn) {
	c.mux.Lock()
	c.freeConnections = append(c.freeConnections, conn)
	c.mux.Unlock()
}

func (c *connectionPool) addConnections(stepSize int) {
	if len(c.freeConnections) > c.maxPoolSize {
		return
	}
	for i := 0; i < stepSize; i++ {
		addr, protocol := configurationServer.GetServerDetails(c.serverName)
		conn, err := net.Dial(protocol, addr)
		if err != nil {
			log.Fatalf("Could not make another connection for %s server\n%s", c.serverName, err.Error())
		}
		c.freeConnections = append(c.freeConnections, conn)
	}
}
