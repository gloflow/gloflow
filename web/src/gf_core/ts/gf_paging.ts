









//---------------------------------------------------
export function init(p_initial_page_int :number,
    p_pages_number_int :number,
    p_page_load_fun :Function,
    p_log_fun       :Function) {

    var current_page_int = p_initial_page_int;
    var page_is_loading_bool = false;
    
    //---------------------------------------------------
    const scroll_handler_fun = async ()=>{

        // $(document).height() - height of the HTML document
        // window.innerHeight   - Height (in pixels) of the browser window viewport including, if rendered, the horizontal scrollbar
        if (window.scrollY >= $(document).height() - (window.innerHeight+50)) {
            
            // IMPORTANT!! - only load 1 page at a time
            if (!page_is_loading_bool) {
                
                page_is_loading_bool = true;
                p_log_fun("INFO", `current_page_int - ${current_page_int}`);

                await p_page_load_fun(current_page_int);
                
                current_page_int += p_pages_number_int;
                page_is_loading_bool = false;
            }
        }
    };

    //---------------------------------------------------
    window.onscroll = scroll_handler_fun;
    return scroll_handler_fun;
}