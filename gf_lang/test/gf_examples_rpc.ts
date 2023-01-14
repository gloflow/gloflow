/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

export function get() {

    // 0.0.5 - compiler and interpreter now run in WASM context
    //       - half of test programs run succssfuly

    // 0.0.6 - introduction of remote-procedure-calls
    // - adding a system rule that makes a GF standardized remote procedure call.
    //      - basic model of Erlang addressing "node/module/function" added.
    
    const rpc_simple_lst = [
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
                    ],
                ]],
            ],

            //------------------------------
            // draw as many spheres as there are elements in the response r.data list
            ["*", ["len", "$r.data"], "sphere"]

            //------------------------------
        ]
    ];

    
    //------------------------------
    // ADD!! - variable support, for creating general user-defined variables.

    //------------------------------


    // ultimatelly rpc_server should be able to run on the browser as well,
    // and get a listening "port" allocated on the GF server so that it can route 
    // public web traffic to it.
    const rpc_server_lst = [
        ["lang_v", "0.0.6"],

        ["rpc_serve", [
                "test_node", // node_name

                // handlers
                ["handlers", [[
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
                ]]]
        ]]
    ];

    return {
        "rpc_simple_lst": rpc_simple_lst,
        "rpc_server_lst": rpc_server_lst,
    }
}