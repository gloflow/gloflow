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




from PIL import Image, ImageDraw

import numpy as np

#----------------------------------------------
def main():

    test_image_paths_lst = [
        './../gf_ml_worker/test/data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg',
        './../gf_ml_worker/test/data/input/canvas.png',
        "./../gf_ml_worker/test/data/input/4b14ca75070ac78323cf2ddef077ae92_thumb_medium.jpeg",
        "./../gf_ml_worker/test/data/input/4838df39722bc2d681b67bf739f29357_thumb_small.jpeg",
        "./../gf_ml_worker/test/data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg",
        "./../gf_ml_worker/test/data/input/3a61e7d68fb17198e8dc0476cc862ddd_thumb_small.jpeg",
    ]


    for i in range(0, len(test_image_paths_lst)):
        img = Image.open(test_image_paths_lst[i])
        # pix = img.load()
        width_int, height_int = img.size

        o_str = f"out_{i}.png"
        run(img, o_str, width_int, height_int)



#----------------------------------------------
def run(p_img, p_pallete__output_file_path_str, p_img_width_int, p_img_height_int):

    img_pixels_arr = np.array(p_img.getdata())
    
    #----------------------------------------------
    def process_img__with_px_coords(p_img_pixels_with_global_index_arr, p_level_int, p_levels_max_int=3):

        # TERMINATION
        if p_level_int >= p_levels_max_int:
            

            img_pixels_global_indexes_arr = p_img_pixels_with_global_index_arr[:, 0]

  
            # for all rows, only include last 3 elements (rgb pixel values), since thats what
            # we want to calculate the average with.
            img_pixels_no_global_indexes_arr = p_img_pixels_with_global_index_arr[:, 1:4]

            # get average pixel, by taking averages of all pixels on each of the channels.
            # this is the last image subdivision (dividing by pixel channel with maximum range), 
            # and so calculate average pixel value for that subdivision and return.
            average_pixel_arr = np.average(img_pixels_no_global_indexes_arr.T, axis=1)

            # round and cast as integer
            a = np.round(average_pixel_arr).astype(int)
            return [(img_pixels_global_indexes_arr, a[0], a[1], a[2])]

        # p_img_pixels_with_global_index_arr.T output:
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


        ranges_per_channel_lst           = img_by_channel_arr.ptp(axis=1)
        channel_with_max_range_index_int = np.argmax(ranges_per_channel_lst)


        # img_by_channel_arr+1 - to account for the fact that the first element is global_index,
        #                        so the color index has to be offset by 1.
        img_max_range_channel_vals_arr = p_img_pixels_with_global_index_arr[:, channel_with_max_range_index_int+1]




        img_max_range_channel_vals_sorted_indexes_arr            = img_max_range_channel_vals_arr.argsort()
        upper_half__px_indicies_arr, lower_half__px_indicies_arr = np.array_split(img_max_range_channel_vals_sorted_indexes_arr, 2)




        upper_half__px_arr = p_img_pixels_with_global_index_arr[upper_half__px_indicies_arr]
        lower_half__px_arr = p_img_pixels_with_global_index_arr[lower_half__px_indicies_arr]


        upper_half__avrg_pixels_arr = process_img__with_px_coords(upper_half__px_arr, p_level_int+1)
        lower_half__avrg_pixels_arr = process_img__with_px_coords(lower_half__px_arr, p_level_int+1)

        r=[]
        r.extend(upper_half__avrg_pixels_arr)
        r.extend(lower_half__avrg_pixels_arr)
        return r

    #----------------------------------------------
    # DEPRECATED!! - use the above function that outputs global_indexes for pixels
    #                belonging to average color values.
    def process_img(p_img_pixels_arr, p_level_int, p_levels_max_int=3):
        
        # p_img_pixels_arr
        # shape - (px_num, 3) - 3 element tuple per pixel
        # [[r1, g1, b1], [r2, g2, b2], ...]


        # TERMINATION
        if p_level_int >= p_levels_max_int:
            
            # get average pixel, by taking averages of all pixels on each of the channels.
            # this is the last image subdivision (dividing by pixel channel with maximum range), 
            # and so calculate average pixel value for that subdivision and return.
            # 
            # p_img_pixels_arr.T - average is calculated on a transpose of the pixels matrix,
            #                      so that each channel is an array, and average is calculated on each channel. 
            average_pixel_arr = np.average(p_img_pixels_arr.T, axis=1)
            a = np.round(average_pixel_arr).astype(int)
            return [(a[0], a[1], a[2])]


    
        # transpose - groups pixel values per channel
        # shape - (3, px_num)
        # [
        #   [r1, r2, ...],
        #   [g1, g2, ...],
        #   [b1, b2, ...]   
        # ]
        img_by_channel_arr = p_img_pixels_arr.T

        #----------
        # PICK_CHANNEL_WITH_MAX_RANGE
        # get the range (max_val-min_val) of each of the r/g/b channels.
        ranges_per_channel_lst           = img_by_channel_arr.ptp(axis=1)
        channel_with_max_range_index_int = np.argmax(ranges_per_channel_lst)
        
        # values of pixels in a particular channel which has the maximum range of values.
        img_max_range_channel_vals_arr = p_img_pixels_arr[:, channel_with_max_range_index_int]

        #----------
        # SORT_BY_CHANNEL
        # argsort() - sorts values and returns their indexes
        img_max_range_channel_vals_sorted_indexes_arr = img_max_range_channel_vals_arr.argsort()

        # split sorted list of pixels (on channel with max range by pixel value) indicies
        # into half/median. 
        # np.array_split() - wont error if the image doesnt split in 2 equal parts.
        upper_half__px_indicies_arr, lower_half__px_indicies_arr = np.array_split(img_max_range_channel_vals_sorted_indexes_arr, 2)
        
        # index into image pixels by sorted indexes of the channel with biggest range/variance
        upper_half__px_arr = p_img_pixels_arr[upper_half__px_indicies_arr]
        lower_half__px_arr = p_img_pixels_arr[lower_half__px_indicies_arr]



        upper_half__avrg_pixels_arr = process_img(upper_half__px_arr, p_level_int+1)
        lower_half__avrg_pixels_arr = process_img(lower_half__px_arr, p_level_int+1)

        r=[]
        r.extend(upper_half__avrg_pixels_arr)
        r.extend(lower_half__avrg_pixels_arr)
        return r

    #----------------------------------------------
    


    # ADD_PIXEL_GLOBAL_INDEX - index of the pixel relative to the whole image.
    #                          this is used to be able to track where in the image particular pixels come from.
    # [[r0, g0, b0], [r1, g1, b1], [r2, g2, b2], ...] -> [[0, r0, g0, b0], [1, r1, g1, b1], [2, r2, g2, b2], ...]
    # 
    # np.column_stack() - Stack 1-D arrays as columns into a 2-D array.
    img_pixels_with_global_index_arr = np.column_stack((np.arange(len(img_pixels_arr)), img_pixels_arr))



    r = process_img__with_px_coords(img_pixels_with_global_index_arr, 0)

    # r = process_img(img_pixels_arr, 0)



    # IMAGE_PALETTE
    r_img = Image.new('RGB', (len(r), 1))

    palette_pixels_only_arr = np.array(r)[:, 1:4]

    print(palette_pixels_only_arr)
    print(palette_pixels_only_arr.shape)

    # r_img.putdata(palette_pixels_only_arr)
    r_img = Image.fromarray(palette_pixels_only_arr.reshape(palette_pixels_only_arr.shape[0], 1, 3).astype(np.uint8))
    r_img.save(p_pallete__output_file_path_str)



    #----------
    




    # colorize original image with quantized color palette,
    # using the global_indexes of every average pixel color calcuated
    # to set those indexes to those exact values.
    # ADD!! - figure out how to vectorize this operation, withyout needing
    #         to iterate through pixels individually.
    for c in r:
        print(c)

        global_indexes_arr = c[0]
        color_arr          = c[1:4]

        for i in global_indexes_arr:
            img_pixels_arr[i] = color_arr

    # height/width/3 - height has to go before width
    img_pixels_3d_arr = img_pixels_arr.reshape(p_img_height_int, p_img_width_int, 3)



    r_img = Image.fromarray(img_pixels_3d_arr.astype(np.uint8))

    #----------------------------------------------
    # def draw_sectors():
    #     draw = ImageDraw.Draw(r_img)
    #     draw.rectangle(((0, 00), (100, 100)), outline='black', width=1)

    #----------------------------------------------
    # draw_sectors()


    r_img.save(f"{p_pallete__output_file_path_str}__sectors.png")

    #----------

#----------------------------------------------
main()