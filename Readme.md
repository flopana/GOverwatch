# GOverwatch
GOverwatch is a tool written in go to get the actual Demo from CS:GO Overwatch and find out who's the suspect.

## Installation
Download the goverwatch.exe including the config.json.
Refer to the Configuration part in this README or instructions by the program.

Also you need to install winpcap which is used to sniff your network traffic.

Alternatively clone this repo and compile it yourself
```go build -o goverwatch.exe main.go```

## Usage
Start GOverwatch, after finding the Demo it will prompt you and ask in which round your overwatch case started at.

![alt text](https://i.imgur.com/lbsgVlp.png "GOverwatch Demo 01")

After the last player has been parsed GOverwatch will prompt you again and asks if it should advance to the next round.

Whith this you are able to find the real suspect by comparing the stats.

## Configuration
To be able to get the profile links aswell you need to fill in your SteamWebApiKey which you can obtain on this site.
https://steamcommunity.com/dev/apikey
```json
{
  "steamWebApiKey": "youre API Key here",
  "networkDevice": "\\Device\\NPF_{ABCDE-EFGHI-JKMLNOP123456}"
}
```
Notice the two backticks to mask the second backtick otherwise json would interpret something there

In conclusion convert:

```\Device\NPF_{ABCDE-EFGHI-JKMLNOP123456}```

to

```\\Device\\NPF_{ABCDE-EFGHI-JKMLNOP123456}```

**How do I find my device?**

Simply don't fill out the "networkDevice" GOverwatch will list all available devices

For further questions don't hesitate to open an issue