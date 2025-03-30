package playfab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	requestTypeMaster = iota
	requestTypeTitleEntity
)

// PlayFab represents an instance of a Minecraft PlayFab client.
type PlayFab struct {
	src                        oauth2.TokenSource
	customId                   string
	client                     *http.Client
	id, token                  string
	titleAccountId, titleToken string
}

// New creates a new PlayFab client with the given token source.
func New(client *http.Client, src oauth2.TokenSource) (*PlayFab, error) {
	p := &PlayFab{
		src:    src,
		client: client,
	}
	if err := p.acquireLoginToken(); err != nil {
		return nil, err
	}
	if err := p.acquireEntityToken(p.id, "master_player_account"); err != nil {
		return nil, err
	}
	if err := p.aquireTitleAccount(minecraftTitleID); err != nil {
		return nil, err
	}
	return p, nil
}

// New creates a new PlayFab client with the given custom ID.
func NewWithCustomID(client *http.Client, customId string) (*PlayFab, error) {
	p := &PlayFab{
		customId: customId,
		client:   client,
	}
	if err := p.loginWithCustomId(); err != nil {
		return nil, err
	}
	if err := p.acquireEntityToken(p.id, "master_player_account"); err != nil {
		return nil, err
	}
	if err := p.aquireTitleAccount(minecraftEduTitleID); err != nil {
		return nil, err
	}
	return p, nil
}

// ID returns the PlayFab ID
func (p *PlayFab) ID() string {
	return p.id
}

// request sends a request to the PlayFab API.
func (p *PlayFab) request(url string, body any, res any, title bool) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s.playfabapi.com/%s", minecraftTitleID, url), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", minecraftUserAgent)
	req.Header.Set("X-PlayFabSDK", minecraftDefaultSDK)
	req.Header.Set("X-ReportErrorAsSuccess", "true")
	if !title && len(p.token) > 0 {
		req.Header.Set("X-EntityToken", p.token)
	} else if title && len(p.titleToken) > 0 {
		req.Header.Set("X-EntityToken", p.titleToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&res)
}
