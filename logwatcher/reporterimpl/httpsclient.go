package reporterimpl

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	appconfig "../appconfig"
	lmsg "github.com/dmuth/google-go-log4go"
)

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

func httpsPOST(burst json.RawMessage, postURL string) (isUploadSuccessfull bool) {
	body := bytes.NewReader(burst)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	req, err := http.NewRequest("POST", postURL, body)
	if err != nil {
		lmsg.Errorf("Create Request Error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Transport: transCfg}
	resp, err := client.Do(req)
	if err != nil {
		lmsg.Errorf("Https Request Error %s", err.Error())
	}
	defer resp.Body.Close()
	if resp != nil && resp.StatusCode != 200 {
		lmsg.Errorf("response %v", resp)
		isUploadSuccessfull = false
	} else {
		lmsg.Debugf("*** === Success === ***")
		lmsg.Debugf("response %s", resp)
		isUploadSuccessfull = true
	}
	return
}

func Upload(data json.RawMessage, ESid string) {
	esURL := fmt.Sprintf("%s/%s/%s/%s", appconfig.ESDomain, appconfig.ESIndex, appconfig.ESType, ESid)
	Block{
		Try: func() {
			httpsPOST(data, esURL)
		},
		Catch: func(e Exception) {
			lmsg.Errorf("Caught Exception := %v\n", e)
		},
	}.Do()
}
