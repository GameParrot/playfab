package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gameparrot/playfab"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
)

func main() {
	src := tokenSource()
	db, err := playfab.New(http.DefaultClient, src)
	if err != nil {
		panic(err)
	}

	t := time.Now()
	accounts, err := db.SearchAccounts("2535447073711921")
	delta := float64(time.Since(t).Microseconds()) / 1000
	if err != nil {
		panic(err)
	}
	fmt.Printf("Took %fms\n", delta)

	for xuid, dat := range accounts {
		fmt.Printf("%s: %s [%s]\n", xuid, dat.PlayFabId, dat.TitlePlayerId)
	}

	db.ListFunctions()

	/* resp, err := db.Search(playfab.Filter{
		Count:   true,
		Filter:  "(contentType eq 'PersonaDurable' and displayProperties/pieceType eq 'persona_emote')",
		OrderBy: "creationDate desc",
		SCID:    "4fc10100-5f7a-4470-899b-280835760c07",
		Limit:   300,
	})
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		panic(err)
	}

	os.WriteFile("test.json", b, 0644) */
}

func tokenSource() oauth2.TokenSource {
	token := new(oauth2.Token)
	data, err := os.ReadFile("token.tok")
	if err == nil {
		_ = json.Unmarshal(data, token)
	} else {
		token, err = auth.RequestLiveToken()
		if err != nil {
			panic(err)
		}
	}
	src := auth.RefreshTokenSource(token)
	_, err = src.Token()
	if err != nil {
		token, err = auth.RequestLiveToken()
		if err != nil {
			panic(err)
		}
		src = auth.RefreshTokenSource(token)
	}
	tok, _ := src.Token()
	b, _ := json.Marshal(tok)
	_ = os.WriteFile("token.tok", b, 0644)
	return src
}
