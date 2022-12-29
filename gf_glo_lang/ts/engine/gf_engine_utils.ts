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
export function get_derived_coords(p_x, p_y, p_z, p_coord_origins_stack_lst) {

    // get global coords that are derived from the latest coordinate system origin
    // that was last added to the stack of coordinate system origins.
    const latest_coord_origin_v3 = p_coord_origins_stack_lst[p_coord_origins_stack_lst.length-1];
    const derived_x_f = latest_coord_origin_v3.x + p_x;
    const derived_y_f = latest_coord_origin_v3.y + p_y;
    const derived_z_f = latest_coord_origin_v3.z + p_z;

    return [derived_x_f, derived_y_f, derived_z_f];
}

//-------------------------------------------------
// export function get_derived_rotation(p_rx, p_ry, p_rz, p_rotations_stack_lst) {
//     const [latest_rx_f, latest_ry_f, latest_rz_f] = p_rotations_stack_lst[p_rotations_stack_lst.length-1];
//     const derived_rx_f = p_rx+latest_rx_f;
//     const derived_ry_f = p_ry+latest_ry_f;
//     const derived_rz_f = p_rz+latest_rz_f;
//     return [derived_rx_f, derived_ry_f, derived_rz_f];
// }

//-------------------------------------------------
export function get_real_world_pos(p_x, p_y, p_z,
    p_rx, p_ry, p_rz,
    p_rotation_pivot_point_stack_lst) {
    
    const real_world_point_v3 = new THREE.Vector3(p_x, p_y, p_z);

    // last pivot point in the stack, around which the rotation should be applied
    const pivot_point_v3 = get_rotation_pivot(p_rotation_pivot_point_stack_lst)

    const rotated_point_v3 = rotate_around_pivot(real_world_point_v3,
        p_rx, p_ry, p_rz,
        pivot_point_v3);
        
    return rotated_point_v3;
}

//-------------------------------------------------
// ROTATIONS
//-------------------------------------------------
export function get_rotation_pivot(p_rotation_pivot_point_stack_lst) {
    return p_rotation_pivot_point_stack_lst[p_rotation_pivot_point_stack_lst.length-1];
}

//-------------------------------------------------
export function rotate_around_pivot(p_point_v3,
    p_rx, p_ry, p_rz,
    p_pivot_point_v3) {

    const new_point_v3 = new THREE.Vector3(p_point_v3.x, p_point_v3.y, p_point_v3.z);

    // first subtract pivot_point position from obj position,
    // so that we're rotating around world origin again.
    new_point_v3.x -= p_pivot_point_v3.x; 
    new_point_v3.y -= p_pivot_point_v3.y;
    new_point_v3.z -= p_pivot_point_v3.z;

    // rotate around world origin
    const rotation_world_matrix_mat4 = rotate_world(p_rx, p_ry, p_rz);
    new_point_v3.applyMatrix4(rotation_world_matrix_mat4);

    // translate object back to its final position
    // by the povot_point delta
    new_point_v3.x += p_pivot_point_v3.x; 
    new_point_v3.y += p_pivot_point_v3.y;
    new_point_v3.z += p_pivot_point_v3.z;

    return new_point_v3;
}

//-------------------------------------------------
export function rotate_world(p_rx, p_ry, p_rz) {

    // WORLD_CENTERED_ROTATION
    // Calling "makeRotationX/Y/Z" will fill up the matrix cells with 
    // values for rotation around the wanted axis,
    // thus overriding any previous value in the matrix. thats why separate matrices
    // are created here for each axis, and then finally multiplied at the end.
    const rx   = new THREE.Matrix4().makeRotationX(p_rx);
    const ry   = new THREE.Matrix4().makeRotationY(p_ry);
    const rz   = new THREE.Matrix4().makeRotationZ(p_rz);
    const rxy  = new THREE.Matrix4().multiplyMatrices(rx, ry);
    const rxyz = new THREE.Matrix4().multiplyMatrices(rxy, rz);
    const rotation_matrix = rxyz;
    return rotation_matrix;
}

//-------------------------------------------------
export function rotate_self(p_mesh, p_rx, p_ry, p_rz) {
    // SELF_ROTATION
    p_mesh.rotation.x = p_rx;
    p_mesh.rotation.y = p_ry;
    p_mesh.rotation.z = p_rz;
}