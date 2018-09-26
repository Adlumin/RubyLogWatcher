package watcherimpl

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"

	appconfig "../appconfig"
	dtparser "github.com/araddon/dateparse"
	lmsg "github.com/dmuth/google-go-log4go"
)

// isEventWithinTimeFrame : is the check if the Timestamp occured in the time window
func isEventWithinTimeFrame(chkTimeStamp, strt, end time.Time) bool {
	rangeInMin := float64(appconfig.WatchlogInterval + 1)
	chk1 := strt.Sub(chkTimeStamp).Minutes() <= rangeInMin
	chk2 := chkTimeStamp.Sub(end).Minutes() <= rangeInMin
	if chk1 && chk2 {
		return true
	}
	return false
}

func isNewLogEntryFatalandWithinTimeFrame(line string, beginTF, endTF time.Time) (isTaggedAsFatal, isNewLogEntry bool) {
	lineSlice := strings.FieldsFunc(line, func(divide rune) bool {
		return divide == ' '
	})

	// Fatal Message Capture ON and OFF
	if len(lineSlice) > 2 {
		// line split successful
		var isFatal, isWithinTF bool
		if lineSlice[0] == "F," {
			isFatal = true
		}

		chkTimeStamp := lineSlice[1][1:]
		tm, etm := dtparser.ParseLocal(chkTimeStamp)
		if etm != nil {
			lmsg.Errorf("%s", etm)
			isNewLogEntry = false
		} else {
			isNewLogEntry = true
		}
		// lmsg.Tracef("string => %s", chkTimeStamp)
		// lmsg.Tracef("toTime() => %s", tm)

		if isEventWithinTimeFrame(tm, beginTF, endTF) {
			isWithinTF = true
		}
		// This confirm line of of this format
		//F, [2018-08-27T14:31:53.788969 #1483] FATAL -- :
		if isFatal && isWithinTF {
			isTaggedAsFatal = true
		}
	}
	return
}

func ScanProductionLog(filePath string) (foundALLResult [][]string) {
	var result []string
	foundALLResult = make([][]string, 0)
	GrabVerbage := false
	// Original  -- start
	endTF := time.Now()
	beginTF := endTF.Add(time.Minute * time.Duration(-(appconfig.WatchlogInterval + 1)))
	// Original  -- stop

	// // TestOnly -- new start
	// now := time.Now()
	// endTF := now.Add(time.Minute * time.Duration(-1230))
	// beginTF := endTF.Add(time.Minute * time.Duration(-61))
	// lmsg.Infof("=================Timer Range====================")
	// lmsg.Infof("beginTF => %s", beginTF)
	// lmsg.Infof("endTF => %s", endTF)
	// lmsg.Infof("=================Timer Range====================")
	// // TestOnly  -- new stop

	//pathRsysLog := "../SampleInputLogs/pro.log"
	file, err := os.Open(filePath)
	if err != nil {
		lmsg.Errorf("Cannot Open File %s Error: %s", filePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineNum := 0
	LastLineCapt := -1
	for {
		lineNum++
		line, err := reader.ReadString('\n')

		isStartofFatal, isNewLogEntry := isNewLogEntryFatalandWithinTimeFrame(line, beginTF, endTF)
		if isStartofFatal && isNewLogEntry {
			if len(result) > 0 {
				//Capture anything from earlier Iteration , if found a new log entry
				foundALLResult = append(foundALLResult, result)
			}
			GrabVerbage = true
			result = make([]string, 0)
			result = append(result, line)
			LastLineCapt = lineNum
		}

		if GrabVerbage {
			if !isNewLogEntry {
				result = append(result, line)
				LastLineCapt = lineNum

			}
			if isNewLogEntry && (LastLineCapt+1 == lineNum) {
				GrabVerbage = false
			}
		}

		if err == io.EOF {
			foundALLResult = append(foundALLResult, result) // capture the lastOne
			break
		}

		if err != nil {
			lmsg.Errorf("Parser Error on File %s Error: %s", filePath, err)
			foundALLResult = append(foundALLResult, result)
			return
		}
	}

	for itr := range foundALLResult {
		lmsg.Infof("==================")
		lmsg.Warnf("%v", foundALLResult[itr])
		lmsg.Infof("==================")
	}
	lmsg.Errorf("%d", len(foundALLResult))
	return
}
