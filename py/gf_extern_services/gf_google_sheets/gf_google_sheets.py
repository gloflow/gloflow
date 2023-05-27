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

#---------------------------------------------------------------------------------
def get_column_by_name(p_column_name_str,
    p_spreadsheet_id_str,
    p_subsheet_name_str,
    p_service,
    p_column_name__row_index_int=1):

    range_name_str = f"{p_subsheet_name_str}!A:Z"

    result = p_service.spreadsheets().values().get(spreadsheetId=p_spreadsheet_id_str,
        # body=data,
        range=range_name_str,

        # format returned data to be in column-first format 
        majorDimension='COLUMNS').execute()



    # iterate over each column, and get the one thats needed
    i=0
    for column_vals_lst in result['values']:
        
        column_name_str = column_vals_lst[p_column_name__row_index_int]
        if column_name_str == p_column_name_str:
            return column_vals_lst, i

        i+=1
    
    return None, 0

#---------------------------------------------------------------------------------
def get_sheet_id(p_sheet_name_str,
    p_spreadsheet_id_str,
    p_service):

    spreadsheet_metadata = p_service.spreadsheets().get(spreadsheetId=p_spreadsheet_id_str).execute()
    sheets_lst           = spreadsheet_metadata.get('sheets', '')

    for sheet in sheets_lst:
        properties = sheet.get('properties', {})
        if properties.get('title', '').lower() == p_sheet_name_str.lower():
            sheet_id_int = properties.get('sheetId', '')
            return sheet_id_int

    return None