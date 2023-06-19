# GloFlow application and media management/publishing platform
# Copyright (C) 2021 Ivan Trajkovic
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

import os
import argparse
import json
from colored import fg, bg, attr
import numpy as np
from PIL import Image, ImageDraw

#----------------------------------------------
def run_from_file(p_image_file_path_str):
	assert os.path.isfile(p_image_file_path_str)

	image = Image.open(p_image_file_path_str)
	width_int, height_int = image.size

	palette_pixels_only_arr, result = run(image, width_int, height_int)

	return palette_pixels_only_arr, result

#----------------------------------------------
def run(p_image,
	p_img_width_int,
	p_img_height_int):

	img_pixels_original_arr = np.array(p_image.getdata())

	# only select the first 3 channels of each pixel, since some images have
	# an alpha channel as well which we dont want.
	img_pixels_arr = img_pixels_original_arr[:, :3]

	#----------------------------------------------
	# levels_max gives a color palette of certain size:
	# 3 - 8  colors
	# 4 - 16 colors
	# 5 - 32 colors
	# 6 - 64 colors

	def process_img__with_px_coords(p_img_pixels_with_global_index_arr,
		p_level_int,
		p_levels_max_int=3):
		
		# TERMINATION
		# check if we reached the maximum level of subdivision/recursion 
		if p_level_int >= p_levels_max_int:
			

			img_pixels_global_indexes_arr = p_img_pixels_with_global_index_arr[:, 0]

  
			# for all rows, only include last 3 elements (rgb pixel values), since thats what
			# we want to calculate the average with.
			img_pixels_no_global_indexes_arr = p_img_pixels_with_global_index_arr[:, 1:4]

			# get average pixel, by taking averages of all pixels on each of the channels.
			# this is the last image subdivision (dividing by pixel channel with maximum range), 
			# and so calculate average pixel value for that subdivision and return.
			# 
			# img_pixels_no_global_indexes_arr.T - average is calculated on a transpose of the pixels matrix,
			#                                      so that each channel is an array, and average is
			#                                      calculated on each channel. 
			average_pixel_arr = np.average(img_pixels_no_global_indexes_arr.T, axis=1)

			# round and cast as integer
			a = np.round(average_pixel_arr).astype(int)

			return [(img_pixels_global_indexes_arr,
				a[0],  # R 
				a[1],  # G
				a[2])] # B

		# remove the pixels global index (index of the pixel relative to the whole image)
		# 
		# input:
		# p_img_pixels_with_global_index_arr
		# [[     0      1      2 ... 115997 115998 115999]
		# [   162    163    164 ...    142    140    140]
		# [   200    201    202 ...    193    191    191]
		# [   203    204    205 ...    194    195    195]]
		# 
		# [1:] - exclude the the first element/row of the transposed matrix,
		#        since its the pixel_global_index and we dont need ranges on that.
		# 
		# output:
		# [   162    163    164 ...    142    140    140]
		# [   200    201    202 ...    193    191    191]
		# [   203    204    205 ...    194    195    195]]
		img_by_channel_arr = p_img_pixels_with_global_index_arr.T[1:]

		#----------
		# PICK_CHANNEL_WITH_MAX_RANGE
		# get the range (max_val-min_val) of each of the r/g/b channels.
		# .ptp() - stands for "peak-to-peak" and it calculates the range of values in a NumPy array along a specified axis.
		#          returns the difference between the maximum and minimum values in the array along the specified axis.
		ranges_per_channel_lst = img_by_channel_arr.ptp(axis=1)

		# argmax() - returns the indices of the maximum values along a specified axis of a NumPy array.
		channel_with_max_range_index_int = np.argmax(ranges_per_channel_lst)

		# values of pixels in a particular channel which has the maximum range of values.
		#
		# img_by_channel_arr+1 - to account for the fact that the first element is global_index,
		#                        so the color index has to be offset by 1.
		img_max_range_channel_vals_arr = p_img_pixels_with_global_index_arr[:, channel_with_max_range_index_int+1]

		#----------
		# SORT_BY_CHANNEL
		# argsort() - sorts values and returns their indexes
		img_max_range_channel_vals_sorted_indexes_arr = img_max_range_channel_vals_arr.argsort()

		# split sorted list of pixels (on channel with max range by pixel value) indicies into half/median. 
		# np.array_split() - wont error if the image doesnt split in 2 equal parts.
		upper_half__px_indicies_arr, lower_half__px_indicies_arr = np.array_split(img_max_range_channel_vals_sorted_indexes_arr, 2)

		# index into image pixels by sorted indexes of the channel with biggest range/variance
		upper_half__px_arr = p_img_pixels_with_global_index_arr[upper_half__px_indicies_arr]
		lower_half__px_arr = p_img_pixels_with_global_index_arr[lower_half__px_indicies_arr]

		# RECURSION
		upper_half__avrg_pixels_arr = process_img__with_px_coords(upper_half__px_arr, p_level_int+1)
		lower_half__avrg_pixels_arr = process_img__with_px_coords(lower_half__px_arr, p_level_int+1)

		r=[]
		r.extend(upper_half__avrg_pixels_arr)
		r.extend(lower_half__avrg_pixels_arr)
		return r

	#----------------------------------------------

	# ADD_PIXEL_GLOBAL_INDEX - add an index of the pixel relative to the whole image, to the pixels rgb color value.
	#                          this is used to be able to track where in the image particular pixels come from.
	# [[r0, g0, b0], [r1, g1, b1], [r2, g2, b2], ...] -> [[0, r0, g0, b0], [1, r1, g1, b1], [2, r2, g2, b2], ...]
	# 
	# np.column_stack() - Stack 1-D arrays as columns into a 2-D array.
	img_pixels_with_global_index_arr = np.column_stack((np.arange(len(img_pixels_arr)), img_pixels_arr))


	start_level_int = 0
	result = process_img__with_px_coords(img_pixels_with_global_index_arr, start_level_int)


	# IMAGE_PALETTE
	r_img = Image.new('RGB', (len(result), 1))



	# get the palette pixels only, removing global pixel indexes from each of the palette pixel part.
	# each of the palette pixels has a format (px_global_indexes_arr, r, g, b), so only get (r, g, b).
	palette_pixels_only_arr = np.empty((len(result), 3), dtype=int)
	for i in range(len(result)):
		pixels_global_indexes_lst, r_int, g_int, b_int = result[i]
		palette_pixels_only_arr[i] = np.array([r_int, g_int, b_int])


		print(pixels_global_indexes_lst)


		print(f"global {len(result[i][0])} {palette_pixels_only_arr[i]}")

	print(palette_pixels_only_arr)
	print(palette_pixels_only_arr.shape)

	return palette_pixels_only_arr, result

