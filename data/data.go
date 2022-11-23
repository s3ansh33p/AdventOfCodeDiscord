package data

import (
	"io/ioutil"
	"net/http"
)

const (
	AocURL = "https://adventofcode.com/2022/leaderboard/"
)

func GetData(boardId string) ([]byte, error) {
	return ioutil.ReadFile("./" + boardId + ".json")
}

func FetchData(boardId, writePath string) error {
	req, err := http.NewRequest("GET", AocURL+boardId, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", "session="+AocToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
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
