package awards

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	conf "github.com/SteMak/house-tyan/config"
)

var (
	endpointUnbelievaBoatAPI = url.URL{
		Scheme: "https",
		Host:   "unbelievaboat.com",
		Path:   "/api/v1/",
	}
)

// Balance balance of user
type responce struct {
	Rank   string `yaml:"rank"`
	UserID string `yaml:"user_id"`
	Cash   int    `yaml:"cash"`
	Bank   int    `yaml:"bank"`
	Total  int    `yaml:"total"`
}

// unbelievaBoat is a magic structure
type unbelievaBoat struct {
	token  string
	client *http.Client
}

func (api *unbelievaBoat) request(protocol, userID string, reqBodyBytes io.Reader) (*responce, error) {
	// RateLimit is a structure for 429 error
	type RateLimit struct {
		Message    string `yaml:"message"`
		RetryAfter int    `yaml:"retry_after"`
	}

	// JSONBalanse is a structure for changing user balance
	type JSONBalanse struct {
		Cash   int    `yaml:"cash"`
		Bank   int    `yaml:"bank"`
		Reason string `yaml:"reason"`
	}

	var (
		err   error
		b     responce
		limit RateLimit
	)

	endpoint := endpointUnbelievaBoatAPI
	endpoint.Path = path.Join("guilds", conf.Bot.GuildID, "users", userID)

	req, err := http.NewRequest(protocol, endpoint.String(), reqBodyBytes)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", api.token)

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		resBodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resBodyBytes, &b)
		if err != nil {
			return nil, err
		}
	}

	if res.StatusCode == 429 {
		resBodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resBodyBytes, &limit)
		if err != nil {
			return nil, err
		}

		time.Sleep(time.Duration(limit.RetryAfter) * time.Millisecond)
		return api.request(protocol, userID, reqBodyBytes)
	}

	if res.StatusCode != http.StatusOK {
		return &b, errors.New("Strange status code: " + strconv.Itoa(res.StatusCode))
	}

	return &b, nil
}

// GetBalance return balance of user
func (api *unbelievaBoat) getBalance(userID string) (*responce, error) {
	return api.request("GET", userID, nil)
}

// SetBalance sets balance of user
func (api *unbelievaBoat) setBalance(userID string, cash, bank int, reason string) error {
	type JSONBalanse struct {
		Cash   int    `yaml:"cash"`
		Bank   int    `yaml:"bank"`
		Reason string `yaml:"reason"`
	}

	jsonBalanse := JSONBalanse{
		Cash:   cash,
		Bank:   bank,
		Reason: reason,
	}

	reqBodyBytes, err := json.Marshal(jsonBalanse)
	if err != nil {
		return err
	}

	_, err = api.request("PUT", userID, bytes.NewBuffer(reqBodyBytes))
	return err
}

func (api *unbelievaBoat) addToBalance(userID string, cash, bank int, reason string) error {
	type JSONBalanse struct {
		Cash   int    `yaml:"cash"`
		Bank   int    `yaml:"bank"`
		Reason string `yaml:"reason"`
	}

	jsonBal := JSONBalanse{
		Cash:   cash,
		Bank:   bank,
		Reason: reason,
	}

	reqBodyBytes, err := json.Marshal(jsonBal)
	if err != nil {
		return err
	}

	_, err = api.request("PATCH", userID, bytes.NewBuffer(reqBodyBytes))
	return err
}
