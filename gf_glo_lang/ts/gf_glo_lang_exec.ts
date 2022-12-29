

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
import * as gf_state          from "./gf_state";
import * as gf_state_setters  from "./gf_state_setters";

//-------------------------------------------------
export function execute_tree(p_expression_ast_lst,
    p_state_parent_map,
    p_rules_defs_map,
    p_shader_defs_map,
    p_state_family_stack_lst,
    p_engine_api_map) {
    
    // console.log("A")

    const symb_map = gf_glo_lang_utils.get_symbols_and_constants();
    
    //--------------------
    // STATE_NEW
    // IMPORTANT!! - on every tree descent a new independent state is constructed
    var state_map = gf_state.create_new(p_state_parent_map);

    //--------------------

    // clone in case of mutations of expression
    const expression_lst = gf_glo_lang_utils.clone_expr(p_expression_ast_lst);

    // iterate over each expression element
    for (var i=0; i<expression_lst.length; i++) {

        const element = expression_lst[i];

        //------------------------------------
        // PROPERTY_MODIFIER
        if (symb_map["predefined_properties_lst"].includes(element)) {

            const property_name_str = element;
            const modifier          = expression_lst[i+1];
            var modifier_val;

            // modifier is a number
            if (typeof modifier == 'number') {
                modifier_val = modifier;
            }
            // modifier is a sub-expression - needs to be evaluated
            else if (Array.isArray(modifier)) {
                const sub_expr_lst = modifier;

                if (gf_glo_lang_utils.is_sys_func(sub_expr_lst)) {
                    const result = gf_glo_lang_utils.sys_func_eval(sub_expr_lst);
                    modifier_val = result;
                }
                else {
                    const result = gf_glo_lang_utils.arithmetic_eval(sub_expr_lst, state_map);
                    modifier_val = result;
                }
            }

            // IMPORTANT!! - incremental modification
            state_map[property_name_str] += modifier_val;

            // if (property_name_str.startsWith("r"))
            //     console.log("rotation", property_name_str, modifier_val, state_map[property_name_str])

            i+=1; // fast-forward, modifiers can be listed sequentially in the same expression
            continue;
        }

        //------------------------------------
        // SUB_EXPRESSION
        else if (Array.isArray(element)) {

            const sub_expr_lst = element;
            
            // recursion
            const [child_state_map, sub_expr_result] = execute_tree(sub_expr_lst,
                state_map,
                p_rules_defs_map,
                p_shader_defs_map,
                p_state_family_stack_lst,
                p_engine_api_map);
                

            

            const merged_state_map = gf_state.merge_child_state(state_map, child_state_map);
            state_map = merged_state_map;
            
            // sub-expression evaluated to a value
            if (sub_expr_result != null) {

                // substitute sub-expression for its results
                expression_lst[i] = sub_expr_result;
                
                // continue looping through the expression elements,
                // without incrementing "i". because we evaluated the expression
                // at position "i" and substituted results of that expression at that slot.
                continue;
            }
        }

        //------------------------------------
        // ARITHMETIC
        else if (Object.keys(symb_map["arithmetic_ops_map"]).includes(element) && i==0) {

            const arithmetic_result = gf_glo_lang_utils.arithmetic_eval(expression_lst, state_map)
            const expr_result       = arithmetic_result;
            return [state_map, expr_result];
        }
        
        //------------------------------------
        // CONDITIONALS
        else if (element == "if") {

            if (i != 0) throw "'if' keyword can ony be the first element in the expression";

            const child_state_map = expr__conditional(expression_lst,
                state_map,
                p_rules_defs_map,
                p_shader_defs_map,
                p_state_family_stack_lst,
                p_engine_api_map);
            
            const merged_state_map = gf_state.merge_child_state(state_map, child_state_map);
            state_map = merged_state_map;

            break;
        }

        //------------------------------------
        // STATE_SETTERS - global state setters
        else if (element == "set" || element == "push" || element == "pop") {
            
            const setter_type_str = element;
            const [, property_name_str, vals] = expression_lst;


            if ((setter_type_str == "push" || setter_type_str == "pop") && property_name_str == "coord_origin") {
                
                // coord_origin state_setter is the only setter so far that returns
                // a new state, all other state setters dont modify the gf_lang state.
                const new_state_map = gf_state_setters.exec__expr(setter_type_str,
                    property_name_str,
                    vals,
                    state_map,
                    p_state_family_stack_lst,
                    p_engine_api_map);

                state_map = new_state_map;
            }

            else {
                gf_state_setters.exec__expr(setter_type_str,
                    property_name_str,
                    vals,
                    state_map,
                    p_state_family_stack_lst,
                    p_engine_api_map);
            }
            break;
        }

        //------------------------------------
        // PRINT
        else if (element == "print") {

            expr__print(expression_lst, state_map);
            break;
        }

        //------------------------------------
        // ANIMATION
        else if (element == "animate") {
            if (i != 0) throw "'animate' keyword can ony be the first element in the expression";

            expr__animation(expression_lst, state_map, p_engine_api_map);
            break;
        }
 
        //------------------------------------
        // RULE_CALL
        else {

            const rule_name_str = element;

            if (i != expression_lst.length-1)
                throw `rule ${rule_name_str} call can only be the last element in expression`;

            const new_state_map = expr__rule_call(rule_name_str,
                expression_lst,
                state_map,
                p_rules_defs_map,
                p_shader_defs_map,
                p_state_family_stack_lst,
                p_engine_api_map);
            
            state_map = new_state_map;
            break;
        }

        //------------------------------------
    }

    return [state_map, null];
}

