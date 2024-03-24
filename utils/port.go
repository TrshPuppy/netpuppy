package utils

import (
	"math/rand"
	"time"
)

func randomPort() int {
	rand.Seed(time.Now().UnixNano())
	const minPort = 1024
	const maxPort = 65535
	return rand.Intn(maxPort-minPort+1) + minPort
}
