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

type Events__register_producer_msg struct {
	Events_id_str string
}

type Events__register_consumer_msg struct {
	Events_id_str string

	// IMPORTANT!! - channel on which the events_consumer is expecting to receive new 
	//               events produced by the events producer,once consumers registration 
	//               (by processing the Events__register_consumer_msg message) is complete.
	Response_ch   chan chan Event__msg 
}

type Event__msg struct {
	Events_id_str string                 `json:"events_id_str"`
	Type_str      string                 `json:"type_str"`
	Msg_str       string                 `json:"msg_str"`
	Data_map      map[string]interface{} `json:"meta_map"`
}

type Events_ctx struct {
	Register_producer_ch chan Events__register_producer_msg
	Register_consumer_ch chan Events__register_consumer_msg
	Events_broker_ch     chan Event__msg
}

//-------------------------------------------------
func Events__send_event(p_events_id_str string,
	p_type_str    string,
	p_msg_str     string,
	p_data_map    map[string]interface{},
	p_events_ctx  *Events_ctx,
	p_runtime_sys *gf_core.RuntimeSys) {

	e := Event__msg{
		Events_id_str: p_events_id_str,
		Type_str:      p_type_str,
		Msg_str:       p_msg_str,
		Data_map:      p_data_map,
	}
	p_events_ctx.Events_broker_ch <- e
}

//-------------------------------------------------
func Events__register_producer(p_events_id_str string,
	p_events_ctx  *Events_ctx,
	p_runtime_sys *gf_core.RuntimeSys) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_ext_events.Events__register_producer()")

	register_producer_msg := Events__register_producer_msg{
		Events_id_str: p_events_id_str,
	}

	p_events_ctx.Register_producer_ch <- register_producer_msg
}

//-------------------------------------------------
func Events__init(p_sse_url_str string, p_runtime_sys *gf_core.RuntimeSys) *Events_ctx {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_ext_events.Events__init()")

	// yellow := color.New(color.FgYellow).SprintFunc()
	// black  := color.New(color.FgBlack).Add(color.BgYellow).SprintFunc()

	register_producer_ch := make(chan Events__register_producer_msg, 50)
	register_consumer_ch := make(chan Events__register_consumer_msg, 50)
	events_broker_ch     := make(chan Event__msg,                    500)

	events_consumers_map := map[string][]chan Event__msg{}
	go func() {
		for ;; {

			select {

				//-----------------
				// REGISTER EVENTS_PRODUCER
				case register_producer_msg := <- register_producer_ch:
					events_id_str                      := register_producer_msg.Events_id_str
					events_consumers_map[events_id_str] = make([]chan Event__msg,0)

				//-----------------
				// REGISTER EVENTS_CONSUMER
				case register_consumer_msg := <- register_consumer_ch:
					events_id_str                      := register_consumer_msg.Events_id_str
					consumer_ch                        := make(chan Event__msg,50)
					events_consumers_map[events_id_str] = append(events_consumers_map[events_id_str], consumer_ch)
				
					register_consumer_msg.Response_ch <- consumer_ch

				//-----------------
				// EVENT_MSG RELAY
				case event_msg := <- events_broker_ch:
					events_id_str := event_msg.Events_id_str

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

	ctx := &Events_ctx{
		Register_producer_ch: register_producer_ch,
		Register_consumer_ch: register_consumer_ch,
		Events_broker_ch:     events_broker_ch,
	}

	events__init_handlers(p_sse_url_str,
		register_consumer_ch,
		ctx,
		p_runtime_sys)
	return ctx
}

//-------------------------------------------------
func events__init_handlers(p_sse_url_str string,
	p_register_consumer_ch chan<- Events__register_consumer_msg,
	p_events_ctx           *Events_ctx,
	p_runtime_sys          *gf_core.RuntimeSys) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_ext_events.events__init_handlers()")

	// yellow := color.New(color.FgYellow).SprintFunc()
	// black  := color.New(color.FgBlack).Add(color.BgYellow).SprintFunc()


	// IMPORTANT!! - new event_consumers (clients) register via this HTTP handler
	http.HandleFunc(p_sse_url_str, func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.LogFun("INFO", "INCOMING HTTP REQUEST -- "+p_sse_url_str+" ----------")

		// start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

		events_id_str := p_req.URL.Query()["events_id"][0]
		p_runtime_sys.LogFun("INFO", "events_id_str - "+events_id_str)

		register_consumer__response_ch := make(chan chan Event__msg)
		register_consumer_msg          := Events__register_consumer_msg{
			Events_id_str: events_id_str,
			Response_ch:   register_consumer__response_ch,
		}

		p_register_consumer_ch <- register_consumer_msg
		events_consumer_ch := <- register_consumer__response_ch

		flusher,gf_err := gf_core.HTTPinitSSE(p_resp, p_runtime_sys)
		if gf_err != nil {
			return
		}

		//-------------
		// SEND_EVENT

		event_type_str := "connection_confirmation"
		msg_str        := "client has successfully connected to a SSE stream"
		data_map       := map[string]interface{}{}

		Events__send_event(events_id_str,
			event_type_str, // p_type_str
			msg_str,        // p_msg_str
			data_map,
			p_events_ctx,
			p_runtime_sys)

		//-------------

		for ;; {

			event_msg,more_bool := <- events_consumer_ch
			// pLogFun("INFO",black("EVENTS >> EVENTS_CONSUMER <- EVENTS_BROKER msg")+" > "+yellow(event_msg.Type_str))

			// channel is not closed, and there are more messages to be received/processed
			if more_bool {
				events__stream_msg(event_msg, p_resp, p_runtime_sys)
				flusher.Flush()
			} else {
				// send this last received message
				events__stream_msg(event_msg, p_resp, p_runtime_sys)
				flusher.Flush()
				break
			}
		}

		//end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

		/*//FIX!! - gf_rpc_lib imports gf_core, which causes a import cycle,
		//        and fails compilation.
		go func() {
			gf_rpc_lib.Store_rpc_handler_run(p_sse_url_str,
								start_time__unix_f,
								end_time__unix_f,
								p_mongodb_coll,
								pLogFun)
		}()*/
	})
}

//-------------------------------------------------
func events__stream_msg(p_event_msg Event__msg,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.RuntimeSys) {


	unix_f       := float64(time.Now().UnixNano())/1000000000.0
	event_id_str := fmt.Sprint(unix_f)
	fmt.Fprintf(p_resp, "id: %s\n",event_id_str)

	event_msg_lst,_ := json.Marshal(p_event_msg)
	fmt.Fprintf(p_resp, "data: %s\n\n", event_msg_lst)
}