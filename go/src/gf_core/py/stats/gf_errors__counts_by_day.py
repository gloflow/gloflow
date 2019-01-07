# GloFlow media management/publishing system
# Copyright (C) 2019 Ivan Trajkovic
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

import datetime
import pprint
import pandas as pd
import matplotlib.pyplot as plt 

#-------------------------------------------------------------
#called to find out how frequently to run the stat
def freq():
	return '1m'
#-------------------------------------------------------------
def run(p_mongo_client,
	p_log_fun,
	p_output_img_str = '../plots/stat__crawler_url_fetches__counts_by_day.png'):

	#-------------------------------------------------------------
	def query():
		coll    = p_mongo_client['prod_db']['data_symphony']
		results = coll.aggregate([
				{'$match':{
					't':'gf_error'}},

				{"$project":{
					"creation_unix_time_f":1,
					"type_str":            1,
					"service_name_str":    1,
					"subsystem_name_str":  1,
					"func_name_str":       1,
				}},
				{"$group":{

					#group by day
					#composite ID to group by - _id.date_str,_id.subsystem_name_str
					"_id":{
						"date_str":{
							"$dateToString": {
							    "format":"%Y-%m-%d",
							    "date"  :{

							    	#convert creation_unix_time_f from seconds to a Date object
							        "$add":[
							            datetime.datetime(1970,1,1), #new Date(0), 

							            #convert creation_unix_time_f seconds to millisecond timestamp first through
							            #multiplying the value by 1000.
							            {"$multiply": [1000, "$creation_unix_time_f"]}
							        ]
							    }
							},
						},
						'subsystem_name_str':'$subsystem_name_str'
			        },
			        'count_int':{'$sum':1}, #gf_error count for a particular day
			    }},

			    #sort all _id.date_str,_id.subsystem_name_str groups by the number of gf_errors in them,
			    #so that once they're grouped by day, subsystems are already sorted by their counts (number
			    #of gf_errors in them). 
			    {'$sort':{'count_int':-1}},

			    {'$group':{
			    	'_id':                  "$_id.date_str",
			    	"subsystems_errors_lst":{"$push":{
			    		"subsystem_name_str":"$_id.subsystem_name_str",
			    		"count_int":         "$count_int" #count for the individual domain in that day
			    	}},
			    	"count_int":{"$sum":"$count_int"}, #add the count for each domain on that date
			    }},

			    #IMPORTANT!! - for plotting, oldest records need to be first, timeseries plots go from left to right
			    {'$sort':{'_id':1}},
			],
			allowDiskUse=True)

		print 'DONE - allowDiskUse'
		return results

		#print results.explain("executionStats")
		#print coll.explain("executionStats")
	#-------------------------------------------------------------
	
	results = query()
	df_lst  = []

	for r in results:
		day_str                     = r['_id']
		gf_error_counts_per_day_int = r['count_int']

		d = {
			'day':                    day_str,
			'gf_error_counts_per_day':gf_error_counts_per_day_int,
		}

		print '============= ---------------'

		for subsystem_errors_map in r['subsystems_errors_lst']:

			print subsystem_errors_map

			subsystem_name_str            = subsystem_errors_map['subsystem_name_str']
			subsystem_gf_errors_count_int = subsystem_errors_map['count_int']

			assert isinstance(subsystem_name_str,           basestring)
			assert isinstance(subsystem_gf_errors_count_int,int)

			d['subsys__%s'%(subsystem_name_str)] = subsystem_gf_errors_count_int
	
		df_lst.append(d)

	df = pd.DataFrame(df_lst)

	df.set_index("day",drop=True,inplace=True)
	print df

	#---------------------------------------------
	#PLOT

	fig = plt.figure(figsize=(30,10))

	#casting subject_alt_names_counts_lst to list() first because its a "multiprocessing.managers.ListProxy"
	#count_s = pd.Series(results_lst)
	#l = count_s.value_counts(sort=False)
	#l.sort_index(inplace=True)
	#l.plot.bar(figsize=(10,6))

	df.plot.line(figsize=(10,6),alpha=0.75) #,rot=0)
	
	plt.title("gf_errors's counts per day (per subsystem)",fontsize=18)
	plt.xlabel("day",                                      fontsize=14)
	plt.ylabel('number of gf_errors',                      fontsize=14)
	plt.xticks(size = 6)
	plt.axes().yaxis.grid() #horizontal-grid

	plt.savefig(p_output_img_str) #save figure to file
	#---------------------------------------------