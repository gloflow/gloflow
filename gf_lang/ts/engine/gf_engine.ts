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

import * as gf_shaders      from "./gf_shaders";
import * as gf_engine_utils from "./gf_engine_utils";
import * as gf_animations   from "./gf_animations";

declare var THREE;
declare var TWEEN;

//-------------------------------------------------
export function init(p_shader_defs_map) {

	// SHADERS
	const compiled_shaders_map = gf_shaders.create_all(p_shader_defs_map);


    const scene_3d = new THREE.Scene();
    
	// CAMERA
	const render_dinstance_int = 10000;
	const camera = new THREE.PerspectiveCamera(75,
		window.innerWidth / window.innerHeight,
		0.1,
		render_dinstance_int);
	camera.position.y = 10;
	camera.position.z = 10;
	camera.position.x = 10;
    camera.lookAt(new THREE.Vector3(0, 5000, 0));

    // RENDERER
	const renderer = new THREE.WebGLRenderer({
		antialias: true
	});
	renderer.setSize( window.innerWidth, window.innerHeight );

	// const bg_color = new THREE.Color(0.9, 0.9, 0.9); // 0x5f2f00;
	// renderer.setClearColor(bg_color, 1);

	document.body.appendChild(renderer.domElement);


	// TRACKING
	// IMPORTANT!! - pass in renderer.domElement instead of nothing. 
	//               otherwise Controls event-handlers will be registered on the window
	//               object and prevent keyboard event listeners from firing.
	//               also any clisk on IDE UI elements outside the canvas will also 
	//               trigger controls logic.
    const controls = new THREE.TrackballControls(camera, renderer.domElement);
    controls.rotateSpeed = 1.0;
	controls.zoomSpeed   = 1.2;
	controls.panSpeed    = 0.8;
	controls.noZoom      = false;
	controls.noPan       = false;
	//controls.staticMoving = true;
	controls.dynamicDampingFactor = 0.3;
	
	// controls.keys = [ 65, 83, 68 ]; //a/s/d
    // controls.addEventListener( 'change', render );

	/*var controls = new THREE.FlyControls(camera, renderer.domElement);
    controls.dragToLook = true;
    controls.movementSpeed = 20;
    controls.rollSpeed = 0.2;*/

	//-------------------------------------------------
	var lt = new Date();
	function render(p_time) {

		/*var now = new Date(),
        secs = (now.getTime() - lt.getTime()) / 1000;
        lt = now;

		controls.update(1 * secs);*/
		controls.update();

		requestAnimationFrame(render);
		renderer.render(scene_3d, camera);


		TWEEN.update(p_time);
	}

	//-------------------------------------------------
	requestAnimationFrame(render);

    
    const cube_geometry = new THREE.BoxGeometry(1, 1, 1);
    const cube_geo      = new THREE.EdgesGeometry(cube_geometry);
	const cube_mat      = new THREE.LineBasicMaterial({color: 0x000000, linewidth: 1.0});


	const radius_int          = 10;
	const width_segments_int  = 50;
	const height_segments_int = 20;
	const sphere_geometry = new THREE.SphereGeometry(radius_int, width_segments_int, height_segments_int);
    const sphere_geo      = new THREE.EdgesGeometry(sphere_geometry);
	const sphere_mat      = new THREE.LineBasicMaterial({color: 0x000000, linewidth: 1.0});

	const line_mat = new THREE.LineBasicMaterial({
		color: 0x000000
	});

	//-------------------------------------------------
	// CREATE_LINE

	var line_points_lst = [];
	function create_line(p_x :number, p_y :number, p_z :number,
		p_rx :number, p_ry :number, p_rz :number,
		p_sx :number, p_sy :number, p_sz :number,
		p_cr :number, p_cg :number, p_cb :number,
		p_rotation_pivot_points_stack_lst) {
		
		//--------------------
		// REAL_WORLD_POSITION

		const real_world_point_v3 = gf_engine_utils.get_real_world_pos(p_x, p_y, p_z,
			p_rx, p_ry, p_rz,
			p_rotation_pivot_points_stack_lst);

		//--------------------

		// a line point already exists, from a previous invocation of create_line()
		// and so with the new point we have the required number of points (2)
		// to draw a line
		if (line_points_lst.length > 0) {
			
			line_points_lst.push(real_world_point_v3); // new_point_v3);
			const geometry = new THREE.BufferGeometry().setFromPoints(line_points_lst);
			const line     = new THREE.Line(geometry, line_mat);

			scene_3d.add(line);
		}
		
		const old_point_v3 = real_world_point_v3; // new_point_v3;
		line_points_lst    = [old_point_v3];
	}

	//-------------------------------------------------
	// CREATE_CUBE

	function create_cube(p_x :number, p_y :number, p_z :number,
		p_rx :number, p_ry :number, p_rz :number,
		p_sx :number, p_sy :number, p_sz :number,
		p_cr :number, p_cg :number, p_cb :number,
		p_material,
		p_wireframe_material_bool,
		p_rotation_pivot_points_stack_lst) {
		
		const cube_mesh = new THREE.Mesh(cube_geometry, p_material);

		/*// APPLY_POSITION
		cube_mesh.position.x = p_x; 
		cube_mesh.position.y = p_y;
		cube_mesh.position.z = p_z;*/
		
		// APPLY_SCALE
		cube_mesh.scale.set(p_sx, p_sy, p_sz );
		
		//--------------------
		// TRANSLATION/ROTATION
		gf_engine_utils.rotate_self(cube_mesh, p_rx, p_ry, p_rz);
		

		const real_world_point_v3 = gf_engine_utils.get_real_world_pos(p_x, p_y, p_z,
			p_rx, p_ry, p_rz,
			p_rotation_pivot_points_stack_lst);

		cube_mesh.position.x = real_world_point_v3.x;
		cube_mesh.position.y = real_world_point_v3.y;
		cube_mesh.position.z = real_world_point_v3.z;

		//--------------------
		// WIREFRAME
		if (p_wireframe_material_bool) {
			const wireframe = new THREE.LineSegments(cube_geo, cube_mat);
			cube_mesh.add(wireframe);
		}

		scene_3d.add(cube_mesh);
		return cube_mesh;
	}

	//-------------------------------------------------
	// CREATE_SPHERE
	
	function create_sphere(p_x :number, p_y :number, p_z :number,
		p_rx :number, p_ry :number, p_rz :number,
		p_sx :number, p_sy :number, p_sz :number,
		p_cr :number, p_cg :number, p_cb :number,
		p_material,
		p_wireframe_material_bool) {

		material = new THREE.MeshBasicMaterial({
			color:               new THREE.Color(p_cr, p_cg, p_cb), // color_int,
			polygonOffset:       true,
			polygonOffsetFactor: 1, // positive value pushes polygon further away
			polygonOffsetUnits:  1
		});

		const sphere_mesh = new THREE.Mesh(sphere_geometry, p_material);
		sphere_mesh.position.x = p_x; 
		sphere_mesh.position.y = p_y;
		sphere_mesh.position.z = p_z;
		
		sphere_mesh.scale.set(p_sx, p_sy, p_sz );

		// ROTATION
		gf_engine_utils.rotate_self(sphere_mesh, p_rx, p_ry, p_rz);
		const rotation_world_matrix_mat4 = gf_engine_utils.rotate_world(p_rx, p_ry, p_rz);
		sphere_mesh.position.applyMatrix4(rotation_world_matrix_mat4);

		// WIREFRAME
		if (p_wireframe_material_bool) {
			const wireframe = new THREE.LineSegments(sphere_geo, sphere_mat);
			sphere_mesh.add(wireframe);
		}

		scene_3d.add(sphere_mesh);
	}

	//-------------------------------------------------

	//------------------------------------
	// ENGINE_STATE
	var material;
	var material_shader_active;
	var shader_material_bool    = false;
	var wireframe_material_bool = false;
	var animation_active_map    = null;

	var color_background;

	const coord_origins_stack_lst = [
		new THREE.Vector3(0, 0, 0)
	];

	const rotations_stack_lst = [
		[0.0, 0.0, 0.0]
	];

	// by default rotate objs around world origin
	const rotation_pivot_points_stack_lst = [
		new THREE.Vector3(0, 0, 0)
	];

	//------------------------------------

	const api_map = {

		"scene_3d":                scene_3d,
		"coord_origins_stack_lst": coord_origins_stack_lst,

		//-------------------------------------------------
		"get_state_fun": (p_state_prop_name_str)=>{
			if (p_state_prop_name_str == "color_background") {
				return color_background.toArray();
			}
		},

		//-------------------------------------------------
		"set_state_fun": (p_state_change_map)=>{
			
			//------------------------------------
			// COORDINATE_ORIGIN
			if (Object.keys(p_state_change_map).includes("property_name_str") && 
				p_state_change_map["property_name_str"] == "coord_origin") {

				const origin_type_str = p_state_change_map["origin_type_str"];

				if (origin_type_str == "current_pos") {
					switch (p_state_change_map["setter_type_str"]) {

						//------------------------------------
						// PUSH
						case "push":

							const latest_coord_origin_v3 = coord_origins_stack_lst[coord_origins_stack_lst.length-1];

							// get absolute world position of the state at the time of the
							// state change, so that it can be set as the new world origin.
							const real_world_point_v3 = gf_engine_utils.get_real_world_pos(latest_coord_origin_v3.x + p_state_change_map["x"],
								latest_coord_origin_v3.y + p_state_change_map["y"],
								latest_coord_origin_v3.z + p_state_change_map["z"],

								p_state_change_map["rx"],
								p_state_change_map["ry"],
								p_state_change_map["rz"],

								// stack is used to get the latest coordinate system origin,
								// since the supplied state_change x/y/z coords are relative to that
								// latest coordinate system.
								// coord_origins_stack_lst,

								// stack is used to get the latest pivot point to apply for rotation
								rotation_pivot_points_stack_lst);

							console.log("NEW COORD_ORIGIN", real_world_point_v3)
							
							const new_coord_origin_point_v3 = real_world_point_v3;
							coord_origins_stack_lst.push(new_coord_origin_point_v3);
							rotations_stack_lst.push([p_state_change_map["rx"], p_state_change_map["ry"], p_state_change_map["rz"]]);

							break;
						
						//------------------------------------
						// POP
						case "pop":

							coord_origins_stack_lst.pop();
							rotations_stack_lst.pop();

							break;

						//------------------------------------
					}
				}
			}

			//------------------------------------
			// ROTATION_PIVOT_POINT
			if (Object.keys(p_state_change_map).includes("property_name_str") && 
				p_state_change_map["property_name_str"] == "rotation_pivot") {

				switch (p_state_change_map["setter_type_str"]) {

					// PUSH
					case "push":
						
						var new_pivot_point_v3;
						if (p_state_change_map["axis_type_str"] == "current_pos") {

							// gf_lang state transforms/rotations have to be applied to get the final
							// point coords of this latest pivot point.
							// p_state_change_map["x"]/y/z only are not enough, since they're not what gets applied
							// to objects directly, instead rotations also get applied to objects
							// (and objects are often the povit points that are used)
							const real_world_point_v3 = gf_engine_utils.get_real_world_pos(p_state_change_map["x"],
								p_state_change_map["y"],
								p_state_change_map["z"],
								p_state_change_map["rx"],
								p_state_change_map["ry"],
								p_state_change_map["rz"],

								// stack is used to get the latest pivot point to apply for rotation
								rotation_pivot_points_stack_lst);

							new_pivot_point_v3 = real_world_point_v3;
						}

						

						// state transforms are applied before the latest pivot point is pushed
						// onto the stack.
						rotation_pivot_points_stack_lst.push(new_pivot_point_v3);
						break;
					
					// POP
					case "pop":
						rotation_pivot_points_stack_lst.pop();
						break;
				}
			}

			//------------------------------------
			// COLOR
			if (Object.keys(p_state_change_map).includes("color_rgb")) {
				
				const color = p_state_change_map["color_rgb"];

				var three_color;

				// RGB_LIST
				if (Array.isArray(color)) {
					const color_lst = color;
					const [r, g, b] = color_lst;
					three_color = new THREE.Color(r, g, b);

				// RGB_HEX
				} else {
					// parse color hex string
					three_color = new THREE.Color(color)
				}
				
				material = new THREE.MeshBasicMaterial({
					color:               three_color,
					polygonOffset:       true,
					polygonOffsetFactor: 1, // positive value pushes polygon further away
					polygonOffsetUnits:  1
				});

				return [three_color.r, three_color.g, three_color.b];
			}

			// COLOR_BACKGROUND
			if (Object.keys(p_state_change_map).includes("color_background")) {
				const color = p_state_change_map["color_background"];
				
				var bg_color;
				if (Array.isArray(color)) {
					const [r, g, b] = color;
					bg_color = new THREE.Color(r, g, b);
				}
				else {
					bg_color = new THREE.Color(color);
				}
				
				color_background = bg_color;
				// renderer.setClearColor(bg_color, 1);
				scene_3d.background = bg_color;
			}

			//------------------------------------
			// MATERIAL
			if (Object.keys(p_state_change_map).includes("material_type_str")) {

				const material_type_str = p_state_change_map["material_type_str"];

				// WIREFRAME
				if (material_type_str == "wireframe") {

					const material_value_bool = p_state_change_map["material_value_bool"];
					wireframe_material_bool = material_value_bool;
				}

				// SHADER
				else if (material_type_str == "shader") {
					const material_value_str = p_state_change_map["material_value_str"];
					const shader_name_str    = material_value_str;
					const shader_material    = compiled_shaders_map[shader_name_str];

					material_shader_active = shader_material;
					shader_material_bool   = true;
				}
			}

			//------------------------------------
			// MATERIAL_PROPERTY

			/*
			const state_change_map = {
				"material_prop_map": {
					"material_shader_name_str":         material_name_str,
					"material_shader_uniform_name_str": uniform_name_str,
					"material_shader_uniform_val":      loaded_val,
				}
			};
			*/

			if (Object.keys(p_state_change_map).includes("material_prop_map")) {
				const material_prop_change_map = p_state_change_map["material_prop_map"];

				if (Object.keys(material_prop_change_map).includes("material_shader_name_str")) {
					const material_shader_name_str         = material_prop_change_map["material_shader_name_str"];
					const material_shader_uniform_name_str = material_prop_change_map["material_shader_uniform_name_str"];
					const material_shader_uniform_val      = material_prop_change_map["material_shader_uniform_val"];

					const material = material_shader_active;

					material.uniforms[material_shader_uniform_name_str].value = material_shader_uniform_val;
				}
			}

			//------------------------------------
			// LINE_START
			if (Object.keys(p_state_change_map).includes("line_cmd_str")) {

				if (p_state_change_map["line_cmd_str"] == "start") {
					line_points_lst = []; // reset line memory
				}
			}

			//------------------------------------
		},

		//-------------------------------------------------
		// CAMERA

		// GET_CAMERA_PROPS
		"camera__get_props_fun": ()=>{
			const camera_props_map = {
				"x": camera.position.x,
				"y": camera.position.y,
				"z": camera.position.z,
			};
			return camera_props_map;
		},

		// CAMERA__ANIMATE
		"camera__animate_fun": (p_animations_lst)=>{

			const animations_packaged_lst = [];
			for (const animation_map of p_animations_lst) {
				const duration_sec_f  = animation_map["duration_sec_f"];
				const start_props_map = animation_map["start_props_map"];
				const end_props_map   = animation_map["end_props_map"];
				const repeat_bool     = animation_map["repeat_bool"];

				const animation_packaged_map = {
					"props_to_animate_lst": [
						["x", start_props_map["x"], end_props_map["x"]],
						["y", start_props_map["y"], end_props_map["y"]],
						["z", start_props_map["z"], end_props_map["z"]],
						["color_background_r", start_props_map["color_background_r"], end_props_map["color_background_r"]],
						["color_background_g", start_props_map["color_background_g"], end_props_map["color_background_g"]],
						["color_background_b", start_props_map["color_background_b"], end_props_map["color_background_b"]]
					],
					"duration_sec_f": duration_sec_f,
					"repeat_bool":    repeat_bool,
				};
				animations_packaged_lst.push(animation_packaged_map);
			}

			gf_animations.apply(animations_packaged_lst, camera, scene_3d);
		},

		//-------------------------------------------------
		// CUBE
		"create_cube_fun": (p_x :number, p_y :number, p_z :number,
			p_rx :number, p_ry :number, p_rz :number,
			p_sx :number, p_sy :number, p_sz :number,
			p_cr :number, p_cg :number, p_cb :number)=>{
			

			if (shader_material_bool) {
				// When cloning a ShaderMaterial, the attributes and vertex/fragment programs
				// are copied by reference. Only the uniforms are copied by value.
				material = material_shader_active.clone();
			}
			else {

				// MATERIAL
				material = new THREE.MeshBasicMaterial({
					color:               new THREE.Color(p_cr, p_cg, p_cb), // color_int,
					polygonOffset:       true,
					polygonOffsetFactor: 1, // positive value pushes polygon further away
					polygonOffsetUnits:  1
				});
			}
			

			// get global coords that are derived from the latest coordinate system origin
			// that was last added to the stack of coordinate system origins.
			const [derived_x_f, derived_y_f, derived_z_f] = 
				gf_engine_utils.get_derived_coords(p_x, p_y, p_z, coord_origins_stack_lst);

			// const [derived_rx_f, derived_ry_f, derived_rz_f] = 
			// 	gf_engine_utils.get_derived_rotation(p_rx, p_ry, p_rz, rotations_stack_lst);

			
			const mesh = create_cube(derived_x_f, derived_y_f, derived_z_f,
				p_rx, p_ry, p_rz,
				p_sx, p_sy, p_sz,
				p_cr, p_cg, p_cb,
				material,
				wireframe_material_bool,
				rotation_pivot_points_stack_lst);

			// ANIMATION - there is an active animation, so apply it to object.
			//             has to be done in the obj creation 
			if (animation_active_map != null) {
				gf_animations.apply([animation_active_map], mesh, scene_3d);
			}

			return mesh;
		},

		//-------------------------------------------------
		// SPHERE
		"create_sphere_fun": (p_x :number, p_y :number, p_z :number,
			p_rx :number, p_ry :number, p_rz :number,
			p_sx :number, p_sy :number, p_sz :number,
			p_cr :number, p_cg :number, p_cb :number)=>{
			
			if (shader_material_bool) {
				// When cloning a ShaderMaterial, the attributes and vertex/fragment programs
				// are copied by reference. Only the uniforms are copied by value.
				material = material_shader_active.clone();

			}
			else {

				// MATERIAL
				material = new THREE.MeshBasicMaterial({
					color:               new THREE.Color(p_cr, p_cg, p_cb), // color_int,
					polygonOffset:       true,
					polygonOffsetFactor: 1, // positive value pushes polygon further away
					polygonOffsetUnits:  1
				});
			}
			

			// get global coords that are derived from the latest coordinate system origin
			// that was last added to the stack of coordinate system origins.
			const [derived_x_f, derived_y_f, derived_z_f] = 
				gf_engine_utils.get_derived_coords(p_x, p_y, p_z, coord_origins_stack_lst);

			
			const mesh = create_sphere(derived_x_f, derived_y_f, derived_z_f,
				p_rx, p_ry, p_rz,
				p_sx, p_sy, p_sz,
				p_cr, p_cg, p_cb,
				material,
				wireframe_material_bool);
		},

		//-------------------------------------------------
		// LINE
		"create_line_fun": (p_x :number, p_y :number, p_z :number,
			p_rx :number, p_ry :number, p_rz :number,
			p_sx :number, p_sy :number, p_sz :number,
			p_cr :number, p_cg :number, p_cb :number)=>{

			
			// get global coords that are derived from the latest coordinate system origin
			// that was last added to the stack of coordinate system origins.
			const [derived_x_f, derived_y_f, derived_z_f] = 
				gf_engine_utils.get_derived_coords(p_x, p_y, p_z, coord_origins_stack_lst);

			
			create_line(derived_x_f, derived_y_f, derived_z_f,
				p_rx, p_ry, p_rz,
				p_sx, p_sy, p_sz,
				p_cr, p_cg, p_cb,
				rotation_pivot_points_stack_lst);
		},

		//-------------------------------------------------
		// ANIMATE
		"animate_fun": (p_props_to_animate_lst, p_duration_sec_f, p_repeat_bool)=>{

			const animation_map = {
				"props_to_animate_lst": [],
				"duration_sec_f":       p_duration_sec_f,
				"repeat_bool":          p_repeat_bool,
			};
			for (const prop_map of p_props_to_animate_lst) {

				const name_str    = prop_map["name_str"];
				const start_val_f = prop_map["start_val_f"];
				const end_val_f   = prop_map["end_val_f"];
				animation_map["props_to_animate_lst"].push([name_str, start_val_f, end_val_f]);
			}

			animation_active_map = animation_map;
		}

		//-------------------------------------------------
	}

	return api_map;
}