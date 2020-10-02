package main

import (
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"github.com/mholt/archiver/v3"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var owStartRound int
var owStartRoundSet bool

//WarningColor ANSI Color yellow
const WarningColor = "\033[1;33m%s\033[0m"

//ErrorColor ANSI Color red
const ErrorColor = "\033[1;31m%s\033[0m"

var (
	snapshotLen int32 = 1024
	promiscuous bool  = false
	err         error
	timeout     time.Duration = 10 * time.Second
	handle      *pcap.Handle
)

func main() {
	welcome := `   __________                                 __       __
  / ____/ __ \_   _____  ______      ______ _/ /______/ /_
 / / __/ / / / | / / _ \/ ___/ | /| / / __  / __/ ___/ __ \
/ /_/ / /_/ /| |/ /  __/ /   | |/ |/ / /_/ / /_/ /__/ / / /
\____/\____/ |___/\___/_/    |__/|__/\____/\__/\___/_/ /_/
                `
	fmt.Println(welcome)

	//https://steamcommunity.com/dev/apikey
	config, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	doc, err := jsonquery.Parse(config)
	if err != nil {
		panic(err)
	}
	steamWebApiKey := jsonquery.FindOne(doc, "steamWebApiKey").InnerText()
	networkDevice := jsonquery.FindOne(doc, "networkDevice").InnerText()
	if steamWebApiKey == "" {
		fmt.Printf(WarningColor, "WARNING Your SteamWebApiKey is empty consider configuring this in the config.json,"+
			"\notherwise you will not get the Profile links"+
			"\nGet your API Key here https://steamcommunity.com/dev/apikey\n\n")
	}
	if networkDevice == "" {
		fmt.Printf(WarningColor, "The ethernet device in the config.json is empty choose one below\nPick a device Name and put it in the networkDevice in config.json\n")
		devices, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatal(err)
		}

		// Print networkDevice information
		fmt.Println("Devices found:")
		for _, ethernetDevice := range devices {
			fmt.Println("\nName: ", ethernetDevice.Name)
			fmt.Println("Description: ", ethernetDevice.Description)
			fmt.Println("Devices addresses: ", ethernetDevice.Description)
			for _, address := range ethernetDevice.Addresses {
				fmt.Println("- IP address: ", address.IP)
				fmt.Println("- Subnet mask: ", address.Netmask)
			}
		}
		os.Exit(0)
	}
	defer config.Close()

	fmt.Println("Searching for the demo\nPlease start downloading the Overwatch case")

	//Open networkDevice
	handle, err = pcap.OpenLive(networkDevice, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	//Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var first string
	var second string
	for packet := range packetSource.Packets() {
		//fmt.Printf("%s",packet.Data())
		if strings.Contains(string(packet.Data()), ".dem.bz2") {
			fmt.Println("Found the demo starting the download")
			first = string(packet.Data()[strings.Index(string(packet.Data()), "Host:")+6 : strings.Index(string(packet.Data()), "Accept:")-2])
			second = string(packet.Data()[strings.Index(string(packet.Data()), "GET")+4 : strings.Index(string(packet.Data()), "HTTP")-1])
			break
		}
	}
	fileUrl := "http://" + first + second
	err = DownloadFile("demo.dem.bz2", fileUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded: " + fileUrl + "\n\n")

	err = archiver.DecompressFile("demo.dem.bz2", "demo.dem")
	if err != nil {
		_ = os.Remove("demo.dem")
		err = archiver.DecompressFile("demo.dem.bz2", "demo.dem")
		if err != nil {
			panic(err)
		}
	}

	demo, err := os.Open("demo.dem")
	if err != nil {
		panic(err)
	}
	defer demo.Close()

	p := dem.NewParser(demo)
	defer p.Close()
	//Register handler on kill events
	p.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		if p.GameState().TotalRoundsPlayed()+1 >= owStartRound && owStartRoundSet {
			allplayers := p.GameState().Participants().Playing()
			fmt.Println("\n##########################################################################")
			fmt.Printf("Current Round: %d\n\n", p.GameState().TotalRoundsPlayed()+1)
			for _, player := range allplayers {
				var profileurl string
				var vacStatus string
				var gameBans int
				var daysSinceLastBan string
				if player.SteamID64 != 0 && steamWebApiKey != "" {
					//https://steamapi.xpaw.me/#ISteamUser/GetPlayerSummaries
					doc, _ := jsonquery.LoadURL("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=" + steamWebApiKey + "&steamids=" + strconv.FormatUint(player.SteamID64, 10))
					for _, n := range jsonquery.Find(doc, "response/players/*/profileurl") {
						profileurl = n.InnerText()
					}
					doc, _ = jsonquery.LoadURL("https://api.steampowered.com/ISteamUser/GetPlayerBans/v1/?key=" + steamWebApiKey + "&steamids=" + strconv.FormatUint(player.SteamID64, 10))
					for _, n := range jsonquery.Find(doc, "players/*/VACBanned") {
						vacStatus = n.InnerText()
					}
					for _, n := range jsonquery.Find(doc, "players/*/NumberOfVACBans") {
						gameBans, _ = strconv.Atoi(n.InnerText())
					}
					if vacStatus == "true" || gameBans > 0 {
						for _, n := range jsonquery.Find(doc, "players/*/DaysSinceLastBan") {
							daysSinceLastBan = n.InnerText()
						}
					}
				}
				var team string
				if player.Team == 2 {
					team = "T"
				} else {
					team = "CT"
				}

				var botName string
				if player.IsBot {
					botName = "BOT "
				} else {
					botName = ""
				}

				var banPrint string
				if vacStatus == "true" {
					banPrint = "VAC "
				}
				if gameBans > 0 {
					banPrint += "GameBans: " + strconv.Itoa(gameBans)
				}

				fmt.Printf("Team: %s ,Player: %s, SteamID64: %d, Profile: %s\n", team, botName+player.Name, player.SteamID64, profileurl)
				if banPrint != "" {
					fmt.Printf(ErrorColor, banPrint+" Days since last ban: "+daysSinceLastBan+"\n")
				}
				fmt.Printf("K: %d, A: %d, D: %d\n\n", player.Kills(), player.Assists(), player.Deaths())
			}
			fmt.Print("Advance to next round? [Press ENTER]")
			fmt.Scanln()
			//reader := bufio.NewReader(os.Stdin)
			//_, _ = reader.ReadString('\n')
		}
	})
	p.RegisterEventHandler(func(e events.MatchStart) {
		fmt.Printf("Map: %s\n", p.Header().MapName)
		fmt.Printf("Server: %s\n\n", p.Header().ServerName)

		fmt.Println("In which round did your Overwatch case start?")
		_, _ = fmt.Scanf("%d", &owStartRound)
		fmt.Scanln()
		owStartRoundSet = true
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\nEnd of Demo")
}

/**
DownloadFile Downloads the file to the provided filepath from the provided url
*/
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
