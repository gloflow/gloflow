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

import random
import numpy as np

import delegator
print(delegator.run("ls -al %s"%("%s/../../build"%(modd_str))).out)

sys.path.append("%s/../../build"%(modd_str))
import gf_images_jobs_py as gf_images_jobs


print("TEST ------------------------------------------------------------------------")
print("\n\nRUN THIS TEST WITH - 'LD_LIBRARY_PATH=. python3 test.py'\n\n")
assert "LD_LIBRARY_PATH" in os.environ.keys()

print(dir(gf_images_jobs))

#---------------------------------------------------------------------------
def test__numpy():

    x       = np.arange(10, dtype=np.float64) # 1_000_000_000, dtype=np.float64)
    x_in_gb = x.nbytes / 1024 / 1024 / 1024
    print(x_in_gb)

    #---------------------------
    # 4D
    x_4d_lst = []

    for a in range(0, 100):

        x_3d_lst = []
        for i in range(0, 100):

            row_lst = []
            for j in range(0, 100):
                
                px_rgb = np.array([
                    random.random()*255.0,
                    random.random()*255.0,
                    random.random()*255.0])

                row_lst.append(px_rgb)
            x_3d_lst.append(row_lst)

        x_4d_lst.append(x_3d_lst)

    x_4d = np.array(x_4d_lst, dtype=np.float64)

    # print("sleep 10s...")
    # import time
    # time.sleep(10)

    print("draw 4d...")

    output_file_path_str = "%s/data/output/test__numpy_4d.jpeg"%(modd_str)
    assert os.path.isdir(os.path.dirname(output_file_path_str))

    gf_images_jobs.view_numpy_arr_4D(x_4d,
        output_file_path_str,
        1000, # collage width
        1000, # collage height
        10,   # rows
        10)   # columns

    #---------------------------
    # 3D
    x_3d_lst = []
    for i in range(0, 1_000):

        row_lst = []
        for j in range(0, 1_000):
            
            px_rgb = np.array([
                random.random()*255.0,
                random.random()*255.0,
                random.random()*255.0])

            row_lst.append(px_rgb)
        x_3d_lst.append(row_lst)

    x_3d = np.array(x_3d_lst, dtype=np.float64)

    # print("sleep 10s...")
    # import time
    # time.sleep(10)

    print("draw 3d...")
    output_file_path_str = "%s/data/output/test__numpy_3d.jpeg"%(modd_str)
    assert os.path.isdir(os.path.dirname(output_file_path_str))

    gf_images_jobs.view_numpy_arr_3D(x_3d,
        output_file_path_str)

    #---------------------------
    # 2D
    mean_int       = 0
    standard_dev_f = 0.5 # 0.002
    x_lst = []
    for i in range(0, 1_000):
        x_lst.append(np.random.normal(mean_int, standard_dev_f, 1_000) * 255.0)

    x = np.array(x_lst, dtype=np.float64)

    # print("sleep 10s...")
    # import time
    # time.sleep(10)

    print("draw...")
    output_file_path_str = "%s/data/output/test__numpy_2d.jpeg"%(modd_str)
    assert os.path.isdir(os.path.dirname(output_file_path_str))

    gf_images_jobs.view_numpy_arr_2D(x,
        output_file_path_str)

    #---------------------------

#---------------------------------------------------------------------------
test__numpy()


#---------------------------------------------------------------------------
def test__collage():


    collage__files_lst       = []
    output_file_path_str = "%s/data/output/test__collage.jpeg"%(modd_str)
    assert os.path.isfile(output_file_path_str)

    for i in range(0, 300):
        collage__files_lst.extend([
            "%s/data/input/test_img_1.jpeg"%(modd_str),
            "%s/data/input/test_img_2.jpeg"%(modd_str),
        ])

    for f in collage__files_lst:
        assert os.path.isfile(f)

    gf_images_jobs.create_collage(collage__files_lst,
        output_file_path_str,
        500,
        500,
        5,  # rows
        5) # columns

#---------------------------------------------------------------------------
test__collage()



#---------------------------------------------------------------------------
def test__transforms():

    img_source_file_path_str = "%s/data/input/test_img_2.jpeg"%(modd_str)
    assert os.path.isfile(img_source_file_path_str)

    # NOISE
    output_f_str = "%s/data/output/test_output__noise.jpeg"%(modd_str)
    assert os.path.isdir(os.path.dirname(output_f_str))

    gf_images_jobs.apply_transforms(["noise"],
        img_source_file_path_str,
        output_f_str)
        
    # CONTRAST
    for i in range(0, 3):

        factor_f     = i * 100.0
        output_f_str = "%s/data/output/test_output__contrast_%s.jpeg"%(modd_str, i)
        assert os.path.isdir(os.path.dirname(output_f_str))

        gf_images_jobs.apply_transforms(["contrast:%s"%(factor_f)],
            img_source_file_path_str,
            output_f_str)


    # SATURATE

    saturation_img_source_file_path_str = "%s/data/input/test_img_1.jpeg"%(modd_str)
    assert os.path.isfile(saturation_img_source_file_path_str)

    for i in range(0, 3):

        factor_f     = i * 0.5
        output_f_str = "%s/data/output/test_output__saturate_%s.jpeg"%(modd_str, i)
        assert os.path.isdir(os.path.dirname(output_f_str))

        gf_images_jobs.apply_transforms(["saturate:%s"%(factor_f)],
            saturation_img_source_file_path_str,
            output_f_str)

#---------------------------------------------------------------------------
test__transforms()