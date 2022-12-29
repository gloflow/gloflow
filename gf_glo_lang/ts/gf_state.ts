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

import * as gf_glo_lang_utils from "./gf_glo_lang_utils";

//-------------------------------------------------
// FAMILY - state family is a new root state and all its children
//          that are merged into each other.
//          they have their own coordinate system origin.
//-------------------------------------------------
export function push_family(p_new_state_map, p_state_family_stack_lst) {
    p_state_family_stack_lst.push(p_new_state_map);
}

export function pop_family(p_state_family_stack_lst) {

    // this is the last state of the last state_family
    const last_family_state_map = p_state_family_stack_lst.pop();

    if (last_family_state_map == undefined)
        throw "family_state_stack has no more states to pop";

    return last_family_state_map;
}

//-------------------------------------------------
// create new state family, which is state like all others except
// it begins a new geometric space (a new coord system) but inherits some of the
// iteration counters and sub-stacks.

export function create_new_family(p_state_parent_map) {

    const state_map = get_empty_state();

    state_map["rules_names_stack_lst"]     = p_state_parent_map["rules_names_stack_lst"];
    state_map["vars_map"]                  = Object.assign({}, p_state_parent_map["vars_map"]); // clone
    state_map["iters_num_global_int"]      = p_state_parent_map["iters_num_global_int"];
    state_map["rules_iters_num_stack_lst"] = Object.assign([], p_state_parent_map["rules_iters_num_stack_lst"]); // clone

    return state_map;
}

//-------------------------------------------------
// VAR
//-------------------------------------------------
export function merge_child_state(p_state_map, p_child_state_map) {
    
    if (p_child_state_map == undefined)
        throw "supplied child_state_map is undefined";

    p_state_map["x"]  = p_child_state_map["x"];
    p_state_map["y"]  = p_child_state_map["y"];
    p_state_map["z"]  = p_child_state_map["z"];
    p_state_map["rx"] = p_child_state_map["rx"];
    p_state_map["ry"] = p_child_state_map["ry"];
    p_state_map["rz"] = p_child_state_map["rz"];
    p_state_map["sx"] = p_child_state_map["sx"];
    p_state_map["sy"] = p_child_state_map["sy"];
    p_state_map["sz"] = p_child_state_map["sz"];
    p_state_map["cr"] = p_child_state_map["cr"];
    p_state_map["cg"] = p_child_state_map["cg"];
    p_state_map["cb"] = p_child_state_map["cb"];
    p_state_map["iters_max"]             = p_child_state_map["iters_max"];
    p_state_map["rules_names_stack_lst"] = p_child_state_map["rules_names_stack_lst"];

    // rule iteration count ("$i") has to propagate up the expression tree,
    // from child expressions to parent expressions, as a part of the state.
    // however $i only travels up to the root of a particular rule.
    p_state_map["vars_map"]["$i"] = gf_glo_lang_utils.var_eval("$i", p_child_state_map); // p_child_state_map["vars_map"]["$i"];

    // what is the global number of iteratios executed relative to the root state
    p_state_map["iters_num_global_int"] = p_child_state_map["iters_num_global_int"];

    p_state_map["rules_iters_num_stack_lst"] = p_child_state_map["rules_iters_num_stack_lst"];

    //----------------------
    // ANIMATIONS - are not merged from children, they can only be propagated
    //              down the execution tree, not up.
    // 
    //----------------------

    return p_state_map;
}



//-------------------------------------------------
export function create_new(p_state_parent_map) {

    const state_map = get_empty_state();

    if (p_state_parent_map != null) {
        state_map["x"]  = p_state_parent_map["x"];
        state_map["y"]  = p_state_parent_map["y"];
        state_map["z"]  = p_state_parent_map["z"];
        state_map["rx"] = p_state_parent_map["rx"];
        state_map["ry"] = p_state_parent_map["ry"];
        state_map["rz"] = p_state_parent_map["rz"];
        state_map["sx"] = p_state_parent_map["sx"];
        state_map["sy"] = p_state_parent_map["sy"];
        state_map["sz"] = p_state_parent_map["sz"];
        state_map["cr"] = p_state_parent_map["cr"];
        state_map["cg"] = p_state_parent_map["cg"];
        state_map["cb"] = p_state_parent_map["cb"];
        state_map["iters_max"]             = p_state_parent_map["iters_max"];
        state_map["rules_names_stack_lst"] = p_state_parent_map["rules_names_stack_lst"];
        state_map["vars_map"]              = Object.assign({}, p_state_parent_map["vars_map"]); // clone

        state_map["iters_num_global_int"]      = p_state_parent_map["iters_num_global_int"];
        state_map["rules_iters_num_stack_lst"] = Object.assign([], p_state_parent_map["rules_iters_num_stack_lst"]); // clone
        state_map["animations_active_map"]     = Object.assign({}, p_state_parent_map["animations_active_map"]);     // clone
    }

    return state_map;
}

//-------------------------------------------------
function get_empty_state() {
    const state_map = {
        "x":  0.0,
        "y":  0.0,
        "z":  0.0,
        "rx": 0.0,
        "ry": 0.0,
        "rz": 0.0,
        "sx": 1.0,
        "sy": 1.0,
        "sz": 1.0,
        "cr": 0.0,
        "cg": 0.0,
        "cb": 0.0,

        // global max number of iterations for any rule
        "iters_max": 250,

        // list of all rules that are executing
        "rules_names_stack_lst": ["root"],

        "vars_map": {
            "$i": 0,
        },

        // global iterations number for a particular root expression
        "iters_num_global_int": 0,

        // stack of iteration numbers for each rule as its entered
        "rules_iters_num_stack_lst": [0],

        // ANIMATIONS - map of animations that are currently active
        //              in a subexpression or its children.
        "animations_active_map": {}
    };
    return state_map;
}