# GloFlow application and media management/publishing platform
# Copyright (C) 2020 Ivan Trajkovic
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

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import numpy as np
import skimage.color
import skimage.io
import requests

#----------------------------------------------
def download_test_img():
    target_test_img_str = "%s/../data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg"%(modd_str)
    url_str             = "http://gf--img.s3-website-us-east-1.amazonaws.com/thumbnails/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg"
    
    if not os.path.isfile(target_test_img_str):
        
        r = requests.get(url_str)
        assert r.status_code == 200
        f = open(target_test_img_str, "wb")
        f.write(r.content)
        f.close()
    
    # LOAD_FROM_FILE
    input_img = skimage.io.imread(target_test_img_str)
    assert isinstance(input_img, np.ndarray)
    print("input img - %s"%(str(input_img.shape)))


    img_grayscale = skimage.color.rgb2gray(input_img)
    img_croped    = img_crop(img_grayscale, 8, 8)

    
    # add another dimension to data, first one, to have a batch-like np array
    input_img_batch = np.expand_dims(img_croped, 0)

    return input_img_batch, target_test_img_str

#----------------------------------------------
def load_data():

    input_img_path_str = "%s/../../../../../rust/gf_images_jobs/test/data/input/test_img_1.jpeg"%(modd_str)
    assert os.path.isfile(input_img_path_str)

    
    # LOAD_FROM_FILE
    input_img = skimage.io.imread(input_img_path_str)
    assert isinstance(input_img, np.ndarray)
    print("input img - %s"%(str(input_img.shape)))

   
    # CROP
    img_croped            = img_crop(input_img, 64, 64)
    input_img_transformed = img_croped.astype("float32")

    # plt.imshow(input_img_transformed.astype("uint8"))
    # plt.show()

    input_img_batch = np.expand_dims(input_img_transformed, 0)
    print(input_img_batch.shape)

    return input_img_batch
    
#----------------------------------------------
def img_crop(p_img_np, p_width_int, p_height_int):
    h = p_img_np.shape[0]
    w = p_img_np.shape[1]

    crop_origin_x = w//2 - p_width_int//2
    crop_origin_y = h//2 - p_height_int//2
    img_crop = p_img_np[crop_origin_x:crop_origin_x+p_width_int, 
        crop_origin_y:crop_origin_y+p_height_int]
    return img_crop