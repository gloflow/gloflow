
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

//-------------------------------------------------
export function expr_eval(p_expr, p_state_map) {

    const symb_map = get_symbols_and_constants();

    // NUMBER
    if (typeof p_expr == "number") {
        return p_expr;
    }

    // VAR_REFERENCE
    else if (typeof p_expr == "string" && p_expr.startsWith("$")) {

        // for now only system-defined vars are available, no user-defined vars yet.
        if (!symb_map["system_vars_lst"].includes(p_expr))
            throw `variable operand ${p_expr} is not one of the system defined vars ${symb_map["system_vars_lst"]}`;

        const var_value = var_eval(p_expr, p_state_map);
        return var_value;
    }

    // SUB_EXPRESSION
    else if (Array.isArray(p_expr)) {

        // ARITHMETIC_OPERATION
        if (Object.keys(symb_map["arithmetic_ops_map"]).includes(p_expr[0])) {
            const val = arithmetic_eval(p_expr, p_state_map);
            return val;
        }

        // SYSTEM_FUNCTION
        else if (is_sys_func(p_expr)) {
            const val = sys_func_eval(p_expr);
            return val;
        }
    }
}

//-------------------------------------------------
// ARITHMETIC_EVALUATION
export function arithmetic_eval(p_expr_lst, p_state_map) {

    const symb_map             = get_symbols_and_constants();
    const arithmetic_ops_map   = symb_map["arithmetic_ops_map"];
    const system_functions_lst = symb_map["system_functions_lst"];

    if (p_expr_lst.length != 3) throw `arithmetic expression ${p_expr_lst} has to be of length 3`;
    if (!Object.keys(arithmetic_ops_map).includes(p_expr_lst[0]))
        throw `arithmetic op ${p_expr_lst[0]} is not supported`;

    const [op_str, operand_1, operand_2] = p_expr_lst;

    const op_1 = eval_op(operand_1);
    const op_2 = eval_op(operand_2);

    // EVALUATE
    const result = arithmetic_ops_map[op_str].eval(op_1, op_2);
    return result;

    //-------------------------------------------------
    function eval_op(p_operand) {
        
        var operand;

        // SUB_EXPRESSION
        if (Array.isArray(p_operand)) {
            const sub_expr_lst = p_operand;

            // system_function sub-expression
            if (is_sys_func(sub_expr_lst)) {
                const sub_result = sys_func_eval(sub_expr_lst);
                operand = sub_result;
            }
            // arithmetic sub-expression
            else {
                const sub_result = arithmetic_eval(sub_expr_lst, p_state_map);
                operand = sub_result;
            }
        }
        // VARIABLE
        else if (typeof p_operand == "string" && p_operand.startsWith("$")) {
            const var_val = var_eval(p_operand, p_state_map);
            operand = var_val;
        }
        
        // NUMBER
        else {
            // if operand is not a var reference, it has to be a number
            if (typeof p_operand != "number") throw `operand ${p_operand} is not a number`;
            operand = p_operand;
        }
        return operand;
    }

    //-------------------------------------------------
}

//-------------------------------------------------
// SYSTEM_FUNCTIONS
export function is_sys_func(p_expr_lst) {
    const sys_funs_lst = get_symbols_and_constants()["system_functions_lst"];
    if (sys_funs_lst.includes(p_expr_lst[0])) {
        return true;
    } else {
        return false;
    }
}
export function sys_func_eval(p_expr_lst) {
    
    const [func_name_str, args_lst] = p_expr_lst;
    
    var val;
    if (func_name_str == "rand") {
        if (args_lst.length != 2) throw "'rand' system function only takes 2 argument";
        const random_range_min_f = args_lst[0];
        const random_range_max_f = args_lst[1];
        val = Math.random()*(random_range_max_f-random_range_min_f) + random_range_min_f;
    }
    return val;
}

//-------------------------------------------------
export function var_eval(p_var_str, p_state_map) {
    if (!p_var_str.startsWith("$")) throw `variable string ${p_var_str} has no "$" prefixed`;
    const var_value = p_state_map["vars_map"][p_var_str];
    return var_value;
}

//-------------------------------------------------
export function rule_get_iters_num(p_state_map) {
    return p_state_map["rules_iters_num_stack_lst"][p_state_map["rules_iters_num_stack_lst"].length-1];
}

