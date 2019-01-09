

namespace gf_crawl_stats__fetches {

declare var c3;

//---------------------------------------------------
export function view__fetches_per_day(p_fetches_by_days_map,
                                    p_parent,
                                    p_log_fun) {
    p_log_fun('FUN_ENTER','gf_crawl_stats__fetches.view__fetches_per_day()');

    const plot = $(`
            <div id='plots'>
                <div id='new_fetches_per_day__plot'>
                    <svg width='2000' height='2000'></svg>
                </div>
            </div>`);

    $(p_parent).append(plot);

    const counts_by_day__sorted_lst        = p_fetches_by_days_map['counts_by_day__sorted_lst'];
    const domain_counts_by_day__sorted_lst = p_fetches_by_days_map['domain_counts_by_day__sorted_lst'];

    /*//------------------
    const daily_count_lst            = [];
    const daily_count_per_domain_map = {}; //:Map<:List<:Int>>

    for (var day_stat_map of p_stats_lst) {

    	const count_int                   = day_stat_map['total_count_int'];
        const total_count__per_domain_map = day_stat_map['total_count__per_domain_map'];


    	daily_count_lst.push(count_int);

        //accumulate up domain counts per day into a specific domains counts lists,
        for (var domain_str in total_count__per_domain_map) {

            const domain_count_int = total_count__per_domain_map[domain_str];

            if (domain_str in daily_count_per_domain_map) {
                daily_count_per_domain_map[domain_str].push(domain_count_int);
            } else {
                daily_count_per_domain_map[domain_str] = [domain_count_int];
            }
        }
    }*/
    //------------------
    //C3_COLUMNS

    counts_by_day__sorted_lst.unshift('per-day fetches count');

    const c3_columns_lst = [
        counts_by_day__sorted_lst
    ]

    for (var domain_counts_map of domain_counts_by_day__sorted_lst) {

        const domain_str      = domain_counts_map['domain_str'];
        const days_counts_lst = domain_counts_map['days_counts_lst'];
        //const domain_counts_lst = daily_count_per_domain_map[domain_str];

        //create a title of the columns for C3
        days_counts_lst.unshift(domain_str);

        c3_columns_lst.push(days_counts_lst);
    }
    //------------------

    const top_c3_columns_lst = c3_columns_lst.slice(0,20);

    console.log(top_c3_columns_lst)

    const chart = c3.generate({
        bindto:'#new_fetches_per_day__plot',
        data  : {
            columns:top_c3_columns_lst
            /*columns: [
            daily_total_count_lst,
            //['data1', 30, 200, 100, 400, 150, 250],
            //['data2', 50, 20, 10, 40, 15, 25]
            ]*/
        }
    });
}
//---------------------------------------------------
}