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

import tensorflow as tf
import tensorflow.keras.layers
tf.compat.v1.enable_eager_execution() # EAGER_EXEC

sys.path.append("%s/.."%(modd_str))
import gf_ml_data

#----------------------------------------------
def main():
    

    examples_num_int = 10
    epochs_num_int   = 1
    batch_size_int   = 1

    # INPUT_IMG
    # input_img_batch, _ = gf_ml_test_data.download_test_img()
    imgs_train, imgs_test, labels_train, labels_test = gf_ml_data.load__mnist()






    input_l   = tensorflow.keras.layers.Input(shape=(64,))
    dense_1_l = tensorflow.keras.layers.Dense(32, activation="relu", name="layer_1")(input_l)
    dense_2_l = tensorflow.keras.layers.Dense(30, activation="relu", name="layer_2")(dense_1_l)
    dense_3_l = tensorflow.keras.layers.Dense(28, activation="relu", name="layer_3")(dense_2_l)


    dense_softmax_l = tensorflow.keras.layers.Dense(10, activation="softmax", name="layer_4")(dense_3_l)


    


    model = tf.keras.models.Model(inputs=input_l, outputs=dense_softmax_l)
    model.summary()

    # MODEL_COMPILE
    model.compile(optimizer = "adam",
        loss    = tf.keras.losses.SparseCategoricalCrossentropy(from_logits=True),
        metrics = ["accuracy"])

    




    #----------------------------------------------
    def on_batch_end_fun(batch, logs):

        # print("---------------")


        batch_size_int = logs["size"]
        batch_loss_f   = logs["loss"]
        acc_loss_f     = logs["acc"]

        print("loss - %s"%(batch_loss_f))


        # [1:] - exlclude input layer, it does not have weights
        for l in model.layers[1:]:
            # print(type(l))

            # print(l.weights)

            w_var, b_var = l.weights

            print(w_var.shape)

        # print(type(dense_1_l))

        w_var, b_var = model.get_layer("layer_1").weights
        # print(w_var)
        
        # print(w_var.shape)



    #----------------------------------------------
    
    # FIT
    train_tf_history = model.fit(imgs_train[:examples_num_int],
        labels_train[:examples_num_int],
        batch_size      = batch_size_int,
        epochs          = epochs_num_int,
        validation_data = (imgs_test[:examples_num_int], labels_test[:examples_num_int]),
        callbacks       = [tf.keras.callbacks.LambdaCallback(on_batch_end=on_batch_end_fun)])









    # MODEL_PREDICT

    print(imgs_test.shape)
    print(imgs_test[0].shape)

    y = model.predict(imgs_test[:1])

    print(y[0])
    print("pred - %s"%(np.argmax(y[0])))
    print("true - %s"%(labels_test[0]))





#----------------------------------------------
main()