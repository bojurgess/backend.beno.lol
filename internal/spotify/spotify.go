// Useful helpers for grabbing Spotify API data.
package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/bojurgess/backend.beno.lol/internal/database"
	"github.com/bojurgess/backend.beno.lol/internal/util"
)

// Gets the user's spotify profile information.
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

// Refreshes spotify access token.
// Does not update the database, this is expected to be done by the consumer.
// Returns a full database.Tokens struct.
func RefreshAccessToken(tokens *database.Tokens, cid string, cs string) error {
	const url = "https://accounts.spotify.com/api/token"
	body := util.MapToQuerystring(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": tokens.RefreshToken,
	})
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": util.EncodeAuth(cid, cs),
	}

	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var tr TokenResponse
	if err = json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return err
	}

	if tr.TokenErrorResponse != nil {
		return errors.New(tr.ErrorDescription)
	}

	t := MapTokenResponse(tr)
	tokens.AccessToken = t.AccessToken
	tokens.ExpiresAt = t.ExpiresAt

	return nil
}

// Grabs the user's currently playing track.
// Everything in the token struct, excluding the access token, can be nil
// Returns a NowPlayingResponse struct.
func GetNowPlaying(tokens database.Tokens) (*NowPlayingResponse, error) {
	// I think the caller should probably be handling refreshes,
	// otherwise we make the behaviour of this helper more opaque.
	if time.Now().After(tokens.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	var response NowPlayingResponse

	url := "https://api.spotify.com/v1/me/player/currently-playing"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}

	return &response, nil
}

// Maps the returned spotify token response to the database token struct
func MapTokenResponse(tr TokenResponse) database.Tokens {
	return database.Tokens{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second),
		TokenType:    tr.TokenType,
		Scope:        tr.Scope,
	}
}
