package main

import (
	"testing"
)

func Test_DealMsg(t *testing.T) {
	line := `{"time":1519812001285,"logid":"bc03b7fa9260fb6c62a4e4a8fb56fb3e","hostname":"www-laoshi-9-94","prelogid":"228a5189d29db3cccffa416f207a3a11","pageid":"","sessid":"","appid":1000006,"token":"","ver":"1","os":"","ua":"okhttp\/3.8.1","ip":"120.193.204.120","xesid":"","userid":"5758005","devid":"","data":[{"req":{"url":"laoshi.xueersi.com\/LiveLecture\/getTestInfoForPlayBack","params":{"url":"LiveLecture\/getTestInfoForPlayBack","pageid":"PublicLiveDetailActivity","datalogid":"35635ee748dcc3f2","systemName":"android","timeStr":"0","appChannel":"oppo","liveId":"39025","systemVersion":"7.1.1","teacherId":"2510","requesttime":"1519812000065","appVersion":"6.3.05","identifierForClient":"329d6de4d4ba1ed4b7959057aa70f67d","logid":"717fe48f49887861","appVersionNumber":"60305"}}}]}`
	log := []byte(line)
	// _, err := DealMsg(log)
	// if err != nil {
	// 	t.Error(err)
	// }
	DealMsg(log)
}
