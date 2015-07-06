package newsGetter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	getTimeout = time.Duration(3 * time.Second)
)

type unMarshalledContent map[string]interface{}

// func httpGet(url_string string) (*http.Response, error) {
// 	response, err := http.Get(url_string)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	return response, nil
// }

// TODO add aditional flag if resp err is nil but fail
func httpGet(urlString string) (*http.Response, error) {

	client := &http.Client{
		Timeout: getTimeout,
	}

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil || resp == nil || resp.StatusCode != 200 {
		// fmt.Println("status code: ", resp)
		// fmt.Println("err: ", err)
		// fmt.Println("-----------------------------------")
		return nil, err
	}
	return resp, nil
}

func responseReader(response *http.Response) ([]byte, error) {
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func unmarshalResponseContent(content []byte, dataContainer interface{}) (interface{}, error) {
	if err := json.Unmarshal(content, &dataContainer); err != nil {
		// fmt.Println(err)
		return nil, err
	}
	return dataContainer, nil
}