//-------------------------------------------------
function expr__rule_call(p_called_rule_name_str :string,
    p_expression_lst,
    p_state_parent_map,
    p_rules_defs_map,
    p_shader_defs_map,
    p_state_family_stack_lst,
    p_engine_api_map) {

    const symb_map = gf_glo_lang_utils.get_symbols_and_constants();
    //------------------------------------
    // SYSTEM_RULE
    // rules predefined in the system

    if (symb_map["system_rules_lst"].includes(p_called_rule_name_str)) {

        const new_state_map = expr__rule_sys_call(p_called_rule_name_str,
            p_state_parent_map,
            p_engine_api_map);
        return new_state_map;
    }

    //------------------------------------
    // USER_RULE
    // rules defined by the user in their program

    else if (Object.keys(p_rules_defs_map).includes(p_called_rule_name_str)) {


        //--------------------
        // STATE_NEW
        // for each rule invocation a new state object is created, that inherits 
        // the values of its parent state, within the same state family.

        const new_state_map = gf_state.create_new(p_state_parent_map);

        //--------------------
        
        const current_rule_name_str = gf_glo_lang_utils.rule_get_name(p_state_parent_map);

        // all user calls (even recursive) are stored in the call stack
        new_state_map["rules_names_stack_lst"].push(p_called_rule_name_str);

        // console.log(`calling rule ${current_rule_name_str}->${p_called_rule_name_str}`, p_state_parent_map["rules_names_stack_lst"]);


        // new rule getting executed
        if (current_rule_name_str !== p_called_rule_name_str) {

            // start a new iterations counter since we're entering a new rule
            // (not recursively iterating within the same rule)
            gf_glo_lang_utils.add_new_iters_num_state(new_state_map);
        }

        const current_rule_iters_num_int = gf_glo_lang_utils.rule_get_iters_num(new_state_map);



        // pick a random definition for a rule, which can have many definitions.
        const [rule_def_map, rule_expressions_lst] = gf_glo_lang_utils.pick_rule_random_def(p_called_rule_name_str,
            p_rules_defs_map);

        //------------------------------------
        // RECURSION_STOP
        //  - prevent infinite rules execution
        //  - if global iters_max limit is reached (global for all rules)
        //  - if local rule-specific (rule modifier) iters_max limit is reached (for current rule only)

        // GLOBAL_LIMIT
        if (new_state_map["iters_num_global_int"] > new_state_map["iters_max"]-1) {

            console.log("global iter limit reached");
            return new_state_map;
        }

        // RULE_LIMIT
        // check if rule has a iters_max rule modifier specified
        else if (Object.keys(rule_def_map["modifiers_map"]).includes("iters_max")) {

            // const rule_iters_num_int   = gf_glo_lang_utils.rule_get_iters_num(p_state_map);
            const rule_iters_limit_int = rule_def_map["modifiers_map"]["iters_max"];

            if (current_rule_iters_num_int > rule_iters_limit_int-1) {

                // console.log(`local iter limit ${current_rule_iters_num_int} for rule ${p_called_rule_name_str} reached`,
                //    p_state_parent_map["rules_names_stack_lst"]);

                //-----------------
                // RULE_EXIT
                // rule naturally ended with iter_num limit, without entering
                // a different rule. the state has to be reset to the callers
                // state.
                new_state_map["rules_iters_num_stack_lst"].pop();
                new_state_map["rules_names_stack_lst"].pop();

                const old_rule_iters_num_int    = gf_glo_lang_utils.rule_get_iters_num(new_state_map);
                new_state_map["vars_map"]["$i"] = old_rule_iters_num_int;
                
                //-----------------

                return new_state_map;
            }
        }

        //------------------------------------
        
        
        

        // RECURSION
        // IMPORTANT!! - rules are not yet treated as expressions, and cant return
        //               results of evaluating its expressions. so ignoring execution results.
        const [child_state_map, ] = execute_tree(rule_expressions_lst,
            new_state_map,
            p_rules_defs_map,
            p_shader_defs_map,
            p_state_family_stack_lst,
            p_engine_api_map);
        
        // remove rule_name from the stack of rules that were executed
        new_state_map["rules_names_stack_lst"].pop();

        if (current_rule_name_str !== p_called_rule_name_str) {

            // RULE_EXIT
            // we returned from a new rule into the old rule context,
            // so the iterations count for that new rule is no longer needed (and removed from stack).
            gf_glo_lang_utils.restore_previous_rules_iters_num(new_state_map);

            return new_state_map;
        }

        // we're still within the same rule, in one of its iterations, so just merge state
        else {

            const merged_state_map = gf_state.merge_child_state(new_state_map, child_state_map);
            return merged_state_map;
        }
    }

    //------------------------------------
    else {
        throw `rule call referencing an unexisting rule - ${p_called_rule_name_str}`;
    }

    //------------------------------------
}

