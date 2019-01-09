
///<reference path="../d/jquery.d.ts" />

namespace gf_crawl_events {

declare var EventSource;

//---------------------------------------------------
export function init_SSE(p_log_fun) {
    p_log_fun("FUN_ENTER","gf_crawl_events.init_SSE()");



    const events_id_str = "crawler_events";
    const event_source  = new EventSource("/a/crawl/events?events_id="+events_id_str)

    $('#view_crawl_events_btn').on('click',(p_e)=>{

        console.log("REGISTER EVENT_SOURCE")
        //const events_id_str = "crawler_events";
        //const event_source  = new EventSource("/a/crawl/events?events_id="+events_id_str)


        event_source.onopen = (p_e)=>{
            console.log('EventSource >> OPEN CONN');
        }

        var i=0;
        event_source.onmessage = (p_e)=>{

            console.log('>>>>> MESSAGE');
            const event_data_map = JSON.parse(p_e.data);
                
            console.log(event_data_map)
            
            view_server_event(event_data_map,
                        p_log_fun);

            i+=1;
        }

        event_source.onerror = (p_e)=>{


            console.log('EventSource >> ERROR - '+event_source.readyState);
            console.log(EventSource.CLOSED)
            console.log(p_e);
              
            //connection was closed
            if (event_source.readyState == EventSource.CLOSED) {
                console.log("EVENT_SOURCE CLOSED")
            }
        }
    });
}
//---------------------------------------------------
function view_server_event(p_event_data_map,
                    p_log_fun) {
    p_log_fun("FUN_ENTER","gf_crawl_events.view_server_event()");



    const event_type_str = p_event_data_map['event_type_str'];


    switch (event_type_str) {
        //--------------
        case 'fetch__http_request__done':
            break;
        //--------------
        case 'image_download__http_request__done':
            break;
        //--------------
        default:
            break;
    }
}
//-------------------------------------------------
}