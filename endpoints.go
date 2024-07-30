package playfab

import (
	"encoding/json"
	"fmt"
	"errors"
)

// Filter represents a filter for the PlayFab catalog search system.
type Filter struct {
	Count   bool   `json:"count"`
	Filter  string `json:"filter"`
	OrderBy string `json:"orderBy"`
	SCID    string `json:"scid"`
	Skip    int    `json:"skip,omitempty"`
	Limit   int    `json:"top"`
}

// Search searches the PlayFab catalog for items matching the given filter.
func (p *PlayFab) Search(filter Filter) (map[string]any, error) {
	m := make(map[string]any)
	if err := p.request("Catalog/Search", filter, &m, false); err != nil {
		return nil, err
	}
	return m["data"].(map[string]any), nil
}

type AccountData struct {
	PlayFabId     string `json:"MasterPlayerAccountId"`
	Namespace     string `json:"NamespaceId"`
	TitleId       string `json:"TitleId"`
	TitlePlayerId string `json:"TitlePlayerAccountId"`
}

// SearchAccounts searches for title player accounts based on their Xbox Live IDs.
func (p *PlayFab) SearchAccounts(xuids ...string) (map[string]AccountData, error) {
	m := make(map[string]any)
	data := struct {
		Sandbox     string   `json:"Sandbox"`
		XboxLiveIds []string `json:"XboxLiveIds"`
		TitleId     string   `json:"TitleId"`
	}{
		Sandbox:     "RETAIL",
		XboxLiveIds: xuids,
		TitleId:     "20CA2",
	}
	if err := p.request("Profile/GetTitlePlayersFromXboxLiveIDs", data, &m, false); err != nil {
		return nil, err
	}

	accounts := make(map[string]AccountData)
	for k, v := range m["data"].(map[string]any)["TitlePlayerAccounts"].(map[string]any) {
		if v == nil {
			// There's no account data for this XUID.
			continue
		}

		dMap := v.(map[string]any)
		accounts[k] = AccountData{
			PlayFabId:     dMap["MasterPlayerAccountId"].(string),
			Namespace:     dMap["NamespaceId"].(string),
			TitleId:       dMap["TitleId"].(string),
			TitlePlayerId: dMap["TitlePlayerAccountId"].(string),
		}
	}

	return accounts, nil
}

// ListFunctions prints out all the cloud functions on the PlayFab API.
func (p *PlayFab) ListFunctions() error {
	// POST https://titleId.playfabapi.com/CloudScript/ListFunctions
	m := make(map[string]any)
	if err := p.request("CloudScript/ListHttpFunctions", nil, &m, true); err != nil {
		return err
	}

	enc, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(enc))
	return nil
}

// TODO: Implement more endpoints.