#----------------------------------------------
def run_extended(p_image_file_path_str,
	p_palette__output_file_path_str):
 
	palette_pixels_only_arr, result = run_from_file(p_image_file_path_str)

	r_img = Image.fromarray(palette_pixels_only_arr.reshape(palette_pixels_only_arr.shape[0], 1, 3).astype(np.uint8))
	r_img.save(p_palette__output_file_path_str)


	#----------
	
	image = Image.open(p_image_file_path_str)
	width_int, height_int = image.size

	img_pixels_original_arr = np.array(image.getdata())

	# only select the first 3 channels of each pixel, since some images have
	# an alpha channel as well which we dont want.
	img_pixels_arr = img_pixels_original_arr[:, :3]

	# colorize original image with quantized color palette,
	# using the global_indexes of every average pixel color calcuated
	# to set those indexes to those exact values.
	# ADD!! - figure out how to vectorize this operation, withyout needing
	#         to iterate through pixels individually.
	colors_lst = []
	for c in result:

		global_indexes_arr = c[0]
		color_arr          = c[1:4]
		colors_lst.append(color_arr)

		# color_4d_arr = np.array(color_arr)
		# color_4d_arr = np.append(color_4d_arr, 1).reshape((4, ))

		for i in global_indexes_arr:
			img_pixels_arr[i] = color_arr

	# height/width/3 - height has to go before width
	img_pixels_3d_arr = img_pixels_arr.reshape(height_int, width_int, 3)

	r_img = Image.fromarray(img_pixels_3d_arr.astype(np.uint8))

	#----------------------------------------------
	# def draw_sectors():
	#     draw = ImageDraw.Draw(r_img)
	#     draw.rectangle(((0, 00), (100, 100)), outline='black', width=1)

	#----------------------------------------------
	# draw_sectors()

	r_img.save(f"{p_palette__output_file_path_str}__sectors.png")

	#----------

	print("K MEANS ------------")
	kmeans(colors_lst)

