package endpoints

import (
	"fmt"
	"io"
	"net/http"
)

func (ge *GoldrushEndpoints) getHistoricalTokenPrices() string {

	url := ge.BaseUrl + "pricing/historical_by_addresses_v2/{chainName}/{quoteCurrency}/{contractAddress}/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}
