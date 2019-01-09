///<reference path="../d/jquery.d.ts" />

namespace gf_posts_browser {
//-----------------------------------------------------
$(document).ready(()=>{
    //-------------------------------------------------
    function log_fun(p_g,p_m) {
        var msg_str = p_g+':'+p_m
        //chrome.extension.getBackgroundPage().console.log(msg_str);

        switch (p_g) {
            case "INFO":
                console.log("%cINFO"+":"+"%c"+p_m,"color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
                break;
            case "FUN_ENTER":
                console.log("%cFUN_ENTER"+":"+"%c"+p_m,"color:yellow; background-color:lightgray","background-color:lightgray");
                break;
        }
    }
    //-------------------------------------------------

    gf_posts_browser.init(log_fun);
});
//-----------------------------------------------------
export function init(p_log_fun) {
    p_log_fun('FUN_ENTER','gf_posts_browser.init()');

    //this app assumes that the first page of the posts is present 
    //in the dom on app startup... subsequent page loads happen from the server
    const static_posts_infos_lst :Object[] = load_data_from_dom(p_log_fun);

    gf_sys_panel.init(p_log_fun);

    gf_posts_browser_view.init(static_posts_infos_lst, p_log_fun);
}
//-----------------------------------------------------
//DATA LOADING
//-----------------------------------------------------
function load_data_from_dom(p_log_fun) {
    p_log_fun('FUN_ENTER','gf_posts_browser.load_data_from_dom()');
        
    const page_posts_infos_lst :Object[] = [];

    $('body #gf_posts_container').find('.gf_post').each((p_i,p_post)=>{

        const post_title_str :string = $(p_post).find('.post_title').text().trim();
        const post_url_str   :string = '/posts/'+post_title_str;

        //---------------------
        //TAGS
        const tags_infos_lst :Object[] = [];
        $(p_post).find('.gf_post_tag').each((p_i,p_tag_element)=>{
            const tag_str     :string = $(p_tag_element).text();
            const tag_url_str :string = $(p_tag_element).attr('href');

            const tag_info_map = {
                'tag_str':    tag_str,
                'tag_url_str':tag_url_str
            };

            tags_infos_lst.push(tag_info_map);
        });
        //--------------------
        //THUMBNAIL URL's
        var thumbnail_url_str :string = $(p_post).find('img').attr('src');
        if (thumbnail_url_str == '' || thumbnail_url_str == 'error') thumbnail_url_str = null;
        //--------------------

        const post_info_map = {
            'post':             p_post,
            'post_title_str':   post_title_str, 
            'post_url_str':     post_url_str,
            'tags_infos_lst':   tags_infos_lst,
            'thumbnail_url_str':thumbnail_url_str
        };

        page_posts_infos_lst.push(post_info_map);
    });

    return page_posts_infos_lst;
}
//-----------------------------------------------------
}