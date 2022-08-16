# GloFlow application and media management/publishing platform
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
	return '5m'
	
#-------------------------------------------------------------
def run(p_mongo_client,
	pLogFun,
	p_output_img_str = '../plots/crawler_url_fetches__counts_by_day.png'):
	
	#-------------------------------------------------------------
	def query():
		coll    = p_mongo_client['prod_db']['gf_crawl']
		results = coll.aggregate([
				{'$match': {
					't': 'crawler_url_fetch'}},

				{"$project": {
					"creation_unix_time_f": 1,
					"domain_str":           1
				}},
				{"$group": {

					#group by day
					"_id": {
						"date_str": {
							"$dateToString": {
								"format": "%Y-%m-%d",
								"date":   {

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
						'domain_str': '$domain_str'
					},
					"count_int": {"$sum": 1},
				}},

				{'$sort': {'count_int': -1}},

 				{'$group':{
					'_id':         "$_id.date_str",
					"domains_lst": {"$push":{
						"domain_str": "$_id.domain_str",
						"count_int":  "$count_int" #count for the individual domain in that day
					}},
					"count_int": {"$sum": "$count_int"}, #add the count for each domain on that date
				}},

				#IMPORTANT!! - for plotting, oldest records need to be first, timeseries plots go from left to right
				{'$sort':{'_id':1}},
			],
			allowDiskUse=True)

		print('DONE - allowDiskUse')
		return results
		#print results.explain("executionStats")
		#print coll.explain("executionStats")
	#-------------------------------------------------------------

	results = query()
	fig     = plt.figure(figsize=(30,10))

	days_lst   = []
	counts_lst = []
	top_domains_counts_per_day_lst = []
	for r in results:

		#pprint.pprint(r)

		days_lst.append(r['_id'])
		counts_lst.append(r['count_int'])

		#------------------
		#IMPORTANT!! - this is relevant because on some days there is a huge number
		#              domains that are crawled, in which case plots become visually
		#              hard to interpret. 

		#IMPORTANT!! - domains are sorted by number of links from that domain,
		#              so first N are the top N most numerous domains.
		domains_lst = r['domains_lst']
		top_domains_lst = None
		if len(domains_lst) > 10:
			top_domains_lst = domains_lst[:10]
		else:
			top_domains_lst = domains_lst
		#------------------

		top_domains__count_int = 0
		for domain_stats_map in top_domains_lst:

			domain_str       = domain_stats_map['domain_str']
			domain_count_int = domain_stats_map['count_int']

			top_domains__count_int += domain_count_int
		top_domains_counts_per_day_lst.append(top_domains__count_int)

	print(len(days_lst))
	print(len(counts_lst))
	print(len(top_domains_counts_per_day_lst))

	df = pd.DataFrame({
		"days":                          days_lst,
		"total_counts":                  counts_lst,
		"top_10_domains_counts_per_day": top_domains_counts_per_day_lst
	})

	df.set_index("days", drop=True, inplace=True)
	print(df)

	#casting subject_alt_names_counts_lst to list() first because its a "multiprocessing.managers.ListProxy"
	#count_s = pd.Series(results_lst)
	#l = count_s.value_counts(sort=False)
	#l.sort_index(inplace=True)
	#l.plot.bar(figsize=(10,6))

	df.plot.line(figsize=(10,6),alpha=0.75) #,rot=0)

	plt.title("crawler_url_fetch's counts per day", fontsize=18)
	plt.xlabel("day",                               fontsize=14)
	plt.ylabel('number of fetches',                 fontsize=14)
	plt.xticks(size = 6)
	plt.axes().yaxis.grid() #horizontal-grid

	plt.savefig(p_output_img_str) #save figure to file