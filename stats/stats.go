package stats

import (
	"fmt"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
)

type Stats struct {
	player	common.Player
	Kills	int
	Assists	int
	Deaths	int
}

func getStats(s Stats) string{
	return fmt.Sprintf("K: %d, A: %d, D: %d", s.Kills, s.Assists, s.Deaths)
}