export function rule_get_name(p_state_map) {
    const rule_name_str = p_state_map["rules_names_stack_lst"][p_state_map["rules_names_stack_lst"].length-1];
    if (rule_name_str == undefined) throw "rule name is undefined";
    return rule_name_str;
}

//-------------------------------------------------
export function increment_iters_num(p_state_map) :number {
    const new_rule_iters_num_int = rule_get_iters_num(p_state_map) + 1;

    p_state_map["rules_iters_num_stack_lst"][p_state_map["rules_iters_num_stack_lst"].length-1] = new_rule_iters_num_int;
    p_state_map["vars_map"]["$i"] = new_rule_iters_num_int;
    
    return new_rule_iters_num_int;
}

export function add_new_iters_num_state(p_state_map) {
    p_state_map["rules_iters_num_stack_lst"].push(0);
    p_state_map["vars_map"]["$i"] = 0;
}

//-------------------------------------------------
// called when one rule exits (finishes executing) and returns execution
// to its parent rule (not when the same rule recurses into itself).
export function restore_previous_rules_iters_num(p_state_map) {
    p_state_map["rules_iters_num_stack_lst"].pop();

    // reinitialize $i to the parents number of iterations
    p_state_map["vars_map"]["$i"] = p_state_map["rules_iters_num_stack_lst"][p_state_map["rules_iters_num_stack_lst"].length-1];
}

//-------------------------------------------------
export function pick_rule_random_def(p_rule_name_str,
    p_rules_defs_map) {

    const rule_defs_lst             = p_rules_defs_map[p_rule_name_str];
    const rule_defs_num_int         = rule_defs_lst.length;
    const random_rule_def_index_int = Math.floor(Math.random() * rule_defs_num_int);
    const rule_def_map              = rule_defs_lst[random_rule_def_index_int];
    const rule_expressions_lst      = rule_def_map["expressions_lst"];
    return [rule_def_map, rule_expressions_lst];
}

//-------------------------------------------------
export function clone_expr(p_expr_lst) {
    return Object.assign([], p_expr_lst);
}

//-------------------------------------------------
export function clone_expr_N_times(p_expr_lst, p_n_int) {
    const clones_lst = [];
    for (var i=0;i<p_n_int;i++) {
        const cloned_expr_lst = clone_expr(p_expr_lst);
        clones_lst.push(cloned_expr_lst);
    }
    return clones_lst;
}

//-------------------------------------------------
export function get_symbols_and_constants() {
    const rule_level_max_int = 250;
    const system_rules_lst = [
        "cube",
        "sphere",
        "line"
    ];
    const predefined_properties_lst = [
        "x",  // x-coordinate
        "y",  // y-coordinate
        "z",  // z-coordinate
        "rx", // x-rotation
        "ry", // y-rotation
        "rz", // z-rotation
        "sx", // x-scale
        "sy", // y-scale
        "sz", // z-scale
        "cr", // red-channel-color
        "cg", // green-channel-color
        "cb"  // blue-channel-color
    ];
    const logic_operators_map = {
        "==": {eval: (p1, p2)=>{return p1==p2}},
        "!=": {eval: (p1, p2)=>{return p1!=p2}},
        "<":  {eval: (p1, p2)=>{return p1<p2}},
        ">":  {eval: (p1, p2)=>{return p1>p2}},
        "<=": {eval: (p1, p2)=>{return p1<=p2}},
        ">=": {eval: (p1, p2)=>{return p1>=p2}},
    };
    const arithmetic_ops_map = {
        "+": {eval: (p1, p2)=>{return p1+p2}},
        "-": {eval: (p1, p2)=>{return p1-p2}},
        "*": {eval: (p1, p2)=>{return p1*p2}},
        "/": {eval: (p1, p2)=>{return p1/p2}},
        "%": {eval: (p1, p2)=>{return p1%p2}},
    }
    const system_vars_lst = [
        "$i", // current rule iteration
    ];
    const system_functions_lst = [
        "rand", // random number generator
    ]

    return {
        "rule_level_max_int":        rule_level_max_int,
        "system_rules_lst":          system_rules_lst,
        "predefined_properties_lst": predefined_properties_lst,
        "logic_operators_map":       logic_operators_map,
        "arithmetic_ops_map":        arithmetic_ops_map,
        "system_vars_lst":           system_vars_lst,
        "system_functions_lst":      system_functions_lst,
    }
}