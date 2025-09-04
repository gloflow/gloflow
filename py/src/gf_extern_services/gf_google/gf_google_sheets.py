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

import json
from google.oauth2 import service_account
from apiclient import discovery
import pandas as pd

#---------------------------------------------------------------------------------
# NAMED_RANGES
#---------------------------------------------------------------------------------
# SET

def set_named_range(p_value_str,
	p_range_name_str,
	p_spreadsheet_id_str,
	p_service_client):
	"""
	Set a named range in a Google Sheet with a string value.

	Parameters:
	- p_value_str: The string value to set in the range.
	- p_range_name_str: The named range to set.
	- p_spreadsheet_id_str: The spreadsheet ID.
	- p_service_client: The authorized Google Sheets API client.
	"""

	body = {
		"range": p_range_name_str,
		"majorDimension": "COLUMNS",
		"values": [[p_value_str]] # single cell with the provided string value
	}

	# UPDATE
	p_service_client.spreadsheets().values().update(
		spreadsheetId=p_spreadsheet_id_str,
		range=p_range_name_str,

		# "USED_ENTERED" - processes the input as if it were entered directly into the cell by a user.
		#                  alternative is "RAW" - processes the input as if it were a literal string.
		valueInputOption='USER_ENTERED',
		body=body
	).execute()

#---------------------------------------------------------------------------------
# GET

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
# ADD_NOTE_TO_NAMED_RANGE

def add_note_to_named_range(p_note_str,
	p_range_name_str,
	p_spreadsheet_id_str,
	p_service_client):
	"""
	Add a note to the first cell of a named range in a Google Sheet.

	Parameters:
	- p_note_str: The note to add to the named range.
	- p_range_name_str: The named range to add a note to.
	- p_spreadsheet_id_str: The spreadsheet ID.
	- p_service_client: The authorized Google Sheets API client.
	"""

	sheet_metadata = p_service_client.spreadsheets().get(
		spreadsheetId=p_spreadsheet_id_str
	).execute()

	ranges = sheet_metadata.get('namedRanges', [])
	target_range = next((r for r in ranges if r['name'] == p_range_name_str), None)

	if target_range:
		sheet_id_str = target_range['range']['sheetId']
		start_row = target_range['range']['startRowIndex']
		start_col = target_range['range']['startColumnIndex']

		# Set the note only for the first cell in the range
		note_body = {
			"requests": [
				{
					"updateCells": {
						"range": {
							"sheetId": sheet_id_str,
							"startRowIndex": start_row,
							"endRowIndex": start_row + 1,
							"startColumnIndex": start_col,
							"endColumnIndex": start_col + 1
						},
						"rows": [
							{
								"values": [
									{
										"note": p_note_str
									}
								]
							}
						],
						"fields": "note"
					}
				}
			]
		}

		p_service_client.spreadsheets().batchUpdate(
			spreadsheetId=p_spreadsheet_id_str,
			body=note_body
		).execute()
	else:
		raise ValueError("Named range not found.")

#---------------------------------------------------------------------------------
def add_hyperlink_above_named_range(p_url_str,
	p_range_name_str,
	p_spreadsheet_id_str,
	p_service_client,
	p_label_str="link"):
	"""
	Add a clickable hyperlink in the cell one row above the first cell of a named range in a Google Sheet.

	Parameters:
	- p_url_str: The URL to link to.
	- p_range_name_str: The named range to add the hyperlink above.
	- p_spreadsheet_id_str: The spreadsheet ID.
	- p_service_client: The authorized Google Sheets API client.
	- p_label_str: The label for the hyperlink (default is 'Click Here').
	"""

	hyperlink_formula_str = f'=HYPERLINK("{p_url_str}", "{p_label_str}")'

	sheet_metadata = p_service_client.spreadsheets().get(
		spreadsheetId=p_spreadsheet_id_str
	).execute()

	ranges = sheet_metadata.get('namedRanges', [])
	target_range = next((r for r in ranges if r['name'] == p_range_name_str), None)

	if target_range:

		sheet_id_int = target_range['range']['sheetId']
		start_row = target_range['range']['startRowIndex']
		start_col = target_range['range']['startColumnIndex']

		sheet_name_str = get_sheet_name(p_spreadsheet_id_str, sheet_id_int, p_service_client)

		col_letter = chr(65 + start_col) # Convert column index to letter (A, B, C...)
		cell_address_str = f"{sheet_name_str}!{col_letter}{start_row}"
		print(f"Adding hyperlink to cell: {cell_address_str}")

		p_service_client.spreadsheets().values().update(
			spreadsheetId=p_spreadsheet_id_str,
			range=cell_address_str,
			valueInputOption='USER_ENTERED',
			body={
				"range": cell_address_str,
				"majorDimension": "COLUMNS",
				"values": [[hyperlink_formula_str]]
			}
		).execute()
	else:
		raise ValueError("Named range not found.")

