package endpoints

import (
	"fmt"
	"io"
	"net/http"
)

func (ge *GoldrushEndpoints) getTokenBalanceForAddress() string {

	url := "https://api.covalenthq.com/v1/{chainName}/address/{walletAddress}/balances_v2/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer <token>")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}
