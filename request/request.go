package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Post(url string, param string, heard *map[string]string) (data []byte) {
	payload := strings.NewReader(param)

	req, _ := http.NewRequest("POST", url, payload)
	if heard != nil {
		for key, value := range *heard {
			req.Header.Add(key, value)
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	return body
}

func Get(url string, heard *map[string]string) ([]byte, error) {

	req, _ := http.NewRequest("GET", url, nil)

	// req.Header.Add("Authorization", "BNbmgAAGI155F6MJ3N2Tk9ruL_6XQpx-uxkkg:tGCY3xCsgybHd5IjcDMi9yZXBvcy93aF9mbG93RGF0YVNvdXJjZTEiLCJleHBpcmVzIjoxNTM2NzU4NjQ3LCJjb250ZW5VudFR5cGUiOiIiLCJoZWFkZXJzIjoiIiwibWV0aG9kIjoiR0VUIn0=")
	if heard != nil {
		for key, value := range *heard {
			req.Header.Add(key, value)
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	fmt.Println(res, 456)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), 456)
	// time.Sleep(time.Duration(3) * time.Second)
	defer res.Body.Close()
	return body, nil
}

func Put(url string, param string, heard *map[string]string) []byte {
	payload := strings.NewReader(param)
	req, _ := http.NewRequest("PUT", url, payload)
	if heard != nil {
		for key, value := range *heard {
			req.Header.Add(key, value)
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	fmt.Println(string(body))
	return body
}
func Put11(url string, filepath string, heard *map[string]string) (data []byte, err error) {
	param, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	data1 := bytes.NewReader(param)
	// data2 := strings.NewReader(string(param))
	// data3 := string(param)
	req, _ := http.NewRequest("PUT", url, data1)
	if heard != nil {
		for key, value := range *heard {
			req.Header.Add(key, value)
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	return body, nil
}

func Delete() {

	url := "http://xxxxx:8080/v2/repos/wh_flowDataSource1"

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("Authorization", "5F6MJ3N2Tk9ruL_6XQpx-uxkkg:o56-nIwtgTzUX80YCNpbcjUL8iM=:eyJyZXNvdXJjZSI6IF9mbG93RGF0YVNvdXJjZTEiLCJleHBpcmVzIjoxNTM2NzU4ODE2LCJjb250ZW50TUQ1IjoiIiwiY29udGVudFR5cGUiOiIiLCJoZWFkZXJzIjoiIiwibWV0aG9kIjoiREVMRVRFIn0=")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
