package endpoints

import (
	"fmt"
	"io"
	"net/http"
)

func (ge *GoldrushEndpoints) getTransaction() string {

	url := ge.BaseUrl + "{chainName}/transaction_v2/{txHash}/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}
