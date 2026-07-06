package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateTicketCode() string {

	rand.Seed(time.Now().UnixNano())

	return fmt.Sprintf(
		"ADU-%d-%04d",
		time.Now().Unix(),
		rand.Intn(10000),
	)
}