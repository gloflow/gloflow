


















export function get() {

    const first_scene__program_ast_lst = [
        ["lang_v", 0.3],
        ["set", "color-background", ["rgb", 0.5, 0.5, 0.5]],
        ["set", "material", ["wireframe", true]],

        [["set", "color", ["rgb", 0, 0, 1.0]], "cube"],                                   // cube
        [["set", "color", ["rgb", 0, 0, 1.0]], ["x", 2.0], "cube"],                       // {x 2} cube
        [["set", "color", ["rgb", 0, 0, 1.0]], ["x", 4.0], "cube"],                       // {x 4} cube
        [["set", "color", ["rgb", 0, 0, 1.0]], ["x", 6.0], "cube"],                       // {x 6} cube
        [["set", "color", ["rgb", 0, 0, 1.0]], ["x", 9.0, "z", 2.0], ["y", 3.0], "cube"], // {x 7 z 1} {y 3} cube
        
        [
            ["set", "color", ["rgb", 0, 0, 1.0]],
            ["*", 3, [["y", -2.0, "z", -2.0, "x", -2.0], "cube"]],  // 3 * {x -2 z -2 x -2} cube
        ],
        
        [
            ["set", "color", ["rgb", 0, 0, 1.0]],
            ["*", 3, [["x", 5.0], ["*", 6, [["y", 2.0], "cube"]]]], // 3 * {x 5} 6 * {y 2} cube
        ],

        /*
        transforms to:
        [

            ["x", 5.0], [
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
            ]
            
            ["x", 5.0], [
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
            ]

            ["x", 5.0], [
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
                [["y", 2.0], "cube"]
            ]
        ]
        */

        [
            ["set", "color", ["rgb", 0, 0, 1.0]],
            ["*", 3, [["x", 5.4], ["*", 1, [["z", 1.6], "cube"]]]], // 3 * {x 5} 6 * {y 2} cube
        ],
        [
            ["set", "color", ["rgb", 0, 0, 1.0]],
            ["*", 2, [["x", 2.1], ["*", 2, [["y", 2.1], "cube"]]]], // 3 * {x 5} 6 * {y 2} cube
        ],
        [
            ["set", "color", ["rgb", 0, 0, 1.0]],
            ["*", 6, [["x", 4.4], ["*", 1, [["z", 2.4], "cube"]]]], // 3 * {x 5} 6 * {y 2} cube
        ],
        [
            ["set", "color", ["rgb", 0, 0, 1.0]],
            ["*", 4, [["x", 5.2], ["*", 1, [["y", 2.2], "cube"]]]], // 3 * {x 5} 6 * {y 2} cube
        ],
        [
            ["set", "color", ["rgb", 0, 0, 1.0]],
            ["*", 1, [["x", 6.1], ["*", 3, [["z", 2.0], "cube"]]]], // 3 * {x 5} 6 * {y 2} cube
        ],
        [
            ["set", "color", ["rgb", 0, 0, 1.0]],
            ["*", 1, [["x", 1.4], ["*", 2, [["y", 2.1], "cube"]]]], // 3 * {x 5} 6 * {y 2} cube
        ]
    ];



    const rules_test__program_ast_lst = [
        ["lang_v", 0.3],
        ["set", "color-background", ["rgb", 0.9, 0.9, 0.9]],
        ["set", "material", ["wireframe", true]],
        
        [["set", "color", ["rgb", 0, 0, 1.0]], "R1"],
        [["set", "color", ["rgb", 0, 0.2, 1.0]], "R2"],
        [["set", "color", ["rgb", 0, 0.4, 1.0]], "R3"],
        [["set", "color", ["rgb", 0, 0.4, 1.0]], "R4"],
        
        ["rule", "R1", [
            ["*", 2.0, [["x", 1.0, "y", 1.5, "rz", 0.05], "cube"]],
            ["R1"]
        ]],
        ["rule", "R2", [
            ["*", 3.0, [["x", 2.0, "y", 2.5, "z", 0.5, "rz", 0.1], "cube"]],
            ["R2"]
        ]],
        ["rule", "R3", [
            ["*", 4.0, [["x", 0.5, "z", 0.5, "rz", 0.2], "cube"]],
            ["R3"]
        ]],
        ["rule", "R4", [
            [["x", 1.0, "z", 1.0, "y", 0.5, "rz", 0.4], "cube"],
            ["R4"]
        ]]
    ];


    const multi_rule_test__program_ast_lst = [
        ["lang_v", 0.3],
        ["set", "color-background", ["rgb", 0.9, 0.9, 0.9]],
        ["set", "material", ["wireframe", true]],

        [["set", "color", ["rgb", 0, 0, 1.0]], "R1"],
        [["set", "color", ["rgb", 0, 0.06, 1.0]], ["x", 2.0], "R1"],
        [["set", "color", ["rgb", 0, 0.11, 1.0]], ["x", 4.0,  "y", 1.0], "R1"],
        [["set", "color", ["rgb", 0, 0.16, 1.0]], ["x", 6.0,  "y", 2.0], "R1"],
        [["set", "color", ["rgb", 0, 0.22, 1.0]], ["x", 8.0,  "y", 3.0], "R1"],
        [["set", "color", ["rgb", 0, 0.25, 1.0]], ["x", 10.0, "y", 4.0], "R1"],
        [["set", "color", ["rgb", 0, 0.32, 1.0]], ["x", 12.0, "y", 5.0], "R1"],
        [["set", "color", ["rgb", 0, 0.37, 1.0]], ["x", 14.0, "y", 6.0], "R1"],
        [["set", "color", ["rgb", 0, 0.42, 1.0]], ["x", 16.0, "y", 7.0], "R1"],
        [["set", "color", ["rgb", 0, 0.48, 1.0]], ["x", 18.0, "y", 8.0], "R1"],
        [["set", "color", ["rgb", 0, 0.53, 1.0]], ["x", 20.0, "y", 9.0], "R1"],
        [["set", "color", ["rgb", 0, 0.57, 1.0]], ["x", 21.0, "y", 10.0], "R1"],
        [["set", "color", ["rgb", 0, 0.62, 1.0]], ["x", 22.0, "y", 11.0], "R1"],
        [["set", "color", ["rgb", 0, 0.67, 1.0]], ["x", 23.0, "y", 12.0], "R1"],
        [["set", "color", ["rgb", 0, 0.75, 1.0]], ["x", 24.0, "y", 13.0], "R1"],

        ["rule", "R1", [
            ["*", 2.0, [["x", 0.5, "y", 0.7, "rz", 0.07], "cube"]],
            ["R1"]
        ]],
        ["rule", "R1", [
            ["*", 1.0, [["x", 0.5, "y", 0.7, "z", 1.0, "rz", 0.05], "cube"]],
            ["R1"]
        ]],
        ["rule", "R1", [
            ["*", 1.0, [["x", 0.5, "y", 0.7, "z", 1.3, "rz", 0.05], "cube"]],
            ["R1"]
        ]],
        ["rule", "R1", [
            ["*", 1.0, [["x", 0.5, "y", 0.7, "z", 1.5, "rz", 0.05], "cube"]],
            ["R1"]
        ]],
    ];

    const color_test__program_ast_lst = [
        ["lang_v", "0.0.3"],
        ["set", "color-background", ["rgb", 0.1, 0.6, 0.9]],
        ["set", "material", ["wireframe", true]],

        [
            // color has to be set as a sub-expression of some root exxpression
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", -9], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", -7], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", -5], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", -3], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", -1], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", 1], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", 3], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", 5], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", 7], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", 9], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", 12], "R1"],
        ],
        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],  // set color {r 0.4 g 0.8 b 0.5}
            [["z", 14], "R1"],
        ],



        ["rule", "R1", [
            ["set", "iters_max", 500],
            ["*", 1, [["x", 1.1, "y", -0.5, "z", 1, "cg", 0.005], "cube"]], // 50 * {x 2 cr 0.05} cube
            ["R1"]
        ]],
        ["rule", "R1", [
            ["set", "iters_max", 500],
            ["*", 1, [["x", 1.2, "y", 0.9, "z", -1, "cg", 0.005], "cube"]], // 50 * {x 2 cr 0.05} cube
            ["R1"]
        ]],
        ["rule", "R1", [
            ["set", "iters_max", 500],
            ["*", 1, [["x", 1.3, "y", -0.9, "z", -1, "cg", 0.005], "cube"]], // 50 * {x 2 cr 0.05} cube
            ["R1"]
        ]],
        ["rule", "R1", [
            ["set", "iters_max", 500],
            ["*", 1, [["x", 1.4, "y", 0.5, "z", 1, "cg", 0.005], "cube"]], // 50 * {x 2 cr 0.05} cube
            ["R1"]
        ]],
    ];



    const conditional_test__program_ast_lst = [
        // 0.0.2 - lang added support for:
        //  - conditional statements
        //      - only 'if' for now
        //      - only single logical expression, not composite
        //      - support for logic-operators (>, <, ==, !=, etc.)
        //  - basic system variables reading ($i - rule iteration counter)
        //  - proper state independence between rules (rule state isolation).
        //    rules state mutation does not propagate to their calling rules.
        //  - rule-modifiers
        //      - iters_max - for specifying how many iterations a rule can run for
        //  - 3d obj scale incremental modifier and scale setter
        //  - hex-color parsing support
        //  - "print" function (for debugging)

        ["lang_v", "0.0.3"], 
        ["set", "color-background", ["rgb", 0.7, 0.7, 0.7]],

        [
            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],
            ["set", "material", ["wireframe", true]],
            ["set", "iters_max", 30],
            ["R1"],
        ],

        //-------------------------------------------------
        // R1
        ["rule", "R1", [["iters_max", 10]], [

            // branch off and run R2
            ["print", ["first condition", "$i"]],
            ["if", ["==", "$i", 2], [
                ["print", ["$i", "running R2 A"]],
                ["R2"]
            ]],

            // branch off and run R2
            ["print", ["second condition", "$i"]],
            ["if", ["==", "$i", 4], [
                ["print", ["$i", "running R2 B"]],
                ["R2"]
            ]],

            // branch off and run R2
            ["print", ["third condition", "$i"]],
            ["if", ["==", "$i", 8], [
                ["print", ["$i", "running R2 C"]],
                ["R2"]
            ]],

            ["set", "color", ["rgb", 0.7, 0.0, 0.0]],
            ["*", 1, [["x", 2], "cube"]],

            ["R1"]
        ]],

        //-------------------------------------------------
        // R2
        ["rule", "R2", [["iters_max", 5]], [
            ["print", ["===============in R2", "$i"]],

            ["y", 4],
            
            ["print", ["first R2 condition", "$i"]],
            ["if", [">", "$i", 2], [                
                ["set", "color", ["rgb", 0.0, 0.9, 0.0]],
                ["print", ["$i", "running R3"]],
                
                ["R3"],
            ]],

            [
                ["print", ["===============in R2 cube", "$i"]],
                ["set", "color", ["rgb", 0.0, 0.0, 0.9]],
                ["*", 1, ["cube"]],

                ["print", ["===============exit R2", "$i"]],
                ["R2"],
            ],

            
        ]],

        //-------------------------------------------------
        // R3
        ["rule", "R3", [["iters_max", 5]], [
            ["print", ["$i", "===============in R3"]],

            ["*", 1, [["z", 4], "cube"]],
            ["print", ["$i", "===============exit R3"]],
            "R3"
        ]]

        //-------------------------------------------------
    ];




    const main_branches_color_1_str = "#334ea2";
    const main_branches_color_2_str = "#4b69c6";
    const main_branches_color_3_str = "#6881cf";
    const secondary_branches_color_1_str = "#118984";
    const secondary_branches_color_2_str = "#79bd5e";
    const offshots_color_1_str = "#ed3f21";
    const offshots_color_2_str = "#e38623";

    const conditional_test_2__program_ast_lst = [
        ["lang_v", "0.0.3"],
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

        [
            ["set", "color", main_branches_color_1_str], // ["rgb", 0.7, 0.0, 0.0]],
            ["set", "material", ["wireframe", true]],
            ["set", "iters_max", 250],

            ["set", "scale", [1.7, 1.7, 1.7]],  
            [["x", 0], "R1"],
            [["x", 2], "R1"],
            [["x", 4], "R1"],
            [["x", 6], "R1"],
            [["x", 8], "R1"],
            [["x", 10],"R1"],
            [["x", 12], "R1"],
            [["x", 14], "R1"],
            [["x", 16], "R1"],
            [["x", 18], "R1"],
            [["x", 20], "R1"],
            [["x", 22],"R1"],
        ],

        //-------------------------------------------------
        // R1
        ["rule", "R1", [["iters_max", 60]], [
            ["set", "color", main_branches_color_1_str], // ["rgb", 0.7, 0.0, 0.0]],
            [["y", 2.4, "x", 0.4, "z", 0.4], "cube"],
            ["R1"]
        ]],

        ["rule", "R1", [["iters_max", 60]], [
            ["set", "color", main_branches_color_2_str], // ["rgb", 0.8, 0.0, 0.0]],
            [["y", 2.4, "x", -0.4, "z", -0.4], "cube"],
            ["R1"]
        ]],
        
        ["rule", "R1", [["iters_max", 60]], [
            ["set", "color", main_branches_color_3_str], // ["rgb", 0.8, 0.0, 0.0]],
            [["y", 2.4, "x", -0.4, "z", -0.4], "cube"],
            ["R1"]
        ]],

        ["rule", "R1", [["iters_max", 60]], [
            ["set", "color", main_branches_color_2_str], // ["rgb", 0.8, 0.0, 0.0]],
            [["y", 2.4, "x", -0.4, "z", -0.4], "cube"],
            ["R1"]
        ]],

        ["rule", "R1", [["iters_max", 60]], [
            ["set", "color", main_branches_color_3_str], // ["rgb", 0.9, 0.0, 0.0]],

            ["if", [">", "$i", 8], [
                ["R2"],

                // so that the cube coming after this gets drawn larger
                // only when R2 rule is activated
                ["set", "scale", [3.0, 3.0, 3.0]],  
            ]],

            [["y", 2.4, "x", -0.4, "z", 0.4], "cube"],

            ["if", [">", "$i", 8], [   
                // return back to normal cube scale for R1
                ["set", "scale", [1.7, 1.7, 1.7]],  
            ]],

            ["R1"]
        ]],

        ["rule", "R1", [["iters_max", 60]], [
            ["set", "color", main_branches_color_1_str], // ["rgb", 0.6, 0.0, 0.0]],

            ["if", [">", "$i", 13], [        
                ["R2"],

                // so that the cube coming after this gets drawn larger
                // only when R2 rule is activated
                ["set", "scale", [3.0, 3.0, 3.0]],  
            ]],

            [["y", 2.4, "x", 0.4, "z", -0.4], "cube"],

            ["if", [">", "$i", 13], [   
                // return back to normal cube scale for R1
                ["set", "scale", [1.7, 1.7, 1.7]],  
            ]],

            ["R1"]
        ]],

        //-------------------------------------------------
        // R2
        ["rule", "R2", [["iters_max", 30]], [
            ["set", "color", secondary_branches_color_1_str], // ["rgb", 0.0, 0.7, 0.0]],
            ["set", "scale", [1.0, 1.0, 1.0]],


            ["if", [">", "$i", 15], [
                ["R4"],
            ]],


            [["x", 0.4, "y", 0.2, "z", 1.2, "cg", 0.05], "cube"],
            ["R2"]
        ]],

        ["rule", "R2", [["iters_max", 30]], [
            ["set", "color", secondary_branches_color_2_str], // ["rgb", 0.0, 0.8, 0.0]],
            ["set", "scale", [1.0, 1.0, 1.0]],  

            ["if", ["==", "$i", 29], [
                ["set", "scale", [0.7, 0.7, 0.7]],
                ["set", "color", offshots_color_1_str], // ["rgb", 0.9, 0.0, 0.9]],              
                ["R3"],
                ["R3"],
                ["R3"],
                ["set", "scale", [1.0, 1.0, 1.0]],
                ["set", "color", secondary_branches_color_2_str], // ["rgb", 0.0, 0.8, 0.0]],
            ]],

            [["x", -0.4, "y", 0.2, "z", 1.2, "cg", 0.05], "cube"],
            ["R2"]
        ]],

        ["rule", "R2", [["iters_max", 30]], [
            ["set", "color", secondary_branches_color_2_str], // ["rgb", 0.0, 0.8, 0.0]],
            ["set", "scale", [1.0, 1.0, 1.0]],  

            ["if", ["==", "$i", 29], [
                ["set", "scale", [0.7, 0.7, 0.7]],
                ["set", "color", offshots_color_2_str], // ["rgb", 0.9, 0.0, 0.9]],          
                ["R3"],
                ["R3"],
                ["set", "scale", [1.0, 1.0, 1.0]],
                ["set", "color", secondary_branches_color_2_str], // ["rgb", 0.0, 0.8, 0.0]],
            ]],

            [["x", 0.4, "y", -0.2, "z", 1.2, "cg", 0.05], "cube"],
            ["R2"]
        ]],

        ["rule", "R2", [["iters_max", 30]], [
            ["set", "color", secondary_branches_color_2_str], // ["rgb", 0.0, 0.8, 0.0]],
            ["set", "scale", [1.0, 1.0, 1.0]],  

            [["x", -0.4, "y", -0.2, "z", 1.2, "cg", 0.05], "cube"],
            ["R2"]
        ]],

        //-------------------------------------------------
        ["rule", "R3", [["iters_max", 20]], [
            
            [["x", 0.3, "y", 0.2, "z", 0.5, "cr", -0.02], ["sx", -0.03, "sy", -0.03, "sz", -0.03], "cube"],
            ["R3"]
        ]],

        ["rule", "R3", [["iters_max", 20]], [
            [["x", -0.3, "y", -0.2, "z", 0.5, "cr", -0.02], ["sx", -0.03, "sy", -0.03, "sz", -0.03], "cube"],
            ["R3"]
        ]],

        ["rule", "R3", [["iters_max", 20]], [
            [["x", 0.3, "y", -0.2, "z", 0.5, "cr", -0.02], ["sx", -0.03, "sy", -0.03, "sz", -0.03], "cube"],
            ["R3"]
        ]],

        ["rule", "R3", [["iters_max", 20]], [
            [["x", -0.3, "y", 0.2, "z", 0.5, "cr", -0.02], ["sx", -0.03, "sy", -0.03, "sz", -0.03], "cube"],
            ["R3"]
        ]],

        //-------------------------------------------------
        ["rule", "R4", [["iters_max", 3]], [
            ["set", "scale", [0.3, 0.3, 0.3]], 
            ["set", "color", secondary_branches_color_2_str], // ["rgb", 0.0, 0.3, 0.0]],
            [["y", 0.3], "cube"],
            ["R4"]
        ]],

        //-------------------------------------------------
    ];






    const mothership_color_str = "#EB4B0F";
    const mothership_module_color_str = "#DC470F";
    const mothership_2_color_str = "#E1761F";
    const connection_uncle_color_str = "#F28D68";
    const main_block_color_str = "#4176B6";
    const small_block_color_str = "#355C8E";
    const small_block_2_color_str = "#29405E";

    const form_experiment__ast_lst = [
        ["lang_v", "0.0.3"],
        ["set", "color-background", ["rgb", 0.70, 0.70, 0.70]],
        
        [
            ["set", "material", ["wireframe", true]],
            ["set", "color", mothership_color_str],
            ["set", "scale", [100, 100, 100]], 
            ["x", -55, "cube"],

            ["set", "scale", [80, 80, 80]],
            ["set", "color", mothership_module_color_str], 
            [["y", 12], "mothership_module"],

            ["set", "material", ["wireframe", true]],
            ["mothership_module_connection"]
        ],
        [
            ["set", "color", mothership_2_color_str], 
            ["set", "scale", [30, 30, 30]],

            ["set", "material", ["wireframe", true]],
            ["x", -15, "cube"]
        ],
       
        ["mothership_module_movement"],

        //-------------------------------------------------
        ["rule", "small_block_base", [["iters_max", 200]], [
            [["x", 0.5, "y", 0.6, "sx", -0.005, "sy", -0.005, "sz", -0.005,
                "cr", 0.005, "cg", 0.005, "cb", 0.005], "cube"],
            ["small_block_base"]
        ]],
        ["rule", "small_block_base", [["iters_max", 200]], [
            [["x", 0.5, "y", -0.6, "sx", -0.005, "sy", -0.005, "sz", -0.005,
                "cr", 0.005, "cg", 0.005, "cb", 0.005], "cube"],
            ["small_block_base"]
        ]],
        ["rule", "small_block_base", [["iters_max", 200]], [
            [["x", 0.5, "y", -0.6, "z", 0.6, "sx", -0.005, "sy", -0.005, "sz", -0.005,
                "cr", 0.005, "cg", 0.005, "cb", 0.005], "cube"],
            ["small_block_base"]
        ]],
        ["rule", "small_block_base", [["iters_max", 200]], [
            [["x", 0.5, "y", 0.6, "z", -0.6, "sx", -0.005, "sy", -0.005, "sz", -0.005,
                "cr", 0.005, "cg", 0.005, "cb", 0.005], "cube"],
            ["small_block_base"]
        ]],

        //-------------------------------------------------
        ["rule", "mothership_module", [["iters_max", 3]], [
            [["y", 6, "sx", -10, "sy", -10, "sz", -10], "cube"],
            ["mothership_module"]
        ]],

        ["rule", "mothership_module_connection", [["iters_max", 100000]], [
            ["set", "scale", [4, 4, 4]], 

            ["if", ["==", "$i", 14], [
                [["z", -10], "mothership_module_connection_uncle"],
                ["uncle_tangent__z_positive"],
                ["z", 10], // revert 
            ]],
            

            ["if", ["==", "$i", 20], [
                [["z", 15], "mothership_module_connection_uncle"],
                ["uncle_tangent__z_negative"],
                ["z", -15], // revert 
            ]],
            

            ["if", ["==", "$i", 10], [
                [["x", 10], "mothership_module_connection_uncle"],
                ["uncle_tangent__z_positive"],
                ["x", -10], // revert 
            ]],
            

            [["y", 7], "cube"],
            ["mothership_module_connection"]
        ]],
        
        ["rule", "mothership_module_connection_uncle", [["iters_max", 60]], [
            ["set", "scale", [2, 2, 2]], 
            ["set", "color", connection_uncle_color_str], 

            [["y", 7], "cube"],
            ["mothership_module_connection_uncle"]
        ]],
        ["rule", "uncle_tangent__z_positive", [["iters_max", 6000]], [
            // ["set", "scale", [0.8, 0.8, 0.8]], 
            [["z", 4], "cube"],
            "uncle_tangent__z_positive"
        ]],
        ["rule", "uncle_tangent__z_negative", [["iters_max", 6000]], [
            // ["set", "scale", [0.8, 0.8, 0.8]], 
            [["z", -4], "cube"],
            "uncle_tangent__z_negative"
        ]],

        //-------------------------------------------------
        ["rule", "mothership_module_movement", [["iters_max", 1]], [
            [
                ["set", "color", main_block_color_str], 
                ["set", "scale", [10, 10, 10]], 
                "cube"
            ],
            ["x", 5.45],
            [
                ["set", "color", small_block_2_color_str], 
                ["set", "scale", [1, 1, 1]], 
                [["y", 3], "cube"],
                ["small_block_base"]
            ],
            [
                ["set", "color", small_block_color_str], 
                ["set", "scale", [1, 1, 1]], 
                [["y", -5, "z", 2], "cube"],
                ["small_block_base"]
            ],
            [
                ["set", "color", small_block_color_str], 
                ["set", "scale", [1, 1, 1]], 
                [["y", -1, "z", -2.2], "cube"],
                ["small_block_base"]
            ],
            [
                ["set", "color", small_block_color_str], 
                ["set", "scale", [1, 1, 1]], 
                [["y", 3, "z", -1], "cube"],
                ["small_block_base"]
            ],
            [
                ["set", "color", small_block_2_color_str], 
                ["set", "scale", [1, 1, 1]], 
                [["y", -3.5, "z", 3], "cube"],
                ["small_block_base"]
            ],
        ]],

        //-------------------------------------------------
    ];



    const extrusion_color_str = "#6881cf";
    const extrusion_color_2_str = "#42b0f5";
    const extrusion_color_3_str = "#f5bf42";
    const extrusion_color_4_str = "#f59842";
    const sphere_test__program_ast_lst = [
        ["lang_v", "0.0.3"],
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],
        [
            ["set", "color", extrusion_color_str],
            ["R1"],
            
        
            ["set", "material", ["wireframe", true]],
            ["set", "color", extrusion_color_2_str],
            ["set", "scale", [1.2, 1.2, 1.2]], 
            [["x", 10], "R1"],
            ["set", "material", ["wireframe", false]],

            
            ["set", "color", extrusion_color_3_str],
            ["set", "scale", [1.3, 1.3, 1.3]], 
            [["x", 20], "R1"],
            

            ["set", "color", extrusion_color_4_str],
            ["set", "scale", [1.4, 1.4, 1.4]], 
            [["x", 40], "R1"],
        ],

        //-------------------------------------------------
        ["rule", "R1", [["iters_max", 200]], [
            [["y", 3, "z", -1.0, "rz", 0.1, "sx", -0.01, "sy", -0.01, "sz", -0.01, "cr", -0.002, "cg", -0.002, "cb", -0.002], "sphere"],
            ["R1"]
        ]],
        ["rule", "R1", [["iters_max", 200]], [
            [["y", 3.5, "z", 1.0, "rz", 0.1, "sx", -0.01, "sy", -0.01, "sz", -0.01, "cr", -0.002, "cg", -0.002, "cb", -0.002], "sphere"],
            ["R1"]
        ]],
        ["rule", "R1", [["iters_max", 200]], [
            [["y", -3.5, "z", -1.0, "rz", 0.1, "sx", -0.01, "sy", -0.01, "sz", -0.01, "cr", -0.002, "cg", -0.002, "cb", -0.002], "sphere"],
            ["R1"]
        ]],

        //-------------------------------------------------
    ];

    const shader_test__program_ast_lst = [
        // 0.0.3 - lang added support for:
        // - shaders
        //      - specify vertex/fragment GLSL shaders along with other program definitions (rules, etc.)
        //      - uniform definitions - declaring uniforms and their types|default values
        //          - for program validation
        //          - mechanism for passing uniform values to shaders from glo-lang code state/variables
        // - material properties
        //      - setting of material properties via state setters
        //          - propagate to all subsequent expressions state
        //          - only disable/enable wireframe property on materials for now.
        //            (other properties for material will be exposed soon).
        // - arbitrary arithmetic (*,/,+,-) expressions:
        //      - two numeric values|variables|sub-expressions can be multiplied now at runtime, in any combination of those.
        //      - arbitrary number of arithmetic operations sub-expressions levels/nesting supported.
        //      - previously only expression multiplication was supported, acting as an compiler expansion macro
        //        that multiplied/cloned a given sub-expression N number of times at compile time.
        // - sphere primitive added
        //      - not just cube anymore
        // - rotation bug fix
        //      - proper world-centered rotation with independent axis-rotation matrices 

        ["lang_v", "0.0.3"], 
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

        [
            ["set", "material", ["wireframe", true]],
            ["set", "color", extrusion_color_2_str],
            ["sphere"],
        ],

        [
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.3]]], // "$i"]]],
            ["x", 30], "sphere"
        ],


        [
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.5]]], // "$i"]]],
            ["sx", -0.2, "sy", -0.2, "sz", -0.2],
            ["x", 40, "rz", 0.4], "sphere"
        ],
        [
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.7]]], // "$i"]]],
            ["sx", -0.4, "sy", -0.4, "sz", -0.4],
            ["x", 50, "rz", 0.6], "sphere"
        ],
        [
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", ["*", 2, 0.35]]]], // "$i"]]],
            ["sx", 0.6, "sy", 0.6, "sz", 0.6],
            ["x", -20, "y", 20, "rx", 1.0], "sphere"
        ],


        [
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.9]]], // "$i"]]],
            "R1"
        ],

        //-------------------------------------------------
        ["rule", "R1", [["iters_max", 5]], [
            ["set", "color", extrusion_color_2_str],
            [["y", 10, 'ry', 0.5], "cube"],
            ["R1"]
        ]],

        //-------------------------------------------------

        ["shader", "gf_shader_test",
            ["uniforms", [
                ["i", "float", 0.1]
            ]],
            ["vertex", `

                varying vec3 fNormal;
                varying vec3 worldPos;
                varying vec3 localPos;
                
                void main() {
                    // vec4 pos = modelViewMatrix * vec4(position, 1.0);
                    // gl_Position = projectionMatrix * pos;

                    fNormal  = normalize(normalMatrix * normal);
                    vec4 pos = modelViewMatrix * vec4(position, 1.0);
                    worldPos = pos.xyz;
                    localPos = position;
                    gl_Position = projectionMatrix * pos;
                }
            `],

            ["fragment", `
                precision highp float;

                // IMPORTANT!! - shader uniform passed from glo-lang code
                uniform float i;

                varying vec3 fNormal;
                varying vec3 localPos;

                float pulse(float val, float dst) {
                    return floor(mod(val * dst, 1.0) + i);
                }

                void main() {
                    vec3 dir   = vec3(0.5, 1, 0); // high noon
                    vec3 cpos  = localPos;
                    vec3 color = vec3(1, pulse(cpos.y, 10.0), 1); 

                    float diffuse = .5 + dot(fNormal, dir);
                    gl_FragColor  = vec4(diffuse * color, 1.0);
                }

                /*void main() {
                    gl_FragColor = vec4(1, 1, 0, 1.0);
                }*/
            `]
        ]
    ];



    const shader_example__program_ast_lst = [

        ["lang_v", "0.0.3"], 
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

        // CENTRAL_SPHERE
        [
            ["set", "material", ["wireframe", true]],
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.3]]],
            ["y", -5],
            ["sx", -0.7, "sy", -0.7, "sz", -0.7],
            ["sphere"]
        ],

        // CENTRAL_SQUARE
        [
            
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.8]]],
            ["cube"]
        ],

        [
            ["sx", -0.8, "sy", -0.8, "sz", -0.8],
            ["y", 1.5, "x", 2],
            "R1"
        ],
        [
            ["sx", -0.8, "sy", -0.8, "sz", -0.8],
            ["y", 1.5, "x", -2],
            "R1"
        ],
        [
            ["sx", -0.8, "sy", -0.8, "sz", -0.8],
            ["y", 1.5, "z", -2],
            "R1"
        ],
        [
            ["sx", -0.8, "sy", -0.8, "sz", -0.8],
            ["y", 1.5, "z", 2],
            "R1"
        ],

        [
            ["set", "scale", [0.13, 0.13, 0.13]],
            ["x", 1],
            "R2"
        ],
        [
            ["set", "scale", [0.1, 0.1, 0.1]],
            ["x", 2],
            "R2"
        ],
        [
            ["set", "scale", [0.05, 0.05, 0.05]],
            ["x", 3],
            "R3"
        ],
        
        //-------------------------------------------------
        ["rule", "R1", [["iters_max", 40]], [
            ["set", "color", extrusion_color_2_str],
            [["sx", -0.005, "sy", -0.005, "sz", -0.005], ["y", 0.2, 'ry', 0.01]],
            ["x", -0.1, "z", -0.1],
            ["cube"],
            ["R1"]
        ]],
        ["rule", "R1", [["iters_max", 40]], [
            ["set", "color", extrusion_color_2_str],
            [["sx", -0.005, "sy", -0.005, "sz", -0.005], ["y", 0.2, 'ry', 0.01]],
            ["x", 0.1, "z", 0.1],

            /*["if", [">", "$i", 14], [
                ["R2"]
            ]],*/

            ["cube"],
            ["R1"]
        ]],
        ["rule", "R1", [["iters_max", 40]], [
            ["set", "color", extrusion_color_2_str],
            [["sx", -0.005, "sy", -0.005, "sz", -0.005], ["y", 0.2, 'ry', 0.01]],
            ["x", -0.1, "z", 0.1],
            ["cube"],
            ["R1"]
        ]],
        ["rule", "R1", [["iters_max", 40]], [
            ["set", "color", extrusion_color_2_str],
            [["sx", -0.005, "sy", -0.005, "sz", -0.005], ["y", 0.2, 'ry', 0.01]],
            ["x", 0.1, "z", -0.1],
            ["cube"],
            ["R1"]
        ]],

        //-------------------------------------------------

        ["rule", "R2", [["iters_max", 31]], [
            ['ry', 0.2],
            ["cube"],
            ["R2"]
        ]],
        ["rule", "R3", [["iters_max", 31]], [
            ['ry', 0.2],
            ["cube"],
            ["R3"]
        ]],

        //-------------------------------------------------

        ["shader", "gf_shader_test",
            ["uniforms", [
                ["i", "float", 0.1]
            ]],
            ["vertex", `

                varying vec3 fNormal;
                varying vec3 worldPos;
                varying vec3 localPos;
                
                void main() {
                    // vec4 pos = modelViewMatrix * vec4(position, 1.0);
                    // gl_Position = projectionMatrix * pos;

                    fNormal  = normalize(normalMatrix * normal);
                    vec4 pos = modelViewMatrix * vec4(position, 1.0);
                    worldPos = pos.xyz;
                    localPos = position;
                    gl_Position = projectionMatrix * pos;
                }
            `],

            ["fragment", `
                precision highp float;

                // IMPORTANT!! - shader uniform passed from glo-lang code
                uniform float i;

                varying vec3 fNormal;
                varying vec3 localPos;

                float pulse(float val, float dst) {
                    return floor(mod(val * dst, 1.0) + i);
                }

                void main() {
                    vec3 dir   = vec3(0, 1, 0); // high noon
                    vec3 cpos  = localPos;
                    vec3 color = vec3(1, pulse(cpos.y, 2.0), 0); 

                    float diffuse = .5 + dot(fNormal, dir);
                    gl_FragColor  = vec4(diffuse * color, 1.0);
                }

                /*void main() {
                    gl_FragColor = vec4(1, 1, 0, 1.0);
                }*/
            `]
        ]
    ];




    const feather_color_1_str = "#e2b31d"; // yellow
    const feather_color_2_str = "#393352"; // dar blue
    const feather_color_3_str = "#761f25"; // red
    const feather_color_4_str = "#d7cdb9";
    const feather_color_5_str = "#9d714b";
    const feather_color_6_str = "#daba70"; // light-brown
    const feather_color_7_str = "#6c7c8c"; // light-blue
    const feather_color_8_str = "#e1cc26"; // yellow
    const feather_color_9_str = "#ac8c7c";
    
    const planes_world__program_ast_lst = [

        ["lang_v", "0.0.3"], 
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],
        ["set", "iters_max", 500000],

        [
            ["set", "color", extrusion_color_2_str],
            ["set", "material", ["wireframe", true]],


            [
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
        ],

        [
            ["set", "color", feather_color_1_str],
            ["set", "material", ["wireframe", true]],

            ["y", 5.3, "x", 0.6, "rz", 0.3],

            [
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
        ],

        [
            ["set", "color", feather_color_3_str],
            ["set", "material", ["wireframe", true]],

            ["y", 10.6, "x", 3, "rz", 0.6],

            [
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
        ],

        [
            ["set", "color", feather_color_9_str],
            ["set", "material", ["wireframe", true]],

            ["y", 15.6, "x", 6, "rz", 0.9],

            [
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
            [
                ["z", 2],
                ["set", "scale", [0.1, 5.0, 1.0]],  
                "feather_plane"
            ],
        ],
        //-----------------------------------------------------------------
        // FEATHER_PLANE
        ["rule", "feather_plane", [["iters_max", 20]], [
            [
                ["x", 0.4, "rx", -0.03],
                ["cr", -0.03, "cg", -0.03, "cb", -0.03],
                "cube"
            ],
            ["feather_plane"]
        ]],
        ["rule", "feather_plane", [["iters_max", 20]], [
            [
                ["x", 0.4, "rx", -0.03],
                ["cr", -0.005, "cg", -0.005, "cb", -0.005],
                "cube"
            ],
            
            ["if", ["==", "$i", 19], [
                [["y", -1, "sy", -4], "feather_child_plane"]
            ]],

            ["feather_plane"]
        ]],
        ["rule", "feather_plane", [["iters_max", 20]], [
            [
                ["x", 0.4, "rx", -0.03],
                ["cr", 0.01, "cg", 0.01, "cb", 0.01],
                "cube"
            ],

            ["if", ["==", "$i", 19], [
                [["y", -1, "sy", -4], "feather_child_plane"]
            ]],

            ["feather_plane"]
        ]],

        //-----------------------------------------------------------------
        // FEATHER_CHILD_PLANE
        ["rule", "feather_child_plane", [["iters_max", 60]], [
            [
                ["y", 1, "x", 0.02, "rz", -0.01],
                "cube"
            ],

            ["if", [">", "$i", 55], [
                [
                    // ["set", "material", ["wireframe", false]],
                    ["set", "scale", [1, 1, 1]],  
                    ["y", 1], 
                    "child"
                ],
                // ["set", "material", ["wireframe", true]],
            ]],

            ["feather_child_plane"]
        ]],
        ["rule", "feather_child_plane", [["iters_max", 60]], [
            [
                ["y", 0.7, "x", 0.08, "rz", -0.015],
                ["cr", 0.01, "cg", 0.01, "cb", 0.01],
                "cube"
            ],
            ["feather_child_plane"]
        ]],

        //-----------------------------------------------------------------
        // CHILD
        ["rule", "child", [["iters_max", 25]], [
            [
                ["y", 0.6, "x", 0.2, "z", -0.2],
                "cube"
            ],
            ["child"]
        ]],
        ["rule", "child", [["iters_max", 25]], [
            [
                // ["set", "material", ["wireframe", true]],
                ["cr", 0.01, "cg", 0.01, "cb", 0.01],
                ["y", 0.7, "x", -0.2, "y", 0.2],
                "cube"
            ],
            ["child"]
        ]],
        ["rule", "child", [["iters_max", 25]], [
            [
                ["cr", -0.01, "cg", -0.01, "cb", -0.01],
                // ["set", "material", ["wireframe", true]],
                ["y", 0.5, "x", -0.2, "z", 0.2],
                "cube"
            ],

            ["if", [">", "$i", 24], [
                [
                    ["set", "scale", [0.2, 0.2, 0.2]],  
                    ["y", 0.1], 
                    "seed"
                ],
                // ["set", "material", ["wireframe", true]],
            ]],

            ["child"]
        ]],
        ["rule", "child", [["iters_max", 25]], [
            [
                ["cr", 0.01, "cg", 0.01, "cb", 0.01],
                // ["set", "material", ["wireframe", false]],
                ["y", 0.4, "x", 0.2, "z", -0.2],
                "cube"
            ],
            ["child"]
        ]],

        //-----------------------------------------------------------------
        ["rule", "seed", [["iters_max", 120]], [
            [
                ["set", "color", ["rgb", 0.2, 0.2, 0.2]],
                ["y", 0.2],
                "cube"
            ],
            ["seed"]
        ]],
        ["rule", "seed", [["iters_max", 120]], [
            [
                ["set", "color", ["rgb", 0.2, 0.2, 0.2]],
                ["z", 0.2],
                "cube"
            ],
            ["seed"]
        ]],
        ["rule", "seed", [["iters_max", 120]], [
            [
                ["set", "color", ["rgb", 0.2, 0.2, 0.2]],
                ["z", -0.2],
                "cube"
            ],
            ["seed"]
        ]],


    ];




    // 0.0.4 - big additions
    // - addition of a mostly generalized animation system for both individual
    //   rule objects and camera/global properties.
    // - added color-animation design that works along the general property-animation pattern.
    // - basic IDE added (integrated development environment);
    //      - allow for setting the background color via color-picker
    //      - display of a helper grid and axis for user origentation
    //      - toggle for view/hide axes/grid
    //      - start/stop button for global/camera animations
    //      - individual animation sequences display and ability to change time durations of each sequence
    //      - playback of animations
    //      - support for repeatable animations
    // - line display
    //      - new state setters for starting new line sequences, when a recursive rule is used that uses "line" rule
    // - state setters can contain arbitrary expressions now for its operands (not just numbers or strings anymore).
    const shader_colors_and_sys_funs_test__program_ast_lst = [

        ["lang_v", "0.0.4"], 
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

        [
            ["set", "color", ["rgb", 0.2, 0.2, 0.2]],

            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.8]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
            ["y", 2, "z", 2], 
            "cube"
        ],
        [
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "color", ["rgb", 0.2, 0.8, 0.2]],
            ["x", 2],
            "seed"
        ],
        [
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "color", ["rgb", 0.5, 0.8, 0.2]],
            ["x", 4],
            "seed"
        ],
        [
            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "color", ["rgb", 0.5, 0.8, 0.6]],
            ["x", 6],
            "seed"
        ],
        [
            ["x", 8, "y", 4],
            ["set", "line", ["start"]],
            ["set", "scale", [0.3, 0.3, 0.3]], 
            "line_draw"
        ],
        [
            ["x", 8.5, "y", 4],
            ["set", "line", ["start"]],
            ["set", "scale", [0.3, 0.3, 0.3]], 
            "line_draw"
        ],
        [
            ["x", 9, "y", 4],
            ["set", "line", ["start"]],
            ["set", "scale", [0.3, 0.3, 0.3]], 
            "line_draw"
        ],
        [
            ["x", 9.5, "y", 4],
            ["set", "line", ["start"]],
            ["set", "scale", [0.3, 0.3, 0.3]], 
            "line_draw"
        ],
        [
            ["x", 10, "y", 4],
            ["set", "line", ["start"]],
            ["set", "scale", [0.3, 0.3, 0.3]], 
            "line_draw"
        ],
        [
            ["x", 10.5, "y", 4],
            ["set", "line", ["start"]],
            ["set", "scale", [0.3, 0.3, 0.3]], 
            "line_draw"
        ],
        [
            ["x", 11, "y", 4],
            ["set", "line", ["start"]],
            ["set", "scale", [0.3, 0.3, 0.3]], 
            "line_draw"
        ],

        //-----------------------------------------------------------------
        ["rule", "line_draw", [["iters_max", 30]], [

            ["if", ["==", "$i", 20], [
                ["set", "line", ["start"]],
            ]],

            [
                [
                    "z", ["rand", [0.0, 2.0]],
                    "x", ["rand", [-1.0, 1.0]],
                    "y", ["rand", [-0.7, 0.7]],
                    "rz", 0.3
                ],
                "line"
            ],
            ["cube"],
            ["line_draw"]
        ]],

        //-----------------------------------------------------------------
        ["rule", "seed", [["iters_max", 10]], [
            [
                
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],

                ["z", ["+", 1, ["rand", [0.0, 0.5]]], "cg", -0.05],
                "cube"
            ],
            ["seed"]
        ]],

        //-----------------------------------------------------------------
        ["shader", "gf_shader_test",
            ["uniforms", [
                ["i",  "float", 0.1],
                ["cr", "float", 0.5],
                ["cg", "float", 0.5],
                ["cb", "float", 0.5],
            ]],
            ["vertex", `

                varying vec3 normal_f;
                varying vec3 local_pos_v3;
                
                void main() {

                    normal_f     = normalize(normalMatrix * normal);
                    local_pos_v3 = position;

                    vec4 pos     = modelViewMatrix * vec4(position, 1.0);
                    gl_Position  = projectionMatrix * pos;
                }
            `],

            ["fragment", `
                precision highp float;

                // IMPORTANT!! - shader uniform passed from glo-lang code
                uniform float i;
                uniform float cr;
                uniform float cg;
                uniform float cb;

                varying vec3 normal_f;
                varying vec3 local_pos_v3;

                void main() {
                    vec3 light_direction_v3 = vec3(0.8, 1, 0.8); // high noon
                    vec3 color_v3           = vec3(cr, cg, cb); 

                    float diffuse_f = .5 + dot(normal_f, light_direction_v3);
                    gl_FragColor    = vec4(diffuse_f * color_v3, 1.0);
                }
            `]
        ]

        //-----------------------------------------------------------------
    ];

    const animation_test__program_ast_lst = [

        ["lang_v", "0.0.4"], 
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

        [
            ["set", "color", ["rgb", 0.2, 0.2, 0.2]],


            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.8]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
            ["x", 2, "y", 2, "z", 2],



            ["animate", [["x", 10], ["z", 2], ["y", 20]], 2, "repeat"],
            "cube"
        ],

        //-----------------------------------------------------------------
        ["shader", "gf_shader_test",
            ["uniforms", [
                ["i",  "float", 0.1],
                ["cr", "float", 0.5],
                ["cg", "float", 0.5],
                ["cb", "float", 0.5],
            ]],
            ["vertex", `

                varying vec3 normal_f;
                varying vec3 local_pos_v3;
                
                void main() {

                    normal_f     = normalize(normalMatrix * normal);
                    local_pos_v3 = position;

                    vec4 pos     = modelViewMatrix * vec4(position, 1.0);
                    gl_Position  = projectionMatrix * pos;
                }
            `],

            ["fragment", `
                precision highp float;

                // IMPORTANT!! - shader uniform passed from glo-lang code
                uniform float i;
                uniform float cr;
                uniform float cg;
                uniform float cb;

                varying vec3 normal_f;
                varying vec3 local_pos_v3;

                void main() {
                    vec3 light_direction_v3 = vec3(0.8, 1, 0.8); // high noon
                    vec3 color_v3           = vec3(cr, cg, cb); 

                    float diffuse_f = .5 + dot(normal_f, light_direction_v3);
                    gl_FragColor    = vec4(diffuse_f * color_v3, 1.0);
                }
            `]
        ]

        //-----------------------------------------------------------------
    ];

    const gray_fabric__program_ast_lst = [

        ["lang_v", "0.0.4"], 
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

        [
            ["set", "color", ["rgb", 0.2, 0.2, 0.2]],

            ["set", "material", ["shader", "gf_shader_test"]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.8]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],

            ["x", 0, "y", 20, "z", 2],
            ["*", 20, ["x", 1.3, "stripe"]]
        ],

        //-----------------------------------------------------------------
        // STRIPE
        ["rule", "stripe", [["iters_max", 50]], [
            ["if", ["==", "$i", 0], [["set", "line", ["start"]]]],
            [
                ["y", -1.3],
                ["z", 0.1],

                // ["rand", [0.0, 0.5]]
                [
                    "sx", ["*", -0.0010, "$i"],
                    "sy", ["*", -0.005, "$i"],
                    "sz", ["*", -0.0010, "$i"],
                ],
                // ["animate", [["z", ["rand", [-0.5, 0.5]]]], 1, "repeat"],
                "cube"
            ],
            ["line"],
            ["stripe"]
        ]],
        ["rule", "stripe", [["iters_max", 50]], [
            ["if", ["==", "$i", 0], [["set", "line", ["start"]]]],
            [
                ["y", -1.3],
                ["z", -0.1],
                [
                    "sx", ["*", -0.0010, "$i"],
                    "sy", ["*", -0.0010, "$i"],
                    "sz", ["*", -0.0010, "$i"],
                ],
                [
                    ["cr", 0.01, "cg", 0.01, "cb", 0.01],
                    ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
                    ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
                    ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
                ],
                "cube"
            ],
            ["line"],
            ["stripe"]
        ]],

        ["rule", "stripe", [["iters_max", 50]], [
            ["if", ["==", "$i", 0], [["set", "line", ["start"]]]],
            [
                ["y", -2.3],
                ["z", 0.1],
                [
                    "sx", ["*", -0.0010, "$i"],
                    "sy", ["*", -0.0010, "$i"],
                    "sz", ["*", -0.0010, "$i"],
                ],
                ["set", "color", ["rgb", 0.1, 0.1, 0.1]],
                
                ["if", [">", "$i", 40], [
                    "branch"
                ]],
                "cube"
            ],
            ["line"],
            ["stripe"]
        ]],
        
        //-----------------------------------------------------------------
        // BRANCH
        ["rule", "branch", [["iters_max", 20]], [
            ["set", "scale", [0.1, 0.1, 0.1]],
            [
                ["set", "color", ["rgb", 0.7, 0.7, 0.7]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
            ],
            [
                ["z", 0.2],
                "cube"
            ],
            ["if", [">", "$i", ["rand", [15, 17]]], [
                ["z", 0.5],

                ["set", "material", ["wireframe", true]],
                ["crown"],
                ["set", "material", ["wireframe", false]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
            ]],
            ["if", ["==", "$i", 19], [
                "crown_top"
            ]],
            "branch"
        ]],

        //-----------------------------------------------------------------
        // CROWN
        ["rule", "crown", [["iters_max", 1]], [
            ["set", "scale", [
                ["rand", [0.3, 0.9]],
                ["rand", [0.3, 0.9]],
                ["rand", [0.3, 0.9]]]],
            [   
                ["set", "color", "#ff5512"],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
                "cube"
            ],
            "crown"
        ]],
        ["rule", "crown", [["iters_max", 1]], [
            ["set", "scale", [
                ["rand", [0.3, 0.9]],
                ["rand", [0.3, 0.9]],
                ["rand", [0.3, 0.9]]]],
            [   
                ["set", "color", "#ff8c00"],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
                "cube"
            ],
            "crown"
        ]],

        //-----------------------------------------------------------------
        ["rule", "crown_top", [["iters_max", 20]], [
            ["if", ["==", "$i", 0], [["set", "line", ["start"]]]],

            ["set", "scale", [0.05, 0.05, 0.05]],
            [
                ["set", "color", ["rgb", 0.7, 0.7, 0.7]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
                ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
            ],
            [
                ["z", 0.2],
                "cube"
            ],
            ["line"],
            "crown_top"
        ]],

        //-----------------------------------------------------------------
        ["shader", "gf_shader_test",
            ["uniforms", [
                ["i",  "float", 0.1],
                ["cr", "float", 0.5],
                ["cg", "float", 0.5],
                ["cb", "float", 0.5],
            ]],
            ["vertex", `

                varying vec3 normal_f;
                varying vec3 local_pos_v3;
                
                void main() {

                    normal_f     = normalize(normalMatrix * normal);
                    local_pos_v3 = position;

                    vec4 pos     = modelViewMatrix * vec4(position, 1.0);
                    gl_Position  = projectionMatrix * pos;
                }
            `],

            ["fragment", `
                precision highp float;

                // IMPORTANT!! - shader uniform passed from glo-lang code
                uniform float i;
                uniform float cr;
                uniform float cg;
                uniform float cb;

                varying vec3 normal_f;
                varying vec3 local_pos_v3;

                void main() {
                    vec3 light_direction_v3 = vec3(0.0, 1, 1.0); // high noon
                    vec3 color_v3           = vec3(cr, cg, cb); 

                    float diffuse_f = .5 + dot(normal_f, light_direction_v3);
                    gl_FragColor    = vec4(diffuse_f * color_v3, 1.0);
                }
            `]
        ]

        //-----------------------------------------------------------------
    ];



    const stripe_color_str = "#39b812";
    const ring_color_str   = "#ff8000";
    const ring_2_color_str = "#004cff";

    function set_shader_color(p_color_hex_str) {
        return [
            ["set", "color", p_color_hex_str],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
        ]
    }
    const rotation_pivot_setters_test__program_ast_lst = [

        ["lang_v", "0.0.4"], 
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

        [
            ["set", "material", ["shader", "gf_shader_test"]],
            set_shader_color(stripe_color_str),
        ],
        
        [
            ["x", -15, "z", -5],
            "stripe"
        ],

        //-----------------------------------------------------------------
        // STRIPE
        ["rule", "stripe", [["iters_max", 5]], [

            ["if", ["==", "$i", 3], [
                ["stripe_2"]
            ]],

            [
                ["y", 2],
                "cube"
            ],
            ["stripe"]
        ]],

        //-----------------------------------------------------------------
        // STRIPE_2
        ["rule", "stripe_2", [["iters_max", 10]], [

            ["if", ["==", "$i", 5], [

                set_shader_color(ring_color_str),

                ["push", "rotation_pivot", "current_pos"],
                ["set", "line", ["start"]],
                ["y", 1.5],
                ["ring"],
                ["y", -1.5],
                ["pop", "rotation_pivot", "current_pos"],

                set_shader_color(stripe_color_str),
            ]],
            ["if", ["==", "$i", 7], [

                set_shader_color(ring_color_str),

                ["push", "rotation_pivot", "current_pos"],
                ["set", "line", ["start"]],
                ["y", 3],
                ["ring"],
                ["y", -3],
                ["pop", "rotation_pivot", "current_pos"],

                set_shader_color(stripe_color_str),
            ]],
            ["if", ["==", "$i", 8], [

                set_shader_color(ring_color_str),

                ["push", "rotation_pivot", "current_pos"],
                ["set", "line", ["start"]],
                ["y", 5],
                ["ring"],
                ["y", -5],
                ["pop", "rotation_pivot", "current_pos"],

                set_shader_color(stripe_color_str),
            ]],
            ["if", ["==", "$i", 9], [

                set_shader_color(ring_color_str),

                ["push", "rotation_pivot", "current_pos"],
                ["set", "line", ["start"]],
                ["y", 8],
                ["ring"],
                ["y", -8],
                ["pop", "rotation_pivot", "current_pos"],
                
                set_shader_color(stripe_color_str),
            ]],

            [["x", 3], "cube"],
            ["stripe_2"]
        ]],

        //-----------------------------------------------------------------
        // RING
        ["rule", "ring", [["iters_max", 25]], [

            ["if", ["==", "$i", 6], [
                
                set_shader_color(ring_2_color_str),

                ["push", "rotation_pivot", "current_pos"],
                //["z", 2],
                ["ring_2"],
                ["line"],
                //["z", -2],
                ["pop", "rotation_pivot", "current_pos"],

                set_shader_color(ring_color_str),
            ]],

            ["set", "scale", [0.3, 0.3, 0.3]],
            [["rx", 0.25], "cube"],

            ["line"],
            ["ring"]
        ]],

        //-----------------------------------------------------------------
        // RING_2
        ["rule", "ring_2", [["iters_max", 50]], [
            ["set", "scale", [0.2, 0.2, 0.2]],
            [["ry", 0.25], "cube"],
            ["line"],
            ["ring_2"]
        ]],

        //-----------------------------------------------------------------
        ["shader", "gf_shader_test",
            ["uniforms", [
                ["cr", "float", 0.5],
                ["cg", "float", 0.5],
                ["cb", "float", 0.5],
            ]],
            ["vertex", `

                varying vec3 normal_f;
                varying vec3 local_pos_v3;
                
                void main() {

                    normal_f     = normalize(normalMatrix * normal);
                    local_pos_v3 = position;

                    vec4 pos     = modelViewMatrix * vec4(position, 1.0);
                    gl_Position  = projectionMatrix * pos;
                }
            `],

            ["fragment", `
                precision highp float;

                // IMPORTANT!! - shader uniform passed from glo-lang code
                uniform float i;
                uniform float cr;
                uniform float cg;
                uniform float cb;

                varying vec3 normal_f;
                varying vec3 local_pos_v3;

                void main() {
                    vec3 light_direction_v3 = vec3(0.8, 1, 0.8); // high noon
                    vec3 color_v3           = vec3(cr, cg, cb); 

                    float diffuse_f = .5 + dot(normal_f, light_direction_v3);
                    gl_FragColor    = vec4(diffuse_f * color_v3, 1.0);
                }
            `]
        ]

        //-----------------------------------------------------------------
    ];



    return {
        "first_scene__program_ast_lst":        first_scene__program_ast_lst,
        "rules_test__program_ast_lst":         rules_test__program_ast_lst,
        "multi_rule_test__program_ast_lst":    multi_rule_test__program_ast_lst,
        "color_test__program_ast_lst":         color_test__program_ast_lst,
        "conditional_test__program_ast_lst":   conditional_test__program_ast_lst,
        "conditional_test_2__program_ast_lst": conditional_test_2__program_ast_lst,



        "form_experiment__ast_lst":        form_experiment__ast_lst,
        "sphere_test__program_ast_lst":    sphere_test__program_ast_lst,
        "shader_test__program_ast_lst":    shader_test__program_ast_lst,
        "shader_example__program_ast_lst": shader_example__program_ast_lst,
        "planes_world__program_ast_lst":   planes_world__program_ast_lst,
    
    
        "shader_colors_and_sys_funs_test__program_ast_lst": shader_colors_and_sys_funs_test__program_ast_lst,
        "animation_test__program_ast_lst": animation_test__program_ast_lst,
        "gray_fabric__program_ast_lst":    gray_fabric__program_ast_lst,
        "rotation_pivot_setters_test__program_ast_lst": rotation_pivot_setters_test__program_ast_lst,
    
    
    }
}