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

import tensorflow as tf

from tensorflow.keras import layers, models
import matplotlib.pyplot as plt

import gf_data





#----------------------------------------------
def main():



    print("d")
    gf_data.load__generated()


    

    # DATA_INPUT
    data_map = gf_data.load__cifar10()

    # CREATE
    model = create()

    # FIT
    train_tf_history = fit(data_map, model)




    evaluate(data_map,
        train_tf_history,
        model)



#----------------------------------------------
# CREATE
def create():


    print(" --- CREATE MODEL ")
    model = models.Sequential()

    # CONV_1
    model.add(layers.Conv2D(32, (3, 3), activation="relu", input_shape=(32, 32, 3)))
    model.add(layers.MaxPooling2D((2, 2)))
    
    # CONV_2
    model.add(layers.Conv2D(64, (3, 3), activation="relu"))
    model.add(layers.MaxPooling2D((2, 2)))

    # CONV_3
    model.add(layers.Conv2D(64, (3, 3), activation="relu"))

    # FLATTEN
    model.add(layers.Flatten())

    # DENSE - 64
    model.add(layers.Dense(64, activation="relu"))
    
    # DENSE - 10
    model.add(layers.Dense(10))



    model.summary()

    # COMPILE
    model.compile(optimizer = "adam",
        loss    = tf.keras.losses.SparseCategoricalCrossentropy(from_logits=True),
        metrics = ["accuracy"])

    return model

#----------------------------------------------
# FIT
def fit(p_data_map,
    p_model):

    print(" --- FIT MODEL ")
    epochs_num_int = 10


    print("run...")
    train_tf_history = p_model.fit(p_data_map["train_images"][:10],
        p_data_map["train_labels"][:10],
        epochs          = epochs_num_int,
        validation_data = (p_data_map["test_images"][:10], p_data_map["test_labels"][:10]))


    print(train_tf_history)




    return train_tf_history

#----------------------------------------------
# EVALUATE
def evaluate(p_data_map,
    p_train_tf_history,
    p_model):



    print(p_train_tf_history.history.keys())



    plt.plot(p_train_tf_history.history["acc"],     label = "accuracy")
    plt.plot(p_train_tf_history.history["val_acc"], label = "val_accuracy")
    plt.xlabel("Epoch")
    plt.ylabel("Accuracy")
    plt.ylim([0.5, 1])
    plt.legend(loc="lower right")

    test_loss, test_acc = p_model.evaluate(p_data_map["test_images"],
        p_data_map["test_labels"],
        verbose=2)





    print(test_loss)
    print(test_acc)



    plt.show()


#----------------------------------------------
main()