//-------------------------------------------------
function expr__rule_sys_call(p_rule_name_str :string,
    p_state_map,
    p_engine_api_map) {

    //----------------------
    p_state_map["iters_num_global_int"] += 1;

    // IMPORTANT!! - rule iterations are counted only for actual rule evaluations.
    //               important not to count "set" statements, expression tree
    //               descending, property modifiers execution, etc.
    gf_glo_lang_utils.increment_iters_num(p_state_map);

    //----------------------

    // CUBE
    if (p_rule_name_str == "cube") {
        const x = p_state_map["x"];
        const y = p_state_map["y"];
        const z = p_state_map["z"];
        const rx = p_state_map["rx"];
        const ry = p_state_map["ry"];
        const rz = p_state_map["rz"];
        const sx = p_state_map["sx"];
        const sy = p_state_map["sy"];
        const sz = p_state_map["sz"];

        const cr = p_state_map["cr"];
        const cg = p_state_map["cg"];
        const cb = p_state_map["cb"];

        const create_cube_fun = p_engine_api_map["create_cube_fun"];
        create_cube_fun(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb);
    }

    // SPHERE
    else if (p_rule_name_str == "sphere") {
        const x = p_state_map["x"];
        const y = p_state_map["y"];
        const z = p_state_map["z"];
        const rx = p_state_map["rx"];
        const ry = p_state_map["ry"];
        const rz = p_state_map["rz"];
        const sx = p_state_map["sx"];
        const sy = p_state_map["sy"];
        const sz = p_state_map["sz"];

        const cr = p_state_map["cr"];
        const cg = p_state_map["cg"];
        const cb = p_state_map["cb"];

        const create_sphere_fun = p_engine_api_map["create_sphere_fun"];
        create_sphere_fun(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb);
    }

    // LINE
    else if (p_rule_name_str == "line") {
        const x = p_state_map["x"];
        const y = p_state_map["y"];
        const z = p_state_map["z"];
        const rx = p_state_map["rx"];
        const ry = p_state_map["ry"];
        const rz = p_state_map["rz"];
        const sx = p_state_map["sx"];
        const sy = p_state_map["sy"];
        const sz = p_state_map["sz"];

        const cr = p_state_map["cr"];
        const cg = p_state_map["cg"];
        const cb = p_state_map["cb"];

        const create_line_fun = p_engine_api_map["create_line_fun"];
        create_line_fun(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb);
    }

    return p_state_map;
}

