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
	"fmt"
	"time"
	"strconv"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)
//-------------------------------------------------

func InitLogs() (func(string, string), func(string, string, string, map[string]interface{})) {

	
	log := logrus.New()
	
	// log to stdout instead of the default stderr
	// log.SetOutput(os.Stdout)
	// log.SetLevel(log.ErrorLevel)
	
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// adds the caller as 'method' to the output
	// log.SetReportCaller(true)

	green  := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red    := color.New(color.FgRed).SprintFunc()
	
	//-------------------------------------------------
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
	logNewFun := func(pMsgStr string, pGroupStr string, pLevelStr string, pMetaMap map[string]interface{}) {
		if pMetaMap != nil {
			logFields := logrus.Fields{}
			for k, v := range pMetaMap {
				logFields[k] = v
			}

			contextLogger := log.WithFields(logFields)

			switch pLevelStr {
			case "INFO":
				contextLogger.Info(pMsgStr)
			case "WARNING":
				contextLogger.Warn(pMsgStr)
			case "ERROR":
				contextLogger.Error(pMsgStr)
			}
		}
	}

	//-------------------------------------------------
	return logFun, logNewFun
}