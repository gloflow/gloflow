# GloFlow application and media management/publishing platform
# Copyright (C) 2023 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

import json
from PIL import Image

import gf_color_palette

#----------------------------------------------
def main(p_test_bool=True):

	# TEST
	if p_test_bool:
		test_image_paths_lst = [
			'./../../gf_ml_worker/test/data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg',
			'./../../gf_ml_worker/test/data/input/canvas.png',
			"./../../gf_ml_worker/test/data/input/4b14ca75070ac78323cf2ddef077ae92_thumb_medium.jpeg",
			"./../../gf_ml_worker/test/data/input/4838df39722bc2d681b67bf739f29357_thumb_small.jpeg",
			"./../../gf_ml_worker/test/data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg",
			"./../../gf_ml_worker/test/data/input/3a61e7d68fb17198e8dc0476cc862ddd_thumb_small.jpeg",
		]


		for i in range(0, len(test_image_paths_lst)):

			image_file_path_str = test_image_paths_lst[i]
			o_str               = f"./output/out_{i}.png"

			# RUN
			gf_color_palette.run_extended(image_file_path_str, o_str)

	# PRODUCTION
	else:

		args_map = parse_args()

		input_images_local_file_paths_lst = args_map["input_images_local_file_paths_str"].split(",")
		
		# RUN
		run_multiple(input_images_local_file_paths_lst)

		out_map = {}
		print(f"GF_OUT:{json.dumps(out_map)}")

#----------------------------------------------
def run_multiple(p_input_images_local_file_paths_lst):

	print(f"RUNNING {fg('green')}GF_COLOR_PALETTE{attr(0)} PLUGIN ")
	
	# VERIFY
	for f in p_input_images_local_file_paths_lst:
		assert os.path.isfile(f)
		print(f"input - {fg('yellow')}{f}{attr(0)}")

	# RUN
	for f in p_input_images_local_file_paths_lst:
		gf_color_palette.run_extended(f)

#----------------------------------------------
if __name__ == "__main__":
	main()