#----------------------------------------------
def kmeans(p_colors_lst, p_k = 5):

	import scipy.spatial.distance
	import matplotlib.pyplot as plt
	import sklearn.cluster

	# array([[ 25, 236, 102],
	#       [129, 147, 103],
	#       [226,  35,  64],
	#       [ 94,  15,  33],
	#       [ 56, 104, 115]], dtype=uint64)
	centroids_arr = (np.random.rand(p_k, 3)*255).round().astype(np.uint)

	print("------------------------")
	print(p_colors_lst)
	print("==")
	print(centroids_arr)

	#----------------------------------------------
	def algo():
		for i in range(0, 10):

			# ASSIGN_COLORS_TO_CENTROIDS
			colors_centroids_lst = []
			for color_tpl in p_colors_lst:
				color_to_centroids_distances_arr = scipy.spatial.distance.cdist(centroids_arr, [color_tpl], 'euclidean')

				# print("result:")
				# print(color_to_centroids_distances_arr)

				# get the index/label of the centroid that the color is closest to
				closest_centroid_index_int = color_to_centroids_distances_arr.argmin()
				colors_centroids_lst.append(closest_centroid_index_int)
			
			# GROUP_COLORS_BY_CENTROIDS
			centroids_colors_map = {}
			for i in range(0, len(colors_centroids_lst)):

				color_arr          = p_colors_lst[i]
				centroid_label_int = colors_centroids_lst[i]
				if centroid_label_int in centroids_colors_map.keys():
					centroids_colors_map[centroid_label_int].append(color_arr)
				else:
					centroids_colors_map[centroid_label_int] = [color_arr]



			for centroid_label_int, colors_lst in centroids_colors_map.items():

				colors_by_channel_arr = np.array(colors_lst).T
				new_r_f = np.average(colors_by_channel_arr[0])
				new_g_f = np.average(colors_by_channel_arr[1])
				new_b_f = np.average(colors_by_channel_arr[2])

				centroid_new_coords_arr = np.array([new_r_f, new_g_f, new_b_f])

				# print(centroid_new_coords_arr)

				centroids_arr[centroid_label_int] = centroid_new_coords_arr
		
		return colors_centroids_lst

	#----------------------------------------------
	def plot_colors(p_colors_centroids_lst):
		

		fig = plt.figure()
		ax1 = fig.add_subplot(1, 2, 1, projection="3d") # plt.axes(projection='3d')
		ax2 = fig.add_subplot(1, 2, 2, projection="3d")

		color_map_lst = [
			"red",
			"green",
			"blue",
			"yellow",
			"orange"
		]

		#----------------
		# CLUSTERS
		colors_rgb_coords_arr = np.array(p_colors_lst).T
		colors_x_arr = colors_rgb_coords_arr[0]
		colors_y_arr = colors_rgb_coords_arr[1]
		colors_z_arr = colors_rgb_coords_arr[2]

		colors_1range_arr = np.array(p_colors_lst)/255
		# ax.scatter3D(colors_x_arr, colors_y_arr, colors_z_arr, c=colors_1range_arr)

		colors_lst = [color_map_lst[l] for l in colors_centroids_lst]
		ax1.scatter3D(colors_x_arr, colors_y_arr, colors_z_arr, c=colors_lst)

		#----------------

		# CENTROIDS
		centroid_coords_arr = centroids_arr.T
		centroids_x_arr = centroid_coords_arr[0]
		centroids_y_arr = centroid_coords_arr[1]
		centroids_z_arr = centroid_coords_arr[2]
		ax1.scatter3D(centroids_x_arr, centroids_y_arr, centroids_z_arr, c="black")

		# plt.show()

		#----------------
		# KMEANS_SKLEARN		
		clusters_labels_lst = sklearn.cluster.KMeans(n_clusters=p_k).fit_predict(p_colors_lst)
		print(clusters_labels_lst)

		
		colors_lst = [color_map_lst[l] for l in clusters_labels_lst]
		ax2.scatter3D(colors_x_arr, colors_y_arr, colors_z_arr, c=colors_lst) # colors_1range_arr)

		#----------------

		plt.show()

	#----------------------------------------------
	colors_centroids_lst = algo()
	plot_colors(colors_centroids_lst)

#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)
	
	#----------------------------
	# INPUT_IMAGES_LOCAL_FILE_PATHS
	arg_parser.add_argument("-input_images_local_file_paths", action = "store", default=None,
		help = "list of image file paths (',' delimited) to process")

	#----------------------------
	# OUTPUT_DIR_PATH
	arg_parser.add_argument("-output_dir_path", action = "store", default=None,
		help = "dir path of the output images")

	#----------------------------
	# MEDIAN_CUT_LEVELS_NUM
	arg_parser.add_argument("-median_cut_levels_num", action = "store", default=None,
		help = "number of levels to use for the median-cut algo subdivisions")

	#----------------------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)

	return {
		"input_images_local_file_paths_str": args_namespace.input_images_local_file_paths,
		"output_dir_path_str":               args_namespace.output_dir_path,
		"median_cut_levels_num":             int(args_namespace.median_cut_levels_num)
	}

#----------------------------------------------