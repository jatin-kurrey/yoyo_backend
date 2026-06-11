package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func NewBookingCode() string {
	return fmt.Sprintf("YOYO-%s-%06d", time.Now().Format("20060102"), rand.Intn(900000)+100000)
}
