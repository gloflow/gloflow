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

import * as gf_engine     from "./engine/gf_engine";
import * as gf_lang_exec  from "./gf_lang_exec";
import * as gf_state      from "./gf_state";
import * as gf_lang_utils from "./gf_lang_utils";
import * as gf_ide        from "./ide/gf_ide";

//-------------------------------------------------
export function run(p_program_ast_lst) {

    //------------------------------------
    // AST_EXPANSION
    const expanded_program_ast_lst = [];
    for (let root_expression_lst of p_program_ast_lst) {
        const expanded_root_expression_lst = expand_tree(root_expression_lst, 0);

        // only include expressions which are not "expanded" to expression of 0 length
        // (expressions which are not marked for deletion)
        if (expanded_root_expression_lst.length > 0) {
            expanded_program_ast_lst.push(expanded_root_expression_lst);
        }
    }

    //------------------------------------
    // only load rules after AST tree expansion is complete,
    // and the rule is ready for execution.
    const rule_defs_map   = load_rule_defs(expanded_program_ast_lst);
    const shader_defs_map = load_shader_defs(expanded_program_ast_lst);

    // ENGINE
    const engine_api_map = gf_engine.init(shader_defs_map);

    // IDE
    gf_ide.init(engine_api_map);

    
    //------------------------------------
    // AST_EXECUTION

    var i=0;
    for (let root_expression_lst of expanded_program_ast_lst) {

        const root_state_map = gf_state.create_new(null);

        //------------------------------------
        // STATE_FAMILY_STACK
        // STATE_FAMILY - is a group of stacks that are related and are treated independently,
        //                without copying/merging back into the states from which it came from.
        //
        // allowing for multiple state objects to exist durring program execution.
        // initially there was only one blank state that the program started with, and all future
        // operations worked with that one state.
        // going forward there are state_setters that allow for pushing new states onto that stack,
        // and popping them.
        const state_family_stack_lst = [];
        // gf_state.push_family(root_state_map, state_family_stack_lst)
        
        //------------------------------------

        gf_lang_exec.execute_tree(root_expression_lst,
            root_state_map,
            rule_defs_map,
            shader_defs_map,
            state_family_stack_lst,
            engine_api_map);

        i+=1;
    }
    
    //------------------------------------
}

//-------------------------------------------------
function expand_tree(p_expression_ast_lst,
    p_tree_level_int :number) {

    var expression_lst = gf_lang_utils.clone_expr(p_expression_ast_lst); // clone in case of mutations of expression
    for (var i=0; i<expression_lst.length;) {

        const element = expression_lst[i];

        //------------------------------------
        // RULE_DEFINITION
        // [rule rule_name, [...]]
        if (element === "rule") {

            if (p_tree_level_int != 0) throw "rule definitions can only exist at the top level";
            if (i != 0) throw "rule definition has to be at the start of the expression";

            // rule definition can be of form:
            // ["rule", name_str, expressions_lst]
            // ["rule", name_str, rule_modifiers_lst, expressions_lst]
            if (expression_lst.length != 3 &&
                expression_lst.length != 4) throw "rule definition expression can only have 3|4 elements";
            if (expression_lst.length == 3) 
                if (!Array.isArray(expression_lst[2])) throw "rule definitions 3rd element has to be a list of its expressions";
            if (expression_lst.length == 4) 
                if (!Array.isArray(expression_lst[3])) throw "rule definitions 4rd element has to be a list of its expressions";

            // fast-forward to the 3rd|4th element of the expression, which represents rules expressions
            // so that tree expansion can be run on that rules expressions.
            if (expression_lst.length == 3) i+=2;
            if (expression_lst.length == 4) i+=3;
            continue;
        }

        //------------------------------------
        // SET
        else if (element == "set") {
            //if (p_tree_level_int != 0) throw "'set' statement can only exist at the top level";
            if (i != 0) throw "'set' declaration has to be at the start of the expression";

        }
        
        //------------------------------------
        // MULTIPLICATION__TOP_LEVEL
        // [* op1 op2]
        else if (element == "*" && i==0) {
            if (i != 0) throw "* operator has to be the first element in the expression";

            const operand_1 = expression_lst[i+1];
            const operand_2 = expression_lst[i+2];
            
            if (!(typeof operand_1 == 'number')) throw "first operand of multiplication expression is not a number";
            if (Array.isArray(operand_2)) {

                const expression_to_multiply_lst = operand_2;
                const factor_int                 = operand_1;
                const expanded_expressions_lst   = gf_lang_utils.clone_expr_N_times(expression_to_multiply_lst, factor_int)

                /*
                ["*", 10, [["y", -2.0], "cube"]], // 10 * {x -2} cube
                transforms to:

                [
                    [["y", -2.0], "cube"],
                    [["y", -2.0], "cube"],
                    [["y", -2.0], "cube"],
                ]
                */
                // multiplication of sub-expression has been projected,
                // and multiplication expression itself eliminated/replaced by new cloned expressions
                expression_lst = expanded_expressions_lst;

                i=0;      // rewind to start, since the expression has a new form, and should be re-processed
                continue; // go straight to new iteration, without incrementing 'i' (keeping it at 0 instead)
            }
        }

        //------------------------------------
        // SUB_EXPRESSION
        // [...]
        else if (Array.isArray(element)) {
            const sub_expression_lst = element;
            
            // RECURSION
            const expanded_sub_expression_lst = expand_tree(sub_expression_lst, p_tree_level_int+1);

            // IMPORTANT!! - splice the expanded sub-expression in the place of the old unexpanded element
            const insertion_index_int = i;
            expression_lst[i] = expanded_sub_expression_lst;
        }

        //------------------------------------
        else if (element == "lang_v") {
            if (i != 0) throw "lang_v expression identifier can only be the first element in the expression";

            // this expression is expanded to expression of 0 length, meaning it should be removed.
            return [];
        }

        //------------------------------------

        i+=1;
    }
    return expression_lst;
}

