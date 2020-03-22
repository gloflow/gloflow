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
from matplotlib import pyplot as plt

# TENSORFLOW
import tensorflow as tf

# EAGER_EXECUTION - necessary for calls to TF API's to work outside of 
#                   a graph session run.
tf.compat.v1.enable_eager_execution(
    config         = None,
    device_policy  = None,
    execution_mode = None)

print("TensorFlow version - %s"%(tf.__version__))
assert tf.__version__ == "1.15.0"


# GF_IMAGES_JOBS
print("\n\nRUN THIS TEST WITH - 'LD_LIBRARY_PATH=../../build/ python3 test.py'\n\n")
assert "LD_LIBRARY_PATH" in os.environ.keys()

sys.path.append("%s/../../build"%(modd_str))
import gf_images_jobs_py as gf_images_jobs

#---------------------------------------------------------------------------
def main():

    print("creating ML dataset...")

    #---------------------------
    # CONFIG
    dataset_name_str            = "test"
    dataset_target_dir_path_str = "%s/data/output_ml/generated"%(modd_str)
    classes_lst                 = ["rect"]
    elements_num_int            = 2000
    img_width_int    = 32
    img_height_int   = 32
    img_channels_int = 4 # RGBA

    test_py_tfrecords_file_str   = "./data/output_ml/gf_py_test.tfrecords"
    test_rust_tfrecords_file_str = "./data/output_ml/gf_rust_test.tfrecords"



    test_ml_tf_records_train__file_str    = "./data/output_ml/generated/test__train.tfrecords"
    test_ml_tf_records_validate__file_str = "./data/output_ml/generated/test__validate.tfrecords"
    #---------------------------

    # GENERATE_ML_DATASET
    gf_images_jobs.generate_ml_dataset(dataset_name_str,
        classes_lst,
        elements_num_int,
        img_width_int,
        img_height_int,
        dataset_target_dir_path_str)

    # GENERATE_AND_REGISTER_ML_DATASET - generates .tfrecords of images and issues
    #                                    HTTP request to "gf_ml" server to register the generated dataset.
    # gf_images_jobs.generate_and_register_ml_dataset(dataset_name_str,
    #     classes_lst,
    #     elements_num_int,
    #     img_width_int,
    #     img_height_int,
    #     dataset_target_dir_path_str)

    
    print("----------------")
    print("test .tfrecords reading")

    tf_example__img_width_int  = img_width_int
    tf_example__img_height_int = img_height_int
    collage__img_width_int     = 1000
    collage__img_height_int    = 1000
    collage__rows_num_int      = 40
    collage__columns_num_int   = 40

    assert os.path.isfile(test_ml_tf_records_train__file_str)
    gf_images_jobs.view_ml_dataset(test_ml_tf_records_train__file_str,
        "./generated_dataset_collage.png",
        tf_example__img_width_int,
        tf_example__img_height_int,
        collage__img_width_int,
        collage__img_height_int,
        collage__rows_num_int,
        collage__columns_num_int)

    exit()


    
    test__tf_record_processing(test_rust_tfrecords_file_str,
        img_width_int,
        img_height_int,
        p_img_channels_int = img_channels_int)
    
    
    # PY_TFRECORDS
    test__py_write_tfrecord(test_py_tfrecords_file_str)
    test__py_read_tfrecord(test_py_tfrecords_file_str, p_view_bool = True)
    
    # RUST_TFRECORDS
    test__py_read_tfrecord(test_rust_tfrecords_file_str, p_view_bool = True)






    



