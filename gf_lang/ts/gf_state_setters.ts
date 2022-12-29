/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

import * as gf_state      from "./gf_state";
import * as gf_lang_utils from "./gf_lang_utils";

//-------------------------------------------------
export function exec__expr(p_setter_type_str :string,
    p_property_name_str :string,
    p_vals,
    p_state_map,
    p_state_family_stack_lst,
    p_engine_api_map) {
    
    const symb_map = gf_lang_utils.get_symbols_and_constants();

    if (p_setter_type_str != "set" && p_setter_type_str != "push" && p_setter_type_str != "pop")
        throw "state setter is not of type 'set|push|pop'";

    //------------------------------------
    // SCALE
    if (p_property_name_str == "scale") {
        
        if (typeof p_vals == "number") {
            const scale_val = p_vals;
            p_state_map["sx"] = scale_val;
            p_state_map["sy"] = scale_val;
            p_state_map["sz"] = scale_val;
        }
        else if (Array.isArray(p_vals)) {

            const vals_lst = p_vals;
            if (vals_lst.length != 3) throw "scale values can only be of length 3 [sx, sy, sz]";

            const [sx, sy, sz] = vals_lst;

            p_state_map["sx"] = gf_lang_utils.expr_eval(sx, p_state_map);;
            p_state_map["sy"] = gf_lang_utils.expr_eval(sy, p_state_map);;
            p_state_map["sz"] = gf_lang_utils.expr_eval(sz, p_state_map);;
        }
    }

    //------------------------------------
    // COLOR
    else if (p_property_name_str == "color") {
        
        if (Array.isArray(p_vals)) {
            const vals_lst = p_vals;
            if (vals_lst[0] != "rgb") throw "only rgb type is allowed";
            if (vals_lst.length != 4) throw "rgb values can only be of length 4 ['rgb', r, g, b]";
            
            const [, r, g, b] = vals_lst;
            const state_change_map = {
                "setter_type_str": p_setter_type_str,
                "color_rgb":       [r, g, b],
            };
            p_engine_api_map["set_state_fun"](state_change_map);

            p_state_map["cr"] = r;
            p_state_map["cg"] = g;
            p_state_map["cb"] = b;
        }

        else {
            if (typeof p_vals != "string") throw "setting color has to be either an array or string";
            if (!p_vals.startsWith("#")) throw "setting color with a string has to be done in hex format starting with #";

            const color_hex_str = p_vals;
            const state_change_map = {
                "setter_type_str": p_setter_type_str,
                "color_rgb":       color_hex_str
            };
            // set_state_fun() will return parsed color hex string, with rgb channels in 0-1 range.
            const [r, g, b] = p_engine_api_map["set_state_fun"](state_change_map);

            p_state_map["cr"] = r;
            p_state_map["cg"] = g;
            p_state_map["cb"] = b;
        }
    }

    //------------------------------------
    // COLOR_BACKGROUND
    else if (p_property_name_str == "color-background") {
        const vals_lst = p_vals;
        if (vals_lst[0] != "rgb") throw "only rgb type is allowed";
        if (vals_lst.length != 4) throw "rgb values can only be of length 4 ['rgb', r, g, b]";
        
        const [, r, g, b] = vals_lst;
        const state_change_map = {
            "setter_type_str":  p_setter_type_str,
            "color_background": [r, g, b],
        };
        p_engine_api_map["set_state_fun"](state_change_map);
    }
    
    //------------------------------------
    // ITERS_MAX
    else if (p_property_name_str == "iters_max") {
        const iterations_max_int = p_vals;
        p_state_map["iters_max"] = iterations_max_int;
    }

    //------------------------------------
    // MATERIAL
    else if (p_property_name_str == "material") {
        const [material_type_str, val_str] = p_vals;
        if (material_type_str != "wireframe" 
            && material_type_str != "shader") throw "only 'wireframe|shader' material types are supported";

        const state_change_map = {
            "setter_type_str":    p_setter_type_str,
            "material_type_str":  material_type_str,
            "material_value_str": val_str
        };
        p_engine_api_map["set_state_fun"](state_change_map);                
    }

    //------------------------------------
    // MATERIAL_PROPERTY
    else if (p_property_name_str == "material_prop") {
        const [material_name_str, material_prop_str, material_prop_val] = p_vals;
        if (material_prop_str != "shader_uniform") throw "only 'shader_uniform' material properties are supported";

        if (material_prop_str == "shader_uniform") {
            const [uniform_name_str, uniform_val] = material_prop_val;

            var loaded_val;

            // VARIABLE_REFERENCE
            if (typeof uniform_val == "string" && uniform_val.startsWith("$")) {

                const possible_prop_name_str = uniform_val.slice(1); // remove "$"

                // SYSTEM_PROPERTY - x|y|z|...|cr|cg|cb
                if (symb_map["predefined_properties_lst"].includes(possible_prop_name_str)) {
                    const property_name_str = possible_prop_name_str;
                    loaded_val = p_state_map[property_name_str];
                } 
                // USER_DEFINED_VARIABLE
                else {
                    // evalue the variable reference to get its value
                    loaded_val = p_state_map["vars_map"][uniform_val];
                }
            }

            // ARITHMETIC_EXPRESSION
            else if (Array.isArray(uniform_val)) {
                const sub_expr_lst = uniform_val;

                const mul_result = gf_lang_utils.arithmetic_eval(sub_expr_lst, p_state_map)
                loaded_val       = mul_result;
            }

            // NUMERIC_VALUE
            else {
                loaded_val = uniform_val;
            }
            const state_change_map = {
                "setter_type_str": p_setter_type_str,
                "material_prop_map": {
                    "material_shader_name_str":         material_name_str,
                    "material_shader_uniform_name_str": uniform_name_str,
                    "material_shader_uniform_val":      loaded_val,
                }
            };
            p_engine_api_map["set_state_fun"](state_change_map);        
        }
    }

    //------------------------------------
    // LINE
    else if (p_property_name_str == "line") {
        const [cmd_str] = p_vals;

        const state_change_map = {
            "setter_type_str": p_setter_type_str,
            "line_cmd_str": "start",
        };
        p_engine_api_map["set_state_fun"](state_change_map);
    }

    //------------------------------------
    // ROTATION_PIVOT
    else if (p_property_name_str == "rotation_pivot") {

        const axis_type_str = p_vals;
        if (axis_type_str == "current_pos") {

            const state_change_map = {
                "property_name_str": "rotation_pivot",
                "setter_type_str":   p_setter_type_str,
                "axis_type_str":     "current_pos",
                "x":  p_state_map["x"],
                "y":  p_state_map["y"],
                "z":  p_state_map["z"],
                "rx": p_state_map["rx"],
                "ry": p_state_map["ry"],
                "rz": p_state_map["rz"]
            };
            p_engine_api_map["set_state_fun"](state_change_map);
        }
    }

    //------------------------------------
    // COORD_ORIGIN - setting where the origin for subsequent operation should be.
    //                it can either be the current_position or world origin.
    else if (p_property_name_str == "coord_origin") {

        const origin_type_str = p_vals;
        if (origin_type_str != "current_pos")
            throw "'coord_origin' setter has to have an type of 'current_pos'";

        

        var new_state_map;
        switch (p_setter_type_str) {
            
            //------------------------------------
            case "push":

                const new_state_change_map = {
                    "property_name_str": "coord_origin",
                    "setter_type_str":   "push",
                    "origin_type_str":   origin_type_str,
        
                    "x":  p_state_map["x"],
                    "y":  p_state_map["y"],
                    "z":  p_state_map["z"],
                    "rx": p_state_map["rx"],
                    "ry": p_state_map["ry"],
                    "rz": p_state_map["rz"]
                }
                p_engine_api_map["set_state_fun"](new_state_change_map);

                //---------------------------
                // NEW_BLANK_STATE - only other place where this is being done
                //                   is at the root of the program execution.
                const new_family_state_map = gf_state.create_new_family(p_state_map);

                //---------------------------

                // this is the last state of the current family, that will need to be restored
                // to as the current state when the current family is popped
                const current_family_last_state_map = p_state_map;
                gf_state.push_family(current_family_last_state_map, p_state_family_stack_lst);
                
                new_state_map = new_family_state_map;
                break;

            //------------------------------------
            case "pop":
                
                const last_family_state_map = gf_state.pop_family(p_state_family_stack_lst);
                new_state_map = last_family_state_map;

                const restore_state_change_map = {
                    "property_name_str": "coord_origin",
                    "setter_type_str":   "pop",
                    "origin_type_str":   origin_type_str,
        
                    "x":  last_family_state_map["x"],
                    "y":  last_family_state_map["y"],
                    "z":  last_family_state_map["z"],
                    "rx": last_family_state_map["rx"],
                    "ry": last_family_state_map["ry"],
                    "rz": last_family_state_map["rz"]
                }
                p_engine_api_map["set_state_fun"](restore_state_change_map);

                break;

            //------------------------------------
        }

        return new_state_map;
    }

    //------------------------------------
}