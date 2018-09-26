package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	appconfig "./appconfig"
	reporter "./reporterimpl"
	watcher "./watcherimpl"

	dtparser "github.com/araddon/dateparse"
	lmsg "github.com/dmuth/google-go-log4go"
)

type ProdFatalLog struct {
	Instance         string
	WrittenTimeStamp string
	Message          string
}

func extractLogTimeStampFromCapturedLogEvents(firstLine string) string {
	lineSlice := strings.FieldsFunc(firstLine, func(divide rune) bool {
		return divide == ' '
	})
	tm, _ := dtparser.ParseLocal(lineSlice[1][1:])
	return fmt.Sprintf("%s", tm.Format("2006-01-02T15:04:05.000000"))
}

//curl -XPUT elasticsearch_domain_endpoint/movies/movie/1 -d '{"director": "Burton, Tim", "genre": ["Comedy","Sci-Fi"], "year": 1996, "actor": ["Jack Nicholson","Pierce Brosnan","Sarah Jessica Parker"], "title": "Mars Attacks!"}' -H 'Content-Type: application/json'
func composeESJSONEntryList(logFilePath string) (finalReportList []ProdFatalLog) {
	finalReportList = make([]ProdFatalLog, 0)

	thisInstance := watcher.GetThisInstanceName()
	lmsg.Tracef("thsiInstance => %s", thisInstance)

	reportList := watcher.ScanProductionLog(logFilePath)
	for itr := range reportList {
		if len(reportList[itr]) > 0 {
			pfl := new(ProdFatalLog)
			pfl.Instance = thisInstance
			pfl.WrittenTimeStamp = extractLogTimeStampFromCapturedLogEvents(reportList[itr][0])
			pfl.Message = strings.Join(reportList[itr], "")
			finalReportList = append(finalReportList, *pfl)

			// Logging
			if lmsg.Level() >= lmsg.TraceLevel {
				lmsg.Trace("======================================START====")
				tJSON, _ := json.Marshal(pfl)
				lmsg.Tracef("%s", tJSON)
				lmsg.Trace("======================================End=======")
			}
		}
	}
	return
}

func reportProgLogs() {
	insertCounter := 0
	//  UI Log
	watchUI := fmt.Sprintf("%s/%s", appconfig.WatchlogUI, appconfig.WatchlogFile)
	finalUIList := composeESJSONEntryList(watchUI)
	for i := range finalUIList {
		insertCounter++
		tmp, _ := json.Marshal(finalUIList[i])
		reporter.Upload(tmp, fmt.Sprintf("%s-%d", time.Now().Format("20060102-150405"), insertCounter))
	}

	//  Injector Log
	watchInjector := fmt.Sprintf("%s/%s", appconfig.WatchlogInjestor, appconfig.WatchlogFile)
	finalInjestorList := composeESJSONEntryList(watchInjector)
	for i := range finalInjestorList {
		insertCounter++
		tmp, _ := json.Marshal(finalInjestorList[i])
		reporter.Upload(tmp, fmt.Sprintf("%s-%d", time.Now().Format("20060102-150405"), insertCounter))
	}
	// Once upload is complete run GC
	runtime.GC()

}

func main() {

	lmsg.SetDisplayTime(false)
	lmsg.SetLevel(lmsg.TraceLevel)
	//lmsg.SetLevel(lmsg.ErrorLevel)

	// Report the very first Time, for all the fatal errors
	// found in the past one hr
	reportProgLogs()

	// report every timeInterval set
	for {
		timerLogWatcher := time.NewTimer(time.Minute * time.Duration(appconfig.WatchlogInterval))
		<-timerLogWatcher.C
		reportProgLogs()
	}

}