#---------------------------------------------------------------------------
def test__tf_record_processing(p_tfrecord_path_str,
    p_img_width_int,
    p_img_height_int,

    # 4 - channels for RGBA
    p_img_channels_int = 4):
    assert os.path.isfile(p_tfrecord_path_str)

    # dataset contains serialized tf.train.Example messages
    dataset = tf.data.TFRecordDataset(
        p_tfrecord_path_str,
        compression_type   = None,
        buffer_size        = None,
        num_parallel_reads = 4) # load the file in parallel
        
    assert isinstance(dataset, tf.data.Dataset)
    print(dataset)
    print(dataset.element_spec)
    print("before parsing...")

    #---------------------------------------------------------------------------
    def map_f(p_example):

        # assert isinstance(p_example, tf.data.Tensor)
        # print(p_example)

        label_shape_lst = []
        img_shape_lst   = [p_img_width_int, p_img_height_int, p_img_channels_int] 

        # FixedLenFeature(shape, dtype) - configuration for parsing a fixed-length input feature
        features_def_map = {
            "label": tf.compat.v1.io.FixedLenFeature(label_shape_lst, tf.int64),
            "img":   tf.compat.v1.io.FixedLenFeature(img_shape_lst,   tf.string)
        }

        # parse_single_example() - Parses a single Example proto
        example = tf.compat.v1.io.parse_single_example(p_example,
            features = features_def_map,
            name     = "test_example")

        return example

    #---------------------------------------------------------------------------
    parsed_dataset = dataset.map(map_f)
    print("element spec - %s"%(parsed_dataset.element_spec)) # inspect the type of each element component

    w, h, channels_int = parsed_dataset.element_spec["img"].shape
    assert p_img_width_int == w
    assert p_img_height_int == h
    assert channels_int == p_img_channels_int
    
#---------------------------------------------------------------------------
def test__py_write_tfrecord(p_tfrecord_path_str):
    #---------------------------------------------------------------------------
    def _int64_feature(value):
        # Returns an int64_list from a bool / enum / int / uint
        return tf.train.Feature(int64_list=tf.train.Int64List(value=[value]))
    
    #---------------------------------------------------------------------------
    def _bytes_feature(value):
        # Returns a bytes_list from a string / byte
        return tf.train.Feature(bytes_list=tf.train.BytesList(value=[value]))
    
    #---------------------------------------------------------------------------
    # Read image raw data, which will be embedded in the record file later.
    image_string = open("data/output_ml/generated/train/rect/test-rect-0.png", "rb").read()
    
    # Manually set the label to 0. This should be set according to your situation.
    label = 0
    
    # For each sample there are two features: image raw data, and label. Wrap them in a single dict.
    feature = {
        "label": _int64_feature(label),
        "img":   _bytes_feature(image_string),
    }
    
    # Create a `example` from the feature dict.
    tf_example = tf.train.Example(features=tf.train.Features(feature=feature))

    # Write the serialized example to a record file.
    with tf.io.TFRecordWriter(p_tfrecord_path_str) as writer:
        writer.write(tf_example.SerializeToString())

#---------------------------------------------------------------------------
def test__py_read_tfrecord(p_tfrecord_path_str,
    p_view_bool = False):

    # dataset contains serialized tf.train.Example messages
    dataset = tf.data.TFRecordDataset(
        p_tfrecord_path_str,
        compression_type   = None,
        buffer_size        = None,
        num_parallel_reads = 4) # load the file in parallel


    print(dataset)

    for i, example in dataset.enumerate():

        assert isinstance(example, tf.Tensor)
        print("=========================================================================")
        print("unparsed example - %s"%(example)) 

        #assert isinstance(raw_record, tf.python.framework.ops.EagerTensor)
        
        print("=========================-------")
        example_parsed = tf.train.Example.FromString(example.numpy())
        
        print("example parsed (%s):"%(p_tfrecord_path_str))
        print(example_parsed)
        print("=========================-------")

        print("example_parsed type - %s"%(type(example_parsed)))
        print("image feature size  - %s"%(example_parsed.features.feature["img"].ByteSize()))

        #import tensorflow.core.example.example_pb2 as example_pb2
        #assert isinstance(parsed, example_pb2.Example)

        print("=========================-------")
        print("image feature (%s):"%(p_tfrecord_path_str))
        img_feature = example_parsed.features.feature["img"]
        print(img_feature)
        print("=========================-------")

        # DECODE_PNG
        img_bytes = img_feature.bytes_list.value[0]
        img_png = tf.compat.v1.image.decode_png(img_bytes)
        print("PNG decoded...")

        # VIEW_IMG
        if p_view_bool:
            plt.imshow(img_png, interpolation = "nearest")

            print("viewing file - %s"%(p_tfrecord_path_str))
            plt.show()

        # reshaped_tensor = np.reshape(example_parsed.numpy(), (img_width_int, img_height_int))
        # print(reshaped_tensor)
    

#---------------------------------------------------------------------------
main()