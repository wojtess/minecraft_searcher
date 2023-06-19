package scanner

import (
	"fmt"
	"time"

	"github.com/Tnze/go-mc/bot"
)

func pingServer(ip string, timeout time.Duration) (string, time.Duration, error) {
	resp, delay, err := bot.PingAndListTimeout(fmt.Sprintf("%s:%d", ip, 25565), timeout)
	if err != nil {
		return "", 0, fmt.Errorf("ping and list server fail: %s", err)
	}
	return string(resp), delay, nil
}
