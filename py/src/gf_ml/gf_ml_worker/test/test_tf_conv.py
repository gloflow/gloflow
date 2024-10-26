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
import matplotlib.pyplot as plt

import tensorflow as tf
import tensorflow.keras.layers
tf.compat.v1.enable_eager_execution() # EAGER_EXEC

import tensorflow.keras.utils

sys.path.append("%s/utils"%(modd_str))
import gf_ml_test_data

#----------------------------------------------
def main():

    print(tf.executing_eagerly())

    #----------------------------------------------
    def tf_set_cpu_target():
        cpu_devices_lst = tf.config.experimental.list_physical_devices(device_type="CPU")
        print(cpu_devices_lst)

        tf.config.experimental.set_visible_devices(devices=cpu_devices_lst, device_type="CPU")
        
        # enable logging of target devices used for execution of ops
        tf.debugging.set_log_device_placement(True)

    #----------------------------------------------

    tf_set_cpu_target()

    input_img_batch       = gf_ml_test_data.load_data()
    input_img_transformed = np.squeeze(input_img_batch, 0) # drop 1st (batch) dimension


    basic_conv_model(input_img_batch)



    test_conv_layer(input_img_batch)



#----------------------------------------------
def basic_conv_model(p_input_img_batch):
            
    input_l = tensorflow.keras.layers.Input(shape=(64, 64, 3))


    tf_conv_l = tensorflow.keras.layers.Conv2D(filters=8,
        kernel_size = (5, 5),
        padding     = "same",
        activation  = "relu")(input_l)

    maxpool_l = tensorflow.keras.layers.MaxPooling2D((2, 2))(tf_conv_l)




    # o = tf_conv_l(p_input_img_batch)
    model = tf.keras.models.Model(inputs=input_l, outputs=maxpool_l)
    model.summary()

    # MODEL_COMPILE
    model.compile(optimizer = "adam",
        loss    = tf.keras.losses.SparseCategoricalCrossentropy(from_logits=True),
        metrics = ["accuracy"])

    # MODEL_PREDICT
    y = model.predict(p_input_img_batch)
    print(y[0].shape)



    # CONV_LAYER_OUTPUT
    tf_conv_l_ouput_fn = tf.keras.backend.function([model.layers[0].input],
        [model.layers[1].output])




    conv_output = tf_conv_l_ouput_fn(p_input_img_batch)
    print(conv_output[0].shape)


    tf_conv_l_output = tf_conv_l_ouput_fn(p_input_img_batch)[0]





    # PLOT_MODEL
    tensorflow.keras.utils.plot_model(model, to_file="%s/data/output/tf_model.png"%(modd_str))


    #----------------------------------------------
    def plot_imgs():
        
        fig  = plt.figure(constrained_layout=False, figsize=(10, 5))
        spec = fig.add_gridspec(ncols=3, nrows=2)


        ax0 = fig.add_subplot(spec[:, 0]) # across all rows and first column
        ax1 = fig.add_subplot(spec[:, 1]) # across all rows and second column
        ax2 = fig.add_subplot(spec[0, 2]) # accross first row and third column


        # fig, (ax0, ax1, ax2) = fig.subplots(ncols=3, nrows=2, figsize=(10, 5))


        # IMPORTANT!! - transpose() - output of a convolutional layer and maxpool layer is transposed:
        #               (channels_num, filter_w, filter_h, filters_num)
        #               but to get individual filter values, we want:
        #               (filters_num, filter_w, filter_h, channels_num)

        first_filter = tf_conv_l_output.transpose()[0].squeeze()
        first_maxpool = y[0].transpose()[0]

        ax0.set_title("input image $x$", fontsize=8)
        ax0.imshow(p_input_img_batch[0].astype("uint8"))

        ax1.set_title("Conv2D layer output image", fontsize=8)
        ax1.imshow(first_filter.squeeze())

        ax2.set_title("MaxPooling2D layer output image", fontsize=8)
        ax2.imshow(first_maxpool.astype("uint8"))

        plt.show()

    #----------------------------------------------
    plot_imgs()

#----------------------------------------------
def test_conv_layer(p_input_img_batch):

    # CONVOLUTION_LAYER
    tf_conv_l = tensorflow.keras.layers.Conv2D(filters=3,
        kernel_size = (5, 5),
        padding     = "same",
        input_shape = (None, None, 3),
        activation  = "relu")

    img_in_np  = p_input_img_batch[0]
    img_out    = tf_conv_l(p_input_img_batch)
    img_out_np = img_out[0].numpy()

    # PLOT_FILTERS
    plot_conv_layer(tf_conv_l)

    #----------------------------------------------
    def plot_imgs():
        fig, (ax0, ax1) = plt.subplots(ncols=2, figsize=(10, 5))

        ax0.imshow(img_in_np.astype("uint8"))
        ax1.imshow(img_out_np.astype("uint8"))

        plt.show()
    
    #----------------------------------------------
    plot_imgs()


    
#----------------------------------------------
def plot_conv_layer(p_tf_conv_l):




    print("conv params #  - %s"%(p_tf_conv_l.count_params()))
    print("conv filters # - %s"%(len(p_tf_conv_l.get_weights())))

    print(dir(p_tf_conv_l))
    print(p_tf_conv_l.filters)
    print(p_tf_conv_l.rank)



    filters_num_int = p_tf_conv_l.filters

    vars_lst = p_tf_conv_l.get_weights()
    w        = vars_lst[0]
    b        = vars_lst[1]

    print(w.shape)
    print(b.shape)
    


    print(w.shape)


    # w             - (width, height, channels_num, filters_num)
    # w.transpose() - (filters_num, width, height, channels_num)
    # [3, 0, 1, 2]  - move the last axis (filters_num) to be the first
    n = np.transpose(w, axes=[3, 0, 1, 2])


    print(n.shape)
 
    
    print(n.shape)
    print(n[0].shape)

    print(n[0])

    print(p_tf_conv_l)




    # PLOT_FILTERS
    fig, axis_tpl = plt.subplots(ncols=filters_num_int, figsize=(10, 5))
    for i, ax in enumerate(axis_tpl):

        filter_weights = n[i]
        ax.set_title("filter %s"%(i))



        ax.imshow(filter_weights)



    plt.show()

#----------------------------------------------
main()