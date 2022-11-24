package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	AocURL = "https://adventofcode.com/2022/leaderboard/private/view/"
)

func GetData(boardId string) (*Data, error) {
	b, err := ioutil.ReadFile("./" + boardId + ".json")
	if err != nil {
		return nil, err
	}

	var D Data
	if err = json.Unmarshal(b, &D); err != nil {
		return nil, err
	}
	return &D, nil
}

func FetchData(boardId, sessionToken, writePath string) error {
	req, err := http.NewRequest("GET", AocURL+boardId+".json", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", "session="+sessionToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("error fetching data from leaderboard: %s", boardId)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(writePath+".json", body, 0777); err != nil {
		return err
	}

	return nil
}
