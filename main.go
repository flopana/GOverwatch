package main

import (
	"fmt"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
)

var allplayers map[int]*common.Player

func main() {
	welcome := `   __________                                 __       __  
  / ____/ __ \_   _____  ______      ______ _/ /______/ /_ 
 / / __/ / / / | / / _ \/ ___/ | /| / / __  / __/ ___/ __ \
/ /_/ / /_/ /| |/ /  __/ /   | |/ |/ / /_/ / /_/ /__/ / / /
\____/\____/ |___/\___/_/    |__/|__/\__,_/\__/\___/_/ /_/
		`
	fmt.Println(welcome)

	f, err := os.Open("003434868561326113214_1747621417.dem")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()
	//Register handler on kill events
	p.RegisterEventHandler(func(e events.Kill) {
		var hs string
		if e.IsHeadshot {
			hs = " (HS)"
		}
		var wallBang string
		if e.PenetratedObjects > 0 {
			wallBang = " (WB)"
		}
		fmt.Printf("%s <%v%s%s> %s\n", e.Killer, e.Weapon, hs, wallBang, e.Victim)
	})
	p.RegisterEventHandler(func(e events.MatchStart) {
		allplayers = p.GameState().Participants().AllByUserID()
		delete(allplayers, 2) // Delete GOTV entity
	})

	// Parse to end
	err = p.ParseToEnd()

	fmt.Println(allplayers)

	if err != nil {
		panic(err)
	}
}