//-------------------------------------------------
function load_rule_defs(p_program_ast_lst) {
    const rule_defs_map = {};
    for (var i=0; i<p_program_ast_lst.length;) {

        const root_expression_lst = p_program_ast_lst[i];
        if (root_expression_lst[0] == "rule") {

            // rule with no modifiers
            if (root_expression_lst.length == 3) {
                const [, rule_name_str, rule_expressions_lst] = root_expression_lst;

                const rule_def_map = {
                    "modifiers_map":   {},
                    "expressions_lst": rule_expressions_lst,
                };
                if (Array.isArray(rule_defs_map[rule_name_str])) {
                    rule_defs_map[rule_name_str].push(rule_def_map);
                } else {
                    rule_defs_map[rule_name_str] = [rule_def_map];
                }
                

                // remove the rule definition element from the program_ast, 
                // as it has been expanded and loaded and ready for execution,
                // it doesnt need to be iterated over during execution.
                p_program_ast_lst.splice(i, 1);

                // run next iteration without incrementing "i"
                continue;
            }

            // rule with modifiers
            else if (root_expression_lst.length == 4) {
                const [, rule_name_str, rule_modifiers_lst, rule_expressions_lst] = root_expression_lst;

                // MODIFIERS
                const rule_modifiers_map = {}
                for (let modifier_lst of rule_modifiers_lst) {
                    const [modifier_name_str, modifier_val] = modifier_lst;
                    rule_modifiers_map[modifier_name_str] = modifier_val;
                }
                const rule_def_map = {
                    "modifiers_map":   rule_modifiers_map,
                    "expressions_lst": rule_expressions_lst,
                };
                if (Array.isArray(rule_defs_map[rule_name_str])) {
                    rule_defs_map[rule_name_str].push(rule_def_map);
                } else {
                    rule_defs_map[rule_name_str] = [rule_def_map];
                }

                // remove the rule definition element from the program_ast, 
                // as it has been expanded and loaded and ready for execution,
                // it doesnt need to be iterated over during execution.
                p_program_ast_lst.splice(i, 1);

                // run next iteration without incrementing "i"
                continue;
            }
        }

        i+=1;
    }
    return rule_defs_map;
}

//-------------------------------------------------
function load_shader_defs(p_program_ast_lst) {
    const shader_defs_map = {};
    for (var i=0; i<p_program_ast_lst.length;) {

        const root_expression_lst = p_program_ast_lst[i];
        if (root_expression_lst[0] == "shader") {

            if (root_expression_lst.length != 4 && root_expression_lst.length != 5)
                throw `shader definition expression ${root_expression_lst} can only have 4|5 elements`;
            
            if (root_expression_lst.length == 4) {
                const [, shader_name_str, vertex_shader_lst, fragment_shader_lst] = root_expression_lst;
                const [, vertex_code_str]   = vertex_shader_lst;
                const [, fragment_code_str] = fragment_shader_lst;

                shader_defs_map[shader_name_str] = {
                    "vertex_code_str":   vertex_code_str,
                    "fragment_code_str": fragment_code_str,
                };
            }
            if (root_expression_lst.length == 5) {

                const [, shader_name_str,
                    uniforms_defs_expr_lst,
                    vertex_shader_lst,
                    fragment_shader_lst] = root_expression_lst;
                const [, vertex_code_str]   = vertex_shader_lst;
                const [, fragment_code_str] = fragment_shader_lst;
                
                const [, uniforms_defs_lst] = uniforms_defs_expr_lst;

                shader_defs_map[shader_name_str] = {
                    "uniforms_defs_lst": uniforms_defs_lst,
                    "vertex_code_str":   vertex_code_str,
                    "fragment_code_str": fragment_code_str,
                };
            }

            // remove the rule definition element from the program_ast, 
            // as it has been expanded and loaded and ready for execution,
            // it doesnt need to be iterated over during execution.
            p_program_ast_lst.splice(i, 1);
        }

        i+=1;
    }
    return shader_defs_map
}