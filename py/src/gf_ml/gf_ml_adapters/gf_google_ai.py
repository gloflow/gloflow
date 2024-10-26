# GloFlow application and media management/publishing platform
# Copyright (C) 2024 Ivan Trajkovic
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

import os

from google.oauth2 import service_account
from google.cloud.aiplatform.gapic.schema import predict
from google.cloud import vision
# from google.cloud import aiplatform

#--------------------------------------------
# DESCRIBE_IMAGE

def describe_image(p_url_str, p_ai_client):
    
    print(f"describing image {p_url_str} ...")

    image = vision.Image()
    image.source.image_uri = p_url_str

    response = p_ai_client.label_detection(image=image)
    labels = response.label_annotations

    print("Labels:")
    labels_lst = []
    for label in labels:
        
        print(">+>>>>>>>>>", label.score, label.description)
        
        label_str = label.description
        labels_lst.append((label_str, label.score))
    return labels_lst
        
#--------------------------------------------
# CLASSIFY_IMAGE
# not used currently

def classify_images(p_image_path_local_str,
    p_project_id_str,
    p_model_id_str,
    p_client):
    assert os.path.isfile(p_image_path_local_str)

    # MODEL
    model_name = p_client.model_path(p_project_id_str, "us-central1", p_model_id_str)

    with open(p_image_path_local_str, "rb") as f:
        file_content = f.read()

    request = predict.instance.ImageClassificationPredictionInstance(
        content=file_content
    ).to_value()

    instances = [
        request
    ]

    # PREDICT
    response = p_client.predict(endpoint=model_name, instances=instances)

    # process responses
    for prediction in response.predictions:
        print("Prediction result:", prediction)

#--------------------------------------------
def init(p_key_info_map):
    
    credentials  = service_account.Credentials.from_service_account_info(p_key_info_map)
    
    # client = bigquery.Client(credentials=credentials, project=credentials.project_id)
    # client = aiplatform.gapic.PredictionServiceClient(credentials=credentials)
    client = vision.ImageAnnotatorClient(credentials=credentials)
    
    return client