//-------------------------------------------------
function expr__animation(p_expression_lst,
    p_state_map,
    p_engine_api_map) {

    const symb_map = gf_glo_lang_utils.get_symbols_and_constants();

    if (p_expression_lst.length != 3 && p_expression_lst.length != 4)
        throw "animation expression can only have 3|4 elements";

        
    var props_lst;
    var duration_sec_f;
    var repeat_bool = false;
    if (p_expression_lst.length == 3) {
        [, props_lst, duration_sec_f] = p_expression_lst;
    }
    else if (p_expression_lst.length == 4) {
        var repeat_str;
        [, props_lst, duration_sec_f, repeat_str] = p_expression_lst;
        if (repeat_str == "repeat") {
            repeat_bool = true;
        }
        else {
            throw "animation can only be enabled with the 'repeat' keyword";
        }
    }
    
    const props_to_animate_lst = [];

    for (const [prop_name_str, change_delta_f] of props_lst) {

        if (!symb_map["predefined_properties_lst"].includes(prop_name_str))
            throw `cant animate property that is not predefined - ${prop_name_str}`;

        
        const start_val_f = p_state_map[prop_name_str];
        const end_val_f   = start_val_f + change_delta_f;
        
        props_to_animate_lst.push({
            "name_str":    prop_name_str,
            "start_val_f": start_val_f,
            "end_val_f":   end_val_f,
        });
    }

    const animate_fun = p_engine_api_map["animate_fun"];
    animate_fun(props_to_animate_lst, duration_sec_f, repeat_bool);
}

//-------------------------------------------------
// EXPRESSION__CONDITIONAL
function expr__conditional(p_expression_lst,
    p_state_map,
    p_rules_defs_map,
    p_shader_defs_map,
    p_state_family_stack_lst,
    p_engine_api_map) {


    const [, condition_lst, sub_expressions_lst] = p_expression_lst;

    if (condition_lst.length > 3)
        throw "'if' condition can only have 3 elements [logic_op, operand1, operand2]";

    //-------------------------------------------------
    function evaluate_logic_expr(p_logic_expr_lst) :boolean  {
        
        const symb_map = gf_glo_lang_utils.get_symbols_and_constants();

        const [logic_op_str, operand_1, operand_2] = condition_lst;

        if (!Object.keys(symb_map["logic_operators_map"]).includes(logic_op_str))
            throw `specified logic operator ${logic_op_str}is not valid`;

        //-------------------------------------------------
        const op_1_val = gf_glo_lang_utils.expr_eval(operand_1, p_state_map);
        const op_2_val = gf_glo_lang_utils.expr_eval(operand_2, p_state_map);

        if (symb_map["logic_operators_map"][logic_op_str].eval(op_1_val, op_2_val)) {
            return true;
        } else {
            return false;
        }
    }
    
    //-------------------------------------------------

    // if condition evaluates to true, execute subexpressions
    if (evaluate_logic_expr(condition_lst)) {

        // recursion
        const [child_state_map, ] = execute_tree(sub_expressions_lst,
            p_state_map,
            p_rules_defs_map,
            p_shader_defs_map,
            p_state_family_stack_lst,
            p_engine_api_map);

        const merged_state_map = gf_state.merge_child_state(p_state_map, child_state_map);
        return merged_state_map;
    } else {
        return p_state_map; // else returned state unchanged
    }
}

//-------------------------------------------------
// EXPRESSION__PRINT
function expr__print(p_expression_lst,
    p_state_map) {
    const [, vals_lst] = p_expression_lst;

    var vals_str="";
    for (let val of vals_lst) {
        if (val.startsWith("$")) {
            const var_ref_str = val;
            const var_val     = gf_glo_lang_utils.var_eval(var_ref_str, p_state_map); // p_state_map["vars_map"][var_ref_str];
            const val_str     = `${var_ref_str}=${var_val}`;
            vals_str += val_str+" ";
        } 
        else {
            const val_str = val;
            vals_str += val_str+" ";
        }
    }

    console.log(`gf %c${vals_str}`, 'background: #222; color: #ffffff');
}