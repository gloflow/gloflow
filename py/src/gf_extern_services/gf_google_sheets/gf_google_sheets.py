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
# You should have received a copy of the GNU G11eneral Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

#---------------------------------------------------------------------------------
# IMPORTANT!! - downloads the whole sheet, all the values of columns

def get_all_columns(p_spreadsheet_id_str,
    p_subsheet_name_str,
    p_service_client):
    
    range_name_str = f"{p_subsheet_name_str}!A:Z"

    result = p_service_client.spreadsheets().values().get(spreadsheetId=p_spreadsheet_id_str,
        # body=data,
        range=range_name_str,

        # format returned data to be in column-first format 
        majorDimension='COLUMNS').execute()
    
    # iterate over each column, and get the one thats needed
    columns_vals_lst = []
    for column_vals_lst in result['values']:
        
        columns_vals_lst.append(column_vals_lst)

    return columns_vals_lst


#---------------------------------------------------------------------------------
# NAMED_RANGE

def get_named_range(p_range_name_str,
    p_spreadsheet_id_str,
    p_service_client):


    result = p_service_client.spreadsheets().values().get(spreadsheetId=p_spreadsheet_id_str,

        range=p_range_name_str,

        # format returned data to be in column-first format 
        majorDimension='COLUMNS').execute()
    
    # iterate over each column, and get the one thats needed.
    columns_vals_lst = []
    for column_vals_lst in result['values']:
        
        if len(column_vals_lst) > 0:
            columns_vals_lst.append(column_vals_lst)

    return columns_vals_lst

#---------------------------------------------------------------------------------
def get_column_by_name_from_list(p_column_name_str,
    p_columns_vals_lst,
    p_column_name__row_index_int=1):
    assert isinstance(p_columns_vals_lst, list)

    # iterate over each column, and get the one thats needed
    i=0
    for column_vals_lst in p_columns_vals_lst:

        # first test if the column has values,
        # and if it has more values than the row number where the column name is stored
        if len(column_vals_lst) > 0 and len(column_vals_lst) > (p_column_name__row_index_int+1):
            
            column_name_str = column_vals_lst[p_column_name__row_index_int]
            if column_name_str == p_column_name_str:
                return column_vals_lst, i

        i+=1

    return None, 0
    
#---------------------------------------------------------------------------------
# block is a user-created sub-section of a sheet.

def get_columns_by_name_from_named_range(p_columns_names_lst,
    p_range_name_str,
    p_spreadsheet_id_str,
    p_service_client,
    p_column_name__row_index_int=0):


    columns_lst = get_named_range(p_range_name_str,
        p_spreadsheet_id_str,
        p_service_client)


    columns_vals_lst = []
    for column_name_str in p_columns_names_lst:

        i=0
        for vals_lst in columns_lst:
             
            if len(vals_lst) > 0 and len(vals_lst) > (p_column_name__row_index_int):
                
                a_column_name_str = vals_lst[p_column_name__row_index_int]

                if column_name_str == a_column_name_str:
                    column_vals_lst = vals_lst
                    columns_vals_lst.append((column_vals_lst, i))

                    # column discovered, go to the next one
                    break

            i+=1

    return columns_vals_lst

#---------------------------------------------------------------------------------
def get_column_by_name(p_column_name_str,
    p_spreadsheet_id_str,
    p_subsheet_name_str,
    p_service_client,

    # row in which the user-titles of columns are expected
    p_column_name__row_index_int=1):


    # get all values of all the columns in a sheet
    columns_vals_lst = get_all_columns(p_spreadsheet_id_str,
        p_subsheet_name_str,
        p_service_client)

    i=0

    # iterate over each column, and get the one thats needed
    for column_vals_lst in columns_vals_lst:
        
        # first test if the column has values,
        # and if it has more values than the row number where the column name is stored
        if len(column_vals_lst) > 0 and len(column_vals_lst) > (p_column_name__row_index_int+1):
            
            column_name_str = column_vals_lst[p_column_name__row_index_int]

            if column_name_str == p_column_name_str:
                return column_vals_lst, i

        i+=1
    
    return None, 0

#---------------------------------------------------------------------------------
def get_sheet_id(p_sheet_name_str,
    p_spreadsheet_id_str,
    p_service_client):

    spreadsheet_metadata = p_service_client.spreadsheets().get(spreadsheetId=p_spreadsheet_id_str).execute()
    sheets_lst           = spreadsheet_metadata.get('sheets', '')

    for sheet in sheets_lst:
        properties = sheet.get('properties', {})
        if properties.get('title', '').lower() == p_sheet_name_str.lower():
            sheet_id_int = properties.get('sheetId', '')
            return sheet_id_int

    return None