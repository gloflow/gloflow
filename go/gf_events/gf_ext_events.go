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

package gf_events

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type EventsRegisterProducerMsg struct {
	EventsIDstr string
}

type EventsRegisterConsumerMsg struct {
	EventsIDstr string

	// IMPORTANT!! - channel on which the events_consumer is expecting to receive new 
	//               events produced by the events producer,once consumers registration 
	//               (by processing the EventsRegisterConsumerMsg message) is complete.
	ResponseCh chan chan EventMsg 
}

type EventMsg struct {
	EventsIDstr string                  `json:"events_id_str"`
	TypeStr     string                 `json:"type_str"`
	MsgStr      string                 `json:"msg_str"`
	DataMap     map[string]interface{} `json:"meta_map"`
}

type EventsCtx struct {
	RegisterProducerCh chan EventsRegisterProducerMsg
	RegisterConsumerCh chan EventsRegisterConsumerMsg
	EventsBrokerCh     chan EventMsg
}

//-------------------------------------------------

func SendEvent(pEventsIDstr string,
	pTypeStr    string,
	pMsgStr     string,
	pDataMap    map[string]interface{},
	pEventsCtx  *EventsCtx,
	pRuntimeSys *gf_core.RuntimeSys) {

	e := EventMsg{
		EventsIDstr: pEventsIDstr,
		TypeStr:     pTypeStr,
		MsgStr:      pMsgStr,
		DataMap:     pDataMap,
	}
	pEventsCtx.EventsBrokerCh <- e
}

//-------------------------------------------------

func RegisterProducer(pEventsIDstr string,
	pEventsCtx  *EventsCtx,
	pRuntimeSys *gf_core.RuntimeSys) {

	registerProducerMsg := EventsRegisterProducerMsg{
		EventsIDstr: pEventsIDstr,
	}

	pEventsCtx.RegisterProducerCh <- registerProducerMsg
}

//-------------------------------------------------

func Init(pSSEurlStr string, pRuntimeSys *gf_core.RuntimeSys) *EventsCtx {

	// yellow := color.New(color.FgYellow).SprintFunc()
	// black  := color.New(color.FgBlack).Add(color.BgYellow).SprintFunc()

	register_producer_ch := make(chan EventsRegisterProducerMsg, 50)
	register_consumer_ch := make(chan EventsRegisterConsumerMsg, 50)
	events_broker_ch     := make(chan EventMsg,                  500)

	events_consumers_map := map[string][]chan EventMsg{}
	go func() {
		for ;; {

			select {

				//-----------------
				// REGISTER EVENTS_PRODUCER
				case register_producer_msg := <- register_producer_ch:
					events_id_str                      := register_producer_msg.EventsIDstr
					events_consumers_map[events_id_str] = make([]chan EventMsg, 0)

				//-----------------
				// REGISTER EVENTS_CONSUMER
				case register_consumer_msg := <- register_consumer_ch:
					events_id_str                      := register_consumer_msg.EventsIDstr
					consumer_ch                        := make(chan EventMsg, 50)
					events_consumers_map[events_id_str] = append(events_consumers_map[events_id_str], consumer_ch)
				
					register_consumer_msg.ResponseCh <- consumer_ch

				//-----------------
				// EVENT_MSG RELAY
				case event_msg := <- events_broker_ch:
					events_id_str := event_msg.EventsIDstr

					// IMPORTANT!! - check that this events_id_str has consumers registered for it.
					//               if yes, then get a list of all consumers for this events_id_str,
					//               and go through that list sending the same event_msg to all of them
					//               (multicast style)
					if consumers_lst, ok := events_consumers_map[events_id_str]; ok {
						for _, consumer_ch := range consumers_lst {
							consumer_ch <- event_msg
						}
					}
					
				//-----------------
			}
		}
	}()

	ctx := &EventsCtx{
		RegisterProducerCh: register_producer_ch,
		RegisterConsumerCh: register_consumer_ch,
		EventsBrokerCh:     events_broker_ch,
	}

	initHandlers(pSSEurlStr,
		register_consumer_ch,
		ctx,
		pRuntimeSys)
	return ctx
}

//-------------------------------------------------

func initHandlers(pSSEurlStr string,
	pRegisterConsumerCh chan<- EventsRegisterConsumerMsg,
	pEventsCtx          *EventsCtx,
	pRuntimeSys         *gf_core.RuntimeSys) {

	// yellow := color.New(color.FgYellow).SprintFunc()
	// black  := color.New(color.FgBlack).Add(color.BgYellow).SprintFunc()


	// IMPORTANT!! - new event_consumers (clients) register via this HTTP handler
	http.HandleFunc(pSSEurlStr, func(p_resp http.ResponseWriter, p_req *http.Request) {
		pRuntimeSys.LogFun("INFO", "INCOMING HTTP REQUEST -- "+pSSEurlStr+" ----------")

		eventsIDstr := p_req.URL.Query()["events_id"][0]
		pRuntimeSys.LogFun("INFO", "events_id_str - "+eventsIDstr)

		register_consumer__response_ch := make(chan chan EventMsg)
		register_consumer_msg          := EventsRegisterConsumerMsg{
			EventsIDstr: eventsIDstr,
			ResponseCh:  register_consumer__response_ch,
		}

		pRegisterConsumerCh <- register_consumer_msg
		eventsConsumerCh := <- register_consumer__response_ch

		flusher, gf_err := gf_core.HTTPinitSSE(p_resp, pRuntimeSys)
		if gf_err != nil {
			return
		}

		//-------------
		// SEND_EVENT

		eventTypeStr := "connection_confirmation"
		msgStr       := "client has successfully connected to a SSE stream"
		dataMap      := map[string]interface{}{}

		SendEvent(eventsIDstr,
			eventTypeStr, // pTypeStr
			msgStr,       // pMsgStr
			dataMap,
			pEventsCtx,
			pRuntimeSys)

		//-------------

		for ;; {

			eventMsg, moreBool := <- eventsConsumerCh

			// channel is not closed, and there are more messages to be received/processed
			if moreBool {
				streamMsg(eventMsg, p_resp, pRuntimeSys)
				flusher.Flush()
			} else {
				// send this last received message
				streamMsg(eventMsg, p_resp, pRuntimeSys)
				flusher.Flush()
				break
			}
		}
	})
}

//-------------------------------------------------

func streamMsg(p_event_msg EventMsg,
	p_resp      http.ResponseWriter,
	pRuntimeSys *gf_core.RuntimeSys) {


	unix_f       := float64(time.Now().UnixNano())/1000000000.0
	event_id_str := fmt.Sprint(unix_f)
	fmt.Fprintf(p_resp, "id: %s\n", event_id_str)

	event_msg_lst,_ := json.Marshal(p_event_msg)
	fmt.Fprintf(p_resp, "data: %s\n\n", event_msg_lst)
}