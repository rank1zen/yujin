package format

import "fmt"

func KDA(kills, deaths, assists int) string {
	return fmt.Sprintf("%d / %d / %d", kills, deaths, assists)
}

func Percentage(x float32) string {
	return fmt.Sprintf("%.2f%", x)
}

func Integer(x int) string {
	return fmt.Sprintf("%.2f%", x)
}

func RiotID(id, tag string) string {
	return fmt.Sprintf("%s#%s", id, tag)
}

func WinsLosses(wins, losses int) string {
	return fmt.Sprintf("%d-%d", wins, losses)
}

func CsPer10() string {
	return ""
}

func Rank() string {
	return ""
}
