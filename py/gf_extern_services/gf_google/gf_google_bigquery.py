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

from google.cloud import bigquery

#--------------------------------------------
def table_exists(p_table_name_str,
    p_dataset_id_str,
    p_bigquery_client):

    table_ref = p_bigquery_client.dataset(p_dataset_id_str).table(p_table_name_str)
    try:
        table = p_bigquery_client.get_table(table_ref)

        # Table exists, do something...
        print(f"Table '{table_ref}' exists!")
        return True

    except Exception as e:

        if "Not found" in str(e):
            print(f"Table '{table_ref}' doesn't exist.")
            return False