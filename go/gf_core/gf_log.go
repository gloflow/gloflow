/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_core

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

//-------------------------------------------------

type GFlogFun func(string, string, map[string]interface{})

//-------------------------------------------------

func LogsInit() (func(string, string), GFlogFun) {
	return LogsInitNew(true, "")
}

// pLogLevelStr - allows for explicit programmatic setting of log_level.
//                if this is set to "" then the ENV var is checked.
//                if this is also not set then the default "info" level is used.
func LogsInitNew(pLogrusBool bool, pLogLevelStr string) (func(string, string), GFlogFun) {

	//--------------------
	// LOGRUS_INIT

	if pLogrusBool {

		//--------------------
		// LOG_LEVEL

		var logLevelStr string
		logLevelDefaultStr := "info"

		if pLogLevelStr != "" {
			logLevelStr = pLogLevelStr
		} else {
			logLevelENVstr := os.Getenv("GF_LOG_LEVEL")
			if logLevelENVstr == "" {
				logLevelStr = logLevelDefaultStr
			} else {
				logLevelStr = logLevelENVstr
			}
		}

		fmt.Printf("log level - %s\n", logLevelStr)
		
		level, err := logrus.ParseLevel(logLevelStr)
		if err != nil {
			fmt.Println(err)
			panic("log level is not valid, has to be : trace, debug, info, warning, error, fatal, panic")
		}

		// set loging severity level, and above. 
		// the level that was set will include that level and all severities higher than that
		logrus.SetLevel(level)
		
		//--------------------

		// log := logrus.New()
		
		// log to stdout instead of the default stderr
		logrus.SetOutput(os.Stdout)

		// Log as JSON instead of the default ASCII formatter.
		// log.SetFormatter(&log.JSONFormatter{})

		logrus.SetFormatter(&logrus.TextFormatter{
			// DisableColors: true,
			FullTimestamp: true,
		})

		// adds the caller function as 'method' to the output
		// log.SetReportCaller(true)
	}

	//--------------------

	green  := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red    := color.New(color.FgRed).SprintFunc()
	
	//-------------------------------------------------
	// DEPRECATED!! - when all logging is migrated to logNewFun delete this function.
	logFun := func(pGroupStr string, pMsgStr string) {
		timeStr := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0, 'f', 10, 64)

		if pGroupStr == "FUN_ENTER" {
			fmt.Printf(timeStr+":"+yellow(pGroupStr)+":"+pMsgStr+"\n")
		} else if pGroupStr == "INFO" {
			fmt.Printf(timeStr+":"+green(pGroupStr)+":"+green(pMsgStr)+"\n")
		} else if pGroupStr == "ERROR" {
			fmt.Printf(timeStr+":"+red(pGroupStr)+":"+pMsgStr+"\n")
		}
	}

	//-------------------------------------------------
	// IMPORTANT!! - migrate all loging to this function

	logNewFun := func(pLevelStr string, pMsgStr string, pMetaMap map[string]interface{}) {

		logFields := logrus.Fields{}
		if pMetaMap != nil {
			for k, v := range pMetaMap {
				logFields[k] = v
			}
			// contextLogger := logrus.WithFields(logFields)
		}

		switch pLevelStr {

		// very low-level logs
		case "TRACE":
			logrus.WithFields(logFields).Trace(pMsgStr)

		// debugging logs for devs or troubleshooting
		case "DEBUG":
			logrus.WithFields(logFields).Debug(pMsgStr)

		// informative logs for general observation
		case "INFO":
			logrus.WithFields(logFields).Info(pMsgStr)
			
		case "WARNING":
			logrus.WithFields(logFields).Warn(pMsgStr)

		// failure occured, but the process is not exiting
		case "ERROR":
			logrus.WithFields(logFields).Error(pMsgStr)

		// process will exit after this log
		case "FATAL":
			logrus.WithFields(logFields).Fatal(pMsgStr)
		}
	}

	//-------------------------------------------------
	return logFun, logNewFun
}

//-------------------------------------------------

func LogsIsDebugEnabled() bool {
	level, _ := logrus.ParseLevel("debug")
	return logrus.IsLevelEnabled(level)
}