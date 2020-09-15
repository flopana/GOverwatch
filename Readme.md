# Goverwatch
GOverwatch is a tool written in go to get the actual Demo from CS:GO Overwatch and find out who's the suspect.

## Installation
TODO

## Usage
Star GOverwatch after finding the Demo it will prompt you and ask in which round your overwatch case started at.

![alt text](https://i.imgur.com/M2loJ4x.png "GOverwatch Demo 01")

After the last player has been parsed GOverwatch will prompt you again and asks if it should advance to the next round.

Whith this you are able to find the real suspect by comparing the stats.

## Configuration
To be able to get the profile links aswell you need to fill in your SteamWebApiKey which you can obtain on this site.
https://steamcommunity.com/dev/apikey
```json
{
  "steamWebApiKey": "youre API Key here"
  "networkDevice": "\\\\Device\\\\NPF_{ABCDE-EFGHI-JKMLNOP123456}"
}
```