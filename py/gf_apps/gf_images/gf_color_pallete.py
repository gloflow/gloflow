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




from PIL import Image

import numpy as np

def main():

    # './../gf_ml_worker/test/data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg'
    # './../gf_ml_worker/test/data/input/canvas.png'
    # "./../gf_ml_worker/test/data/input/4b14ca75070ac78323cf2ddef077ae92_thumb_medium.jpeg"
    # "./../gf_ml_worker/test/data/input/4838df39722bc2d681b67bf739f29357_thumb_small.jpeg
    # "./../gf_ml_worker/test/data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg"
    img = Image.open("./../gf_ml_worker/test/data/input/3a61e7d68fb17198e8dc0476cc862ddd_thumb_small.jpeg")
    # pix = img.load()





    # width_int, height_int = img.size




    img_pixels_arr = np.array(img.getdata())


    def process_img(p_img_pixels_arr, p_i_int, p_subdivisons_int=5):
        
        # p_img_pixels_arr
        # shape - (px_num, 3) - 3 element tuple per pixel
        # [[r1, g1, b1], [r2, g2, b2], ...]

        if p_i_int >= p_subdivisons_int:
            
            # get average pixel, by taking averages of all pixels on each of the channels.
            # this is the last image subdivision (dividing by pixel channel with maximum range), 
            # and so calculate average pixel value for that subdivision and return.
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



        # get the range (max_val-min_val) of each of the r/g/b channels.
        ranges_per_channel_lst     = img_by_channel_arr.ptp(axis=1)
        channel_with_max_range_int = np.argmax(ranges_per_channel_lst)
        # print(f"ranges per channel (RGB) - {ranges_per_channel_lst}")
        # print(f"channel with max range   - {channel_with_max_range_int}")
        

        # values of pixels in a particular channel which has the maximum range of values.
        img_channel_vals_arr = p_img_pixels_arr[:, channel_with_max_range_int]

        # argsort() - sorts values and returns their indexes
        img_channel_vals_sorted_indexes_arr = img_channel_vals_arr.argsort()

        # split sorted list of pixels (on channel with max range by pixel value) indicies
        # into half/median. 
        # np.array_split() - wont error if the image doesnt split in 2 equal parts.
        upper_half__px_indicies_arr, lower_half__px_indicies_arr = np.array_split(img_channel_vals_sorted_indexes_arr, 2)
        
        # index into image pixels by sorted indexes of the channel with biggest range/variance
        upper_half__px_arr = p_img_pixels_arr[upper_half__px_indicies_arr]
        lower_half__px_arr = p_img_pixels_arr[lower_half__px_indicies_arr]



        upper_half__avrg_pixels_arr = process_img(upper_half__px_arr, p_i_int+1)
        lower_half__avrg_pixels_arr = process_img(lower_half__px_arr, p_i_int+1)

        r=[]
        r.extend(upper_half__avrg_pixels_arr)
        r.extend(lower_half__avrg_pixels_arr)
        return r

    r = process_img(img_pixels_arr, 0)


    print("result")
    print(r)


    
    r_img = Image.new('RGB', (len(r), 1))

    r_img.putdata(r)
        

    r_img.save('test.png')

main()