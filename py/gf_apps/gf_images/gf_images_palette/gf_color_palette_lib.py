# GloFlow application and media management/publishing platform
# Copyright (C) 2023 Ivan Trajkovic
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

from collections import Counter

from PIL import Image
import numpy as np
from scipy.cluster.vq import kmeans, vq
import webcolors

#----------------------------------------------
# RUN
# gets the color palette for an image by using external libs

def run(p_image_path_str,
    p_num_clusters=3):
    
    image = Image.open(p_image_path_str)
    image = image.resize((150, 150)) # resize the image for faster processing

    # force convert image to RGB image always, in case its a RGBA
    image = image.convert('RGB')

    # convert image to RGB array
    image_array = np.array(image)
    image_array = image_array.reshape((image_array.shape[0] * image_array.shape[1], 3))

    # FIND_CLUSTERS
    centroids, _ = kmeans(image_array.astype(float), p_num_clusters)

    # VECTOR_QUANTIZATION
    # vq() - assigns each observation in a given dataset to the nearest centroid of the clusters
    #   identified by k-means; it maps each data point to the nearest cluster center.
    # cluster - array of indices indicating the closest centroid for each observation in the dataset.
    cluster, _ = vq(image_array, centroids)

    # count frequencies of clusters, to determine which cluster is the most dominant,
    # as a value of the RGB of each pixel.
    cluster_count = Counter(cluster)

    # sort clusters according to their frequency
    dominant_clusters = cluster_count.most_common(p_num_clusters)

    # get the colors for the top clusters
    top_colors = [tuple(centroids[cluster_idx].astype(int)) for cluster_idx, _ in dominant_clusters]



    


    return top_colors

#----------------------------------------------
def get_color_name(p_rgb_color):
    try:
        return webcolors.rgb_to_name(p_rgb_color)
    except ValueError:
        return closest_color(p_rgb_color)

#----------------------------------------------
def closest_color(p_color):

    delta_to_css3_colors_map = {}
    for key, name_str in webcolors.CSS3_HEX_TO_NAMES.items():
        
        r_c, g_c, b_c = webcolors.hex_to_rgb(key)
        
        rd = (r_c - p_color[0]) ** 2
        gd = (g_c - p_color[1]) ** 2
        bd = (b_c - p_color[2]) ** 2

        delta_to_css3_colors_map[(rd + gd + bd)] = name_str

    color_name_str = delta_to_css3_colors_map[min(delta_to_css3_colors_map.keys())]

    return color_name_str

