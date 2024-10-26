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

# IMPORTANT!! - this module is run by GloFlow platform, not directly by user.
#       its run by the GF golang code, and its output thats prefixed 
#       with "GF_OUT:" is parsed by the GF golang code.

import os
import json
import argparse

import sentry_sdk

#----------------------------------------------
def run():

    print("\nGF_IMAGES_CLASSIFY >>>>> PY \n")

    #--------------------
    # INPUT
    
    parser = argparse.ArgumentParser(description='image GF IDs to classify')
    parser.add_argument('-images_ids', type=str, help='list of image GF IDs to classify, comma-separated')
    args = parser.parse_args()

    sentry_dsn_str = os.getenv('SENTRY_DSN')
    sentry_env_str = os.getenv('SENTRY_ENV')

    #--------------------
    # SENTRY

    if not sentry_dsn_str == None:
        print("sentry enabled...")
        sentry_sdk.init(
            dns=sentry_dsn_str,
            environment=sentry_env_str,
            traces_sample_rate=1.0
        )

    #--------------------
    images_ids_lst = args.images_ids.split(',')


    print("images classify...")
    print(f"image id's: {images_ids_lst}")



    classes_lst = ['cat', 'dog', 'bird', 'fish']

    output_str = f"GF_OUT:{json.dumps({'classes_lst': classes_lst})}" 
    print(output_str)

#----------------------------------------------


if __name__ == "__main__":
    run()