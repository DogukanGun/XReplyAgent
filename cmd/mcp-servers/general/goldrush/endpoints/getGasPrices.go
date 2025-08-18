package endpoints

import (
	"io"
	"net/http"
)

func (ge *GoldrushEndpoints) getGasPrices() string {

	url := ge.BaseUrl + "{chainName}/event/{eventType}/gas_prices/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body)
}
