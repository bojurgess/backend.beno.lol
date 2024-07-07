package spotify

import (
	"encoding/json"
	"net/http"

	"github.com/bojurgess/backend.beno.lol/internal/database"
)

// this is NOT the same as database.GetUser()
func GetUser(tokens database.Tokens) (*database.User, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + tokens.AccessToken,
	}

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var ur MeResponse
	if err = json.NewDecoder(res.Body).Decode(&ur); err != nil {
		return nil, err
	}

	return &database.User{
		ID:          ur.Id,
		DisplayName: ur.DisplayName,
		Email:       ur.Email,
		Tokens:      tokens,
	}, nil
}
