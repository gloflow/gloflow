[
    ["lang_v", "0.0.6"],
    ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

    [   
        ["y", 1.5, "x", 2],

        //------------------------------
        // "$r" - some random variable name in which to place results of the rpc_call rule
        ["$r", ["rpc_call", [
                "gf",   // node_name 
                "test", // module_name
                "return_simple_list", // function_name

                // args
                [
                    ["list_len_to_return_int", 10]
                ]
            ]]
        ],

        //------------------------------
        // draw as many spheres as there are elements in the response r.data list
        ["*", ["len", ["$r.data"]], "sphere"]

        //------------------------------
    ]
],


[
    ["lang_v", "0.0.6"],

    ["rpc_serve", [
        "test_node", // node_name

        // handlers
        ["handlers", [
            [
                "/v1/rpc/return_simple_list",
                "test",               // module_name
                "return_simple_list", // function_name
            
                // args_spec
                ["args_spec", [
                    "list_len_to_return_int"
                ]],

                // code
                ["code", [

                    // create variable, from RPC call input argument.
                    // RPC args are declared in "args" section
                    ["$list_len_to_return_lst", "$rpc_in.list_len_to_return_int"],

                    // draw N spheres
                    ["*", "$list_len_to_return_lst", "sphere"],

                    // "make() - creates a new resource, a "list" in this instance, of a certain length
                    // return a new list
                    ["$new_list", ["make", ["list", "$list_len_to_return_lst"]]],
                    ["return", "$new_list"]
                ]]
            ]
        ]]
    ]]
]