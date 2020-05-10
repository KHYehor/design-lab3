package main

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func TestBalancer(t *testing.T) {
	assert := assert.New(t)
	serversPoolTraffic[0] = 100
	serversPoolTraffic[1] = 50
	serversPoolTraffic[2] = 200

	var bestServer1, _ = getTheBestServer()
	// Getting the minimalServer
	assert.Equal(serversPool[1], bestServer1)
	serverStatus[0] = false
	serversPoolTraffic[0] = 0
	serversPoolTraffic[1] = 200
	serverStatus[2] = false
	serversPoolTraffic[2] = 0

	var bestServer2, _ = getTheBestServer()
	// Getting the server without Error
	assert.Equal(serversPool[1], bestServer2)
	serverStatus[0] = false
	serverStatus[1] = false
	serverStatus[2] = false
	serversPoolTraffic[0] = -1
	serversPoolTraffic[1] = -1
	serversPoolTraffic[2] = -1

	var _, err3 = getTheBestServer()
	// Getting the error if error is everywhere
	assert.Equal("every server is not healthy", err3.Error())
}
