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

import os
from urllib.parse import urlparse
import requests

#---------------------------------------------------------------------------------
def download_file(p_url_str, p_local_file_path_str):
    response = requests.get(p_url_str)
    with open(p_local_file_path_str, 'wb') as f:
        f.write(response.content)

#---------------------------------------------------------------------------------
def download_file_chunked(p_url_str,
    p_local_file_path_str):
    """Download an image from a URL and save it to the local filesystem."""
    
    # HTTP_GET
    response = requests.get(p_url_str, stream=True)

    if response.status_code == 200:
        
        # ensure the local directory exists
        if not os.path.dirname(p_local_file_path_str) == "":
            os.makedirs(os.path.dirname(p_local_file_path_str), exist_ok=True)

        # open the local file in binary write mode
        with open(p_local_file_path_str, 'wb') as file:

            for chunk in response.iter_content(1024):
                file.write(chunk)
        
        print(f"image downloaded: {p_local_file_path_str}")
    else:
        print(f"failed to download image. Status code: {response.status_code}")

#---------------------------------------------------------------------------------
def is_absolute_url(url):
    parsed_url = urlparse(url)
    return bool(parsed_url.netloc)