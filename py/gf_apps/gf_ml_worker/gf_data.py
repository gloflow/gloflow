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

import tensorflow as tf
from tensorflow.keras import datasets

#---------------------------
# GF_IMAGES_JOBS

print('''
RUN WITH - 'LD_LIBRARY_PATH=../../build/ python3 gf_simple_model.py'
''')

# os.environ["LD_LIBRARY_PATH"] = "%s/../../../rust/build"%(modd_str)
# assert "LD_LIBRARY_PATH" in os.environ.keys()

# FIX!! - this is super temporary!! as soon as possible fix this, where the 
#         PY extension lib file is not copied explicitly to the current directory.
#         see why using LD_LIBRARY_PATH doesnt work for gf_images_jobs_py.so, 
#         but does work for libtensorflow.so
import delegator
gf_py_libs_path_str = "%s/../../../rust/build"%(modd_str)
delegator.run("cp %s/gf_images_jobs_py.so %s"%(gf_py_libs_path_str, modd_str))

print("loading gf_images_jobs_py.so")
import gf_images_jobs_py as gf_images_jobs

#---------------------------

#----------------------------------------------
# LOAD_GENERATED
def load__generated(p_generate_bool = False):

    print("load generated...")

    dataset_target_dir_path_str = "%s/test/data/output/generated"%(modd_str)
    dataset_train_file_path_str = "%s/test__train.tfrecords"%(dataset_target_dir_path_str)

    # GENERATE
    if p_generate_bool:
        dataset_name_str = "test"
        classes_lst      = [
            "rect",
            "circle"
        ]

        elements_num_int = 100
        img_width_int    = 32
        img_height_int   = 32




        # GENERATE_ML_DATASET
        gf_images_jobs.generate_ml_dataset(dataset_name_str,
            classes_lst,
            elements_num_int,
            img_width_int,
            img_height_int,
            dataset_target_dir_path_str)


    print("loading dataset...")
    print(dataset_train_file_path_str)


    assert os.path.isfile(dataset_train_file_path_str)
    dataset = tf.data.TFRecordDataset([dataset_train_file_path_str],

        # scalar representing the number of files to read in parallel.
        # If greater than one, the records of files read in parallel
        # are outputted in an interleaved order
		num_parallel_reads = 32)

    print(dataset)
    return dataset

#----------------------------------------------
# LOAD__CIFAR10
def load__cifar10():

    (train_images, train_labels), (test_images, test_labels) = datasets.cifar10.load_data()

    # Normalize pixel values to be between 0 and 1
    train_images, test_images = train_images / 255.0, test_images / 255.0


    #----------------------------------------------
    def show_dataset():
        class_names = [
            "airplane",
            "automobile",
            "bird",
            "cat",
            "deer",
            "dog",
            "frog",
            "horse",
            "ship",
            "truck"
        ]

        plt.figure(figsize=(10, 10))
        for i in range(25):

            plt.subplot(5,5,i+1)
            plt.xticks([])
            plt.yticks([])
            plt.grid(False)

            plt.imshow(train_images[i], cmap = plt.cm.binary)

            # The CIFAR labels happen to be arrays, which is why you need the extra index
            plt.xlabel(class_names[train_labels[i][0]])

        plt.show()

    #----------------------------------------------
    # show_dataset()




    data_map = {
        "train_images": train_images,
        "train_labels": train_labels,
        "test_images":  test_images,
        "test_labels":  test_labels
    }



    return data_map