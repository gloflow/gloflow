













import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir




print("%s/../../build"%(modd_str))


import delegator
print(delegator.run("ls -al %s"%("%s/../../build"%(modd_str))).out)



import tensorflow as tf


# EAGER_EXECUTION - necessary for calls to TF API's to work outside of 
#                   a graph session run.
tf.compat.v1.enable_eager_execution(
    config         = None,
    device_policy  = None,
    execution_mode = None)

print("TensorFlow version - %s"%(tf.__version__))
assert tf.__version__ == "1.15.0"



print("\n\nRUN THIS TEST WITH - 'LD_LIBRARY_PATH=. python3 test.py'\n\n")
assert "LD_LIBRARY_PATH" in os.environ.keys()

sys.path.append("%s/../../build"%(modd_str))
import gf_images_jobs_py as gf_images_jobs




print("TEST ------------------------------")
print(dir(gf_images_jobs))


import numpy as np

#---------------------------------------------------------------------------
def test__tensorflow():


    print("creating ML dataset...")
    dataset_name_str            = "test"
    dataset_target_dir_path_str = "%s/data/output_ml/generated"%(modd_str)
    gf_images_jobs.generate_ml_dataset_to_tfrecords(dataset_name_str,
        128, 128,
        dataset_target_dir_path_str)



    test__file_path_str = "%s/data/output_ml/gf_test.tfrecord"%(modd_str)

    # dataset contains serialized tf.train.Example messages
    dataset = tf.data.TFRecordDataset(
        test__file_path_str,
        compression_type   = None,
        buffer_size        = None,
        num_parallel_reads = 4) # load the file in parallel
        
    assert isinstance(dataset, tf.data.Dataset)
    print(dataset)
    print(dataset.element_spec)

    def map_f(p_example):
        # assert isinstance(p_example, tf.data.Tensor)

        print("map F")
        print(p_example)


        label_shape_lst = []
        img_shape_lst   = [100, 100]

        # "label" - has to be of type tf.int64. TensorFlow only allows float32, int64, string types.
        #           event though in Rust examples/record are packed as a list of bytes (u8).
        features_def_map = {
            "label": tf.compat.v1.io.FixedLenFeature(label_shape_lst, tf.int64),
            "img":   tf.compat.v1.io.FixedLenFeature(img_shape_lst,   tf.string)
        }

        # FixedLenFeature(shape, dtype) - configuration for parsing a fixed-length input feature
        example = tf.compat.v1.io.parse_single_example(p_example,
            features = features_def_map)

        #print(example)

        return example

    dataset.map(map_f)

    
    # inspect the type of each element component
    print("element spec")
    print(dataset.element_spec)

    for tensor in dataset.take(10):
        print("-------")
        #assert isinstance(raw_record, tf.python.framework.ops.EagerTensor)
        
        
        print(type(tensor))
        print(tensor.dtype)
        print(tensor.shape)
        print(dir(tensor))


        # reshaped_tensor = np.reshape(tensor.numpy(), (1000, 1000))
        # print(reshaped_tensor)


        print(len(tensor.numpy()))


        # print(repr(raw_record))

#---------------------------------------------------------------------------
test__tensorflow()



exit()



x = np.arange(10, dtype=np.float64) # 1_000_000_000, dtype=np.float64)

x_in_gb = x.nbytes/1024/1024/1024
print(x_in_gb)

#---------------------------------------------------------------------------
def test__numpy():

    import random


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

    output_file_path_str = "%s/test/output/test__numpy_4d.jpeg"%(modd_str)
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
    output_file_path_str = "%s/test/output/test__numpy_3d.jpeg"%(modd_str)
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
    output_file_path_str = "%s/test/output/test__numpy_2d.jpeg"%(modd_str)
    gf_images_jobs.view_numpy_arr_2D(x,
        output_file_path_str)

    #---------------------------
#---------------------------------------------------------------------------
test__numpy()




#---------------------------------------------------------------------------
def test__collage():


    collage__files_lst       = []
    output_file_path_str = "%s/test/output/test__collage.jpeg"%(modd_str)
    for i in range(0, 300):
        collage__files_lst.extend([
            "%s/data/input/50b230e8933860a01cd5da61d082887a.jpeg"%(modd_str),
            "%s/data/input/49a180a9ab8548b69f50e0bb2c96b4d0_thumb_small.jpeg"%(modd_str),
        ])

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

    img_source_file_path_str = "%s/data/input/50b230e8933860a01cd5da61d082887a.jpeg"%(modd_str)

    # NOISE
    output_f_str = "%s/test/output/test_output__noise.jpeg"%(modd_str)
    gf_images_jobs.apply_transforms(["noise"],
        img_source_file_path_str,
        output_f_str)
        
    # CONTRAST
    for i in range(0, 3):

        factor_f     = i * 100.0
        output_f_str = "%s/test/output/test_output__contrast_%s.jpeg"%(modd_str, i)
        gf_images_jobs.apply_transforms(["contrast:%s"%(factor_f)],
            img_source_file_path_str,
            output_f_str)


    # SATURATE

    saturation_img_source_file_path_str = "%s/data/input/49a180a9ab8548b69f50e0bb2c96b4d0_thumb_small.jpeg"%(modd_str)
    for i in range(0, 3):

        factor_f     = i * 0.5
        output_f_str = "%s/test/output/test_output__saturate_%s.jpeg"%(modd_str, i)
        gf_images_jobs.apply_transforms(["saturate:%s"%(factor_f)],
            saturation_img_source_file_path_str,
            output_f_str)

#---------------------------------------------------------------------------
test__transforms()