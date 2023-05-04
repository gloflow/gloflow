/*
MIT License

Copyright (c) 2023 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_rpc_lib

import (
	"fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------------------------------
type GFmmsOnSessionJoinFun func(map[string]interface{}) map[string]interface{}

type GFmmsServerInfo struct {
	OnSessionJoinFun GFmmsOnSessionJoinFun
}

type GFmmsSession struct {
	ClientsJoinedThresholdInt int // number of clients that are expected to join the session at a minimum
	ClientsJoinedLst          []int 
	ClientNextAvailableIDint  int
	ClientsCountInt int
}

//------------------------------------------------------------------------
// SESSION_JOIN
func MMSsessioJoin(pUserNameStr string,
	pMetaMap    map[string]interface{},
	pSession    *GFmmsSession,
	pServerInfo GFmmsServerInfo,
	pRuntimeSys *gf_core.RuntimeSys) (int, map[string]interface{}) {

	newClientIDint := pSession.ClientNextAvailableIDint

	pSession.ClientsJoinedLst = append(pSession.ClientsJoinedLst, newClientIDint)
	pSession.ClientNextAvailableIDint += 1
	pSession.ClientsCountInt += 1

	// call user-defined session_join callback
	sessionMetaMap := pServerInfo.OnSessionJoinFun(pMetaMap)

	return newClientIDint, sessionMetaMap
}

//------------------------------------------------------------------------
// SESSION_STATUS
func MMSsessionStatus() {


}

//------------------------------------------------------------------------
// SESSION_RESET
func MMSsessionReset(pSessionNameStr string,
	pUserNameStr string,
	pSession     *GFmmsSession,
	pRuntimeSys  *gf_core.RuntimeSys) {


}

//------------------------------------------------------------------------

func MMSsessionGetOrCreate(pNameStr string,
	pSessionMap map[string]*GFmmsSession) *GFmmsSession {

	clientsJoinedThresholdInt := 2

	if _, ok := pSessionMap[pNameStr]; !ok {

		session := &GFmmsSession{
			ClientsJoinedThresholdInt: clientsJoinedThresholdInt,
			ClientNextAvailableIDint:  0,
			ClientsCountInt: 0,
		}
		pSessionMap[pNameStr] = session
		return session
	} else {
		return pSessionMap[pNameStr]
	}
	return nil
}

//------------------------------------------------------------------------
func MMSinitHandlers(pServerInfo GFmmsServerInfo,
	pHTTPmutex  *http.ServeMux,
	pRuntimeSys *gf_core.RuntimeSys) {




	sessionsMap := map[string]*GFmmsSession{}

	
 
	//------------------------------------------------------------------------
	CreateHandlerHTTPwithMux("/v1/session/reset",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//------------------
				// INPUT
				iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				sessionNameStr := iMap["session_name_str"].(string)
				fmt.Println("session_name", sessionNameStr)

				userNameStr := iMap["user_name_str"].(string)
				fmt.Println("user_name", userNameStr)

				//------------------

				session := MMSsessionGetOrCreate(sessionNameStr, sessionsMap)

				//------------------
				// MMS_SESSION_JOIN
				MMSsessionReset(sessionNameStr,
					userNameStr,
					session,
					pRuntimeSys)

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
				}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		pHTTPmutex,
		nil, // metrics,
		false, // pStoreRunBool
		nil, 
		pRuntimeSys)

	//------------------------------------------------------------------------
	CreateHandlerHTTPwithMux("/v1/session/join",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//------------------
				// INPUT
				iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				sessionNameStr := iMap["session_name_str"].(string)
				fmt.Println("session_name", sessionNameStr)

				userNameStr := iMap["user_name_str"].(string)
				fmt.Println("user_name", userNameStr)

				metaMap := iMap["meta_map"].(map[string]interface{})

				//------------------

				session := MMSsessionGetOrCreate(sessionNameStr, sessionsMap)

				//------------------
				// MMS_SESSION_JOIN
				newClientIDint, sessionMetaMap := MMSsessioJoin(userNameStr,
					metaMap,
					session,
					pServerInfo,
					pRuntimeSys)

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"client_id_int":    newClientIDint,
					"session_meta_map": sessionMetaMap,
				}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		pHTTPmutex,
		nil, // metrics,
		false, // pStoreRunBool
		nil, 
		pRuntimeSys)
	
	//------------------------------------------------------------------------
	// used to get status of a session, by clients that have already joined a session

	CreateHandlerHTTPwithMux("/v1/session/status",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//------------------
				// INPUT
 				iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				
				sessionNameStr := iMap["session_name_str"].(string)
				fmt.Println("session_name", sessionNameStr)

				clientIDint := int(iMap["client_id_int"].(float64))
				fmt.Println("client_id", clientIDint)

				//------------------

				session := MMSsessionGetOrCreate(sessionNameStr, sessionsMap)

				var allUsersJoinedBool bool
				if session.ClientsJoinedThresholdInt >= session.ClientsCountInt {
					allUsersJoinedBool = true
				}


				MMSsessionStatus()


				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"game_map": map[string]interface{}{
						"ALL_USERS_JOINED": allUsersJoinedBool,
					},
				}
				return dataMap, nil

				//------------------
			}
			return nil, nil
		},
		pHTTPmutex,
		nil, // metrics,
		false, // pStoreRunBool
		nil, 
		pRuntimeSys)

	//------------------------------------------------------------------------
}