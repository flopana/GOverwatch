package main

import (
	"bufio"
	"fmt"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"os"
)

func main() {
	welcome := `   __________                                 __       __  
  / ____/ __ \_   _____  ______      ______ _/ /______/ /_ 
 / / __/ / / / | / / _ \/ ___/ | /| / / __  / __/ ___/ __ \
/ /_/ / /_/ /| |/ /  __/ /   | |/ |/ / /_/ / /_/ /__/ / / /
\____/\____/ |___/\___/_/    |__/|__/\__,_/\__/\___/_/ /_/
		`
	fmt.Println(welcome)

	var owStartRound int
	fmt.Println("In which round did your Overwatch case start?")
	_, err := fmt.Scanf("%d", &owStartRound)

	f, err := os.Open("003435053515502780722_0826630968.dem")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()
	//Register handler on kill events
	p.RegisterEventHandler(func(e events.RoundEnd) {
		allplayers := p.GameState().Participants().Playing()
		fmt.Println("##########################################################################")
		if p.GameState().IsWarmupPeriod(){
			fmt.Printf("Current Round: Warmup\n\n")
		}else{
			fmt.Printf("Current Round: %d\n\n", p.GameState().TotalRoundsPlayed()-1)
		}
		for _, player := range allplayers{
			fmt.Printf("Player: %s, SteamID64: %d\n", player.Name, player.SteamID64)
			fmt.Printf("K: %d, A: %d, D: %d\n\n", player.Kills(), player.Assists(), player.Deaths())
		}
		if p.GameState().TotalRoundsPlayed()-1 >= owStartRound {
			fmt.Print("Advance to next round?")
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadString('\n')
		}
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		fmt.Println(err)
	}
}