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

declare var THREE;

//-------------------------------------------------
export function create_all(p_shaders_defs_map) {

    const compiled_shaders_map = {};
    for (const shader_name_str in p_shaders_defs_map) {
        const shader_def_map   = p_shaders_defs_map[shader_name_str];
        
        const shader_material = create_from_def(shader_def_map);
        compiled_shaders_map[shader_name_str] = shader_material;
    }
    return compiled_shaders_map;
}

//-------------------------------------------------
export function create_from_def(p_shader_def_map) {

    const uniforms_defs_lst = p_shader_def_map["uniforms_defs_lst"];
    const vertex_code_str   = p_shader_def_map["vertex_code_str"];
    const fragment_code_str = p_shader_def_map["fragment_code_str"];

    const uniforms_map = {};
    for (const uniform_def_lst of uniforms_defs_lst) {
        // ["i", "float", 0.1]
        const [uniform_name_str, , default_val] = uniform_def_lst;

        uniforms_map[uniform_name_str] = {value: default_val};
    }

    const shader_material = create(uniforms_map, vertex_code_str, fragment_code_str);
    return shader_material;
}

//-------------------------------------------------
function create(p_uniforms_defs_map,
    p_vertex_shader_str   :string,
    p_fragment_shader_str :string) {

    const material = new THREE.ShaderMaterial({
        
        //---------------------
        // UNIFORMS
        uniforms: p_uniforms_defs_map,
        
        /*uniforms: {
            // time: {value: 1.0},
            // resolution: {value: new THREE.Vector2()}
        },*/
        
        //---------------------
        
        vertexShader:   p_vertex_shader_str,
        fragmentShader: p_fragment_shader_str
    } );

    return material;
}