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
declare var TWEEN;

//-------------------------------------------------
export function apply(p_animations_packaged_lst, p_obj_3d, p_scene_3d) {


	var tweens_lst = [];
	for (const animation_packaged_map of p_animations_packaged_lst) {
		const props_to_animate_lst = animation_packaged_map["props_to_animate_lst"];
		const duration_sec_f       = animation_packaged_map["duration_sec_f"];
		const repeat_bool          = animation_packaged_map["repeat_bool"];

		// animation state thats changed over time, initialized to its start values
		const anim_state_map   = {};
		const anim_targets_map = {};

		for (const [name_str, start_val_f, end_val_f] of props_to_animate_lst) {
			anim_state_map[name_str]   = start_val_f; // initial value of animation_state
			anim_targets_map[name_str] = end_val_f;   // final value that animation_state will arrive at
		}

		const tween = create__simple(p_obj_3d,
			p_scene_3d,
			anim_targets_map,
			anim_state_map,
			duration_sec_f,
			repeat_bool);

		tweens_lst.push(tween);
	}

	// if there are multiple animations then chain them together
	if (p_animations_packaged_lst.length > 1) {
		for (var i=1; i<tweens_lst.length; i++) {
			tweens_lst[i-1].chain(tweens_lst[i]);
		}
	}

	// start the first animation in the sequence
	tweens_lst[0].start();
}

//-------------------------------------------------
export function create__simple(p_obj_3d,
	p_scene_3d,
	p_anim_targets_map,
	p_anim_state_map,
	p_animation_durration_sec_f,
	p_repeat_bool) {

	const tween = new TWEEN.Tween(p_anim_state_map)
		.to(p_anim_targets_map, p_animation_durration_sec_f*1000)
		.easing(TWEEN.Easing.Quadratic.Out)
		.onUpdate(() => {
			
			if (p_anim_state_map.x != undefined) {
				p_obj_3d.position.x = p_anim_state_map.x;
			}
			if (p_anim_state_map.y != undefined) {
				p_obj_3d.position.y = p_anim_state_map.y;
			}
			if (p_anim_state_map.z != undefined) {
				p_obj_3d.position.z = p_anim_state_map.z;
			}
			if (p_anim_state_map.color_background_r != undefined) {
				const color_background = new THREE.Color(
					p_anim_state_map.color_background_r,
					p_anim_state_map.color_background_g,
					p_anim_state_map.color_background_b
				);

				p_scene_3d.background = color_background;
			}
		})
	
	if (p_repeat_bool) tween.repeat(Infinity);

	return tween;
}