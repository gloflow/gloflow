
///<reference path="../../../../../gloflow/web/src/d/jquery.d.ts" />

import * as gf_lang     from "../ts/gf_lang";
import * as gf_examples from "./gf_examples";

//-------------------------------------------------
$(document).ready(()=>{

    const examples_map                        = gf_examples.get();
    const first_scene__program_ast_lst        = examples_map["first_scene__program_ast_lst"];
    const rules_test__program_ast_lst         = examples_map["rules_test__program_ast_lst"];
    const multi_rule_test__program_ast_lst    = examples_map["multi_rule_test__program_ast_lst"];
    const color_test__program_ast_lst         = examples_map["color_test__program_ast_lst"];
    const conditional_test__program_ast_lst   = examples_map["conditional_test__program_ast_lst"];
    const conditional_test_2__program_ast_lst = examples_map["conditional_test_2__program_ast_lst"];
    const form_experiment__ast_lst            = examples_map["form_experiment__ast_lst"];
    const sphere_test__program_ast_lst        = examples_map["sphere_test__program_ast_lst"];
    const shader_test__program_ast_lst        = examples_map["shader_test__program_ast_lst"];
    const shader_example__program_ast_lst     = examples_map["shader_example__program_ast_lst"];
    const planes_world__program_ast_lst       = examples_map["planes_world__program_ast_lst"];
    const gray_fabric__program_ast_lst        = examples_map["gray_fabric__program_ast_lst"];
    const rotation_pivot_setters_test__program_ast_lst = examples_map["rotation_pivot_setters_test__program_ast_lst"];


    const stripe_color_str   = "#ff4400";
    const stripe_2_color_str = "#ff1100";
    const stripe_3_color_str = "#9955ff";
    const stripe_4_color_str = "#44ffff";

    function set_shader_color(p_color_hex_str) {
        return [
            ["set", "color", p_color_hex_str],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]],
            ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]],
        ]
    }
    const origin_setters_test__program_ast_lst = [

        ["lang_v", "0.0.4"], 
        ["set", "color-background", ["rgb", 0.90, 0.90, 0.90]],

        [
            ["set", "material", ["shader", "gf_shader_test"]],
            set_shader_color(stripe_color_str),
        ],
        
        [
            ["z", 5],
            "stripe"
        ],

        //-----------------------------------------------------------------
        // STRIPE
        ["rule", "stripe", [["iters_max", 6]], [

            ["if", ["==", "$i", 5], [

                set_shader_color(stripe_2_color_str),

                // COORD_ORIGIN
                ["push", "coord_origin", "current_pos"],

                ["set", "scale", [0.8, 0.8, 0.8]],
                ["stripe_2"],
                ["set", "scale", [1.0, 1.0, 1.0]],

                ["pop", "coord_origin", "current_pos"],

                set_shader_color(stripe_color_str),
            ]],

            [
                ["z", 0.3, "rx", -0.05, "ry", 0.5],
                "cube"
            ],
            ["stripe"]
        ]],

        //-----------------------------------------------------------------
        // STRIPE_2
        ["rule", "stripe_2", [["iters_max", 10]], [

            ["if", ["==", "$i", 6], [
                
                set_shader_color(stripe_3_color_str),

                // COORD_ORIGIN
                ["push", "coord_origin", "current_pos"],

                ["set", "scale", [0.6, 0.6, 0.6]],
                ["stripe_3"],
                ["set", "scale", [0.8, 0.8, 0.8]],

                ["pop", "coord_origin", "current_pos"],

                set_shader_color(stripe_2_color_str),
            ]],

            [["y", 0.2, "rz", 0.2, "ry", 0.05], "cube"],
            ["stripe_2"]
        ]],

        //-----------------------------------------------------------------
        // STRIPE_3
        ["rule", "stripe_3", [["iters_max", 10]], [

            ["if", ["==", "$i", 5], [
                
                set_shader_color(stripe_4_color_str),

                // COORD_ORIGIN
                ["push", "coord_origin", "current_pos"],

                ["set", "scale", [0.4, 0.4, 0.4]],
                ["stripe_4"],
                ["set", "scale", [0.6, 0.6, 0.6]],

                ["pop", "coord_origin", "current_pos"],

                set_shader_color(stripe_3_color_str),
            ]],

            [["z", 0.1, "rx", 0.2, "ry", 0.3], "cube"],
            ["stripe_3"]
        ]],

        //-----------------------------------------------------------------
        // STRIPE_4
        ["rule", "stripe_4", [["iters_max", 10]], [
            [["x", -1], "cube"],
            ["stripe_4"]
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




    gf_lang.run(origin_setters_test__program_ast_lst);

});