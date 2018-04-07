package helpers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// URL ...
const URL string = "https://statsapi.web.nhl.com/api/v1/schedule?expand=schedule.teams,schedule.scoringplays"

func Get() JSON {
	res, err := http.Get(URL)

	if err != nil {
		// handle error
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// handle error
	}

	return getJSON(body)

}

func getJSON(body []byte) JSON {

	var response JSON

	json.Unmarshal(body, &response)

	return response
}
