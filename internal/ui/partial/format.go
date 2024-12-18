package partial

import (
	"fmt"
	"time"
)

func raw(x int) string {
	return fmt.Sprintf("%.2f%", x)
}

func per(x float32) string {
	return fmt.Sprintf("%.2f%", x)
}

func winloss(wins, losses int) string {
	return fmt.Sprintf("%d-%d", wins, losses)
}

func csper10() string {
	return ""
}

func rank() string {
	return ""
}

func kda(kills, deaths, assists int) string {
	return ""
}

func riotid(id, tag string) string {
	return fmt.Sprintf("%s#%s", id, tag)
}

func date(t time.Time) string {
	return ""
}

func duration(d time.Duration) string {
	return ""
}
