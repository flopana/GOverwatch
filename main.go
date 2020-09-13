package main

import (
	"bufio"
	"fmt"
	"github.com/antchfx/jsonquery"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"os"
	"strconv"
)

var owStartRound int
//https://steamapi.xpaw.me/#ISteamUser/GetPlayerSummaries
const SteamWebApiKey = "5FF5EF4778DF00E92FB0B76676DFE822"

func main() {
	welcome := `   __________                                 __       __  
  / ____/ __ \_   _____  ______      ______ _/ /______/ /_ 
 / / __/ / / / | / / _ \/ ___/ | /| / / __  / __/ ___/ __ \
/ /_/ / /_/ /| |/ /  __/ /   | |/ |/ / /_/ / /_/ /__/ / / /
\____/\____/ |___/\___/_/    |__/|__/\__,_/\__/\___/_/ /_/
		`
	fmt.Println(welcome)

	f, err := os.Open("003435053515502780722_0826630968.dem")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()
	//Register handler on kill events
	p.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		if p.GameState().TotalRoundsPlayed()+1 >= owStartRound{
			allplayers := p.GameState().Participants().Playing()
			fmt.Println("##########################################################################")
			fmt.Printf("Current Round: %d\n\n", p.GameState().TotalRoundsPlayed()+1)
			for _, player := range allplayers {
				var profileurl string
				if player.SteamID64 != 0 {
					doc, _ := jsonquery.LoadURL("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=" + SteamWebApiKey + "&steamids=" + strconv.FormatUint(player.SteamID64, 10))

					for _, n := range jsonquery.Find(doc, "response/players/*/profileurl") {
						profileurl = n.InnerText()
					}
				}
				var team string
				if player.Team == 2 {
					team = "T"
				} else {
					team = "CT"
				}
				fmt.Printf("Team: %s ,Player: %s, SteamID64: %d, Profile: %s\n", team, player.Name, player.SteamID64, profileurl)
				fmt.Printf("K: %d, A: %d, D: %d\n\n", player.Kills(), player.Assists(), player.Deaths())
			}
				fmt.Print("Advance to next round?")
				reader := bufio.NewReader(os.Stdin)
				_, _ = reader.ReadString('\n')
		}
	})
	p.RegisterEventHandler(func(e events.MatchStart) {
		fmt.Printf("Map: %s\n", p.Header().MapName)
		fmt.Printf("Server: %s\n\n", p.Header().ServerName)

		fmt.Println("In which round did your Overwatch case start?")
		_, _ = fmt.Scanf("%d", &owStartRound)
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		fmt.Println(err)
	}
}