#---------------------------------------------------------------------------------
# VAR
#---------------------------------------------------------------------------------
def get_client(p_google_service_key_jobs_json_str):

	SCOPES = ['https://www.googleapis.com/auth/spreadsheets']

	print("google-sheets connecting...")
	credentials = service_account.Credentials.from_service_account_info(json.loads(p_google_service_key_jobs_json_str),
		scopes=SCOPES)
	service = discovery.build("sheets", "v4", credentials=credentials)
	return service

#---------------------------------------------------------------------------------
# IMPORTANT!! - downloads the whole sheet, all the values of columns

def get_all_columns(p_spreadsheet_id_str,
	p_subsheet_name_str,
	p_service_client):

	range_name_str = f"{p_subsheet_name_str}"

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

#---------------------------------------------------------------------------------
def get_sheet_name(p_spreadsheet_id_str,
	p_sheet_id_int,
	p_service_client):
	sheets_metadata = p_service_client.spreadsheets().get(spreadsheetId=p_spreadsheet_id_str).execute()
	for sheet in sheets_metadata['sheets']:
		if sheet['properties']['sheetId'] == p_sheet_id_int:
			return sheet['properties']['title']
	return None

#---------------------------------------------------------------------------------
def get_subcolumns_by_names(p_supercolumn_name_str,
	p_subcolumn_names_lst,
	p_spreadsheet_id_str,
	p_subsheet_name_str,
	p_service_client,
	p_supercolumn_row_index_int=0,
	p_subcolumn_row_index_int=1):
	"""
	Get multiple subcolumns by name that are located within a supercolumn.

	Parameters:
	- p_supercolumn_name_str: The name of the supercolumn (merged column)
	- p_subcolumn_names_lst: List of names of the subcolumns within the supercolumn
	- p_spreadsheet_id_str: The spreadsheet ID
	- p_subsheet_name_str: The name of the sheet
	- p_service_client: The authorized Google Sheets API client
	- p_supercolumn_row_index_int: Row index where the supercolumn name is located (default: 0)
	- p_subcolumn_row_index_int: Row index where the subcolumn name is located (default: 1)

	Returns:
	- Dictionary where keys are subcolumn names and values are dictionaries with "values" (column values)
	  and "column_index" (column index)
	"""
	assert isinstance(p_supercolumn_name_str, str), "p_supercolumn_name_str must be a string"
	assert isinstance(p_subcolumn_names_lst, list), "p_subcolumn_names_lst must be a list"
	assert isinstance(p_spreadsheet_id_str, str), "p_spreadsheet_id_str must be a string"
	assert isinstance(p_subsheet_name_str, str), "p_subsheet_name_str must be a string"

	# Get all columns from the sheet
	columns_vals_lst = get_all_columns(p_spreadsheet_id_str,
		p_subsheet_name_str,
		p_service_client)

	# Get merged cells information to handle supercolumns properly
	spreadsheet_metadata = p_service_client.spreadsheets().get(
		spreadsheetId=p_spreadsheet_id_str,
		includeGridData=False
	).execute()

	# Find the sheet ID
	sheet_id_int = None
	for sheet in spreadsheet_metadata.get('sheets', []):
		if sheet.get('properties', {}).get('title') == p_subsheet_name_str:
			sheet_id_int = sheet.get('properties', {}).get('sheetId')
			break

	if sheet_id_int is None:
		print(f"Sheet '{p_subsheet_name_str}' not found in spreadsheet '{p_spreadsheet_id_str}'.")
		return {}

	# Get the merged ranges for the target sheet.
	# Explanation:
	# In Google Sheets, when cells are merged, the API stores information about these merged regions
	# in the sheet's metadata under the "merges" key. Each merged range is represented as a dictionary
	# with start/end row and column indices. This is crucial for identifying supercolumns, because only
	# the top-left cell of a merged region contains the displayed value, while the rest are empty.
	# By retrieving these ranges, we can determine which columns are part of a merged supercolumn,
	# and thus correctly locate subcolumns within that supercolumn.
	merged_ranges_lst = []
	for sheet in spreadsheet_metadata.get('sheets', []):
		if sheet.get('properties', {}).get('sheetId') == sheet_id_int:
			merged_ranges_lst = sheet.get('merges', [])
			break

	# First, find the column that contains the supercolumn name
	# In merged cells, only the leftmost cell contains the value
	supercolumn_start_col_int = None
	i = 0
	for column_vals_lst in columns_vals_lst:
		if (len(column_vals_lst) > p_supercolumn_row_index_int and
			column_vals_lst[p_supercolumn_row_index_int] == p_supercolumn_name_str):
			supercolumn_start_col_int = i
			break
		i += 1

	# If we didn't find the supercolumn name, return empty dict
	if supercolumn_start_col_int is None:
		print(f"Supercolumn '{p_supercolumn_name_str}' not found in the specified row index.")
		return {}

	#---------------------------------------------------------------------------------
	def get_supercolumn_columns():

		# Find the merged range that contains this column at the supercolumn row
		# This will help identify all columns that are part of the merged supercolumn
		supercolumn_cols_lst = [supercolumn_start_col_int]  # Start with the column that has the name

		for merged_range in merged_ranges_lst:
			start_row_index = merged_range.get('startRowIndex', 0)
			end_row_index = merged_range.get('endRowIndex', 0)
			start_col_index = merged_range.get('startColumnIndex', 0)
			end_col_index = merged_range.get('endColumnIndex', 0)

			# Check if this merged range is at the supercolumn row and includes our starting column
			if (start_row_index <= p_supercolumn_row_index_int < end_row_index and
				start_col_index <= supercolumn_start_col_int < end_col_index):

				# This merged range represents our supercolumn
				# Add all columns in this range
				for col_idx in range(start_col_index, end_col_index):
					if col_idx not in supercolumn_cols_lst:
						supercolumn_cols_lst.append(col_idx)

		return supercolumn_cols_lst

	#---------------------------------------------------------------------------------
	supercolumn_cols_lst = get_supercolumn_columns()

	# Initialize an empty DataFrame with columns in the order of supercolumn_cols_lst
	results_df = pd.DataFrame()

	# Keep track of which subcolumn names we've found
	found_subcolumns_set = set()

	# Now, within the identified supercolumn columns, find all specified subcolumns
	for col_index_int in supercolumn_cols_lst:

		if col_index_int >= len(columns_vals_lst):
			print(f"Column index {col_index_int} is out of range for the columns list.")
			continue

		column_vals_lst = columns_vals_lst[col_index_int]

		if len(column_vals_lst) <= p_subcolumn_row_index_int:
			print(f"Column {col_index_int} does not have enough rows to contain subcolumn names.")
			continue

		subcolumn_name_str = column_vals_lst[p_subcolumn_row_index_int]

		if subcolumn_name_str in p_subcolumn_names_lst:
			# Extract only the values after the header rows
			# Add 1 to p_subcolumn_row_index_int to start from the row after the subcolumn name
			data_start_index = p_subcolumn_row_index_int + 1
			data_values_lst = column_vals_lst[data_start_index:] if data_start_index < len(column_vals_lst) else []

			# Add the column to the DataFrame
			results_df[subcolumn_name_str] = pd.Series(data_values_lst)

			# Mark this subcolumn as found
			found_subcolumns_set.add(subcolumn_name_str)

			# If we've found all requested subcolumns, we can stop searching
			if len(found_subcolumns_set) == len(p_subcolumn_names_lst):
				break

	# Reorder the DataFrame columns to match the order in supercolumn_cols_lst
	results_df = results_df[[name for name in p_subcolumn_names_lst if name in results_df.columns]]

	return results_df
