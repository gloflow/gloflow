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

import pandas as pd
import matplotlib.pyplot as plt 
#-------------------------------------------------------------
#called to find out how frequently to run the stat
def freq():
	return '5m'
	
#-------------------------------------------------------------
def run(p_mongo_client,
	pLogFun,
	p_output_img_str = '../plots/crawler_page_outgoing_links__per_crawler.png'):

	fig = plt.figure(figsize=(30,10))

	results = p_mongo_client['prod_db']['gf_crawl'].aggregate([
			{'$match':{'t': 'crawler_page_outgoing_link'}},
			{'$group':{
				'_id'      : '$crawler_name_str',
				'count_int': {'$sum': 1}}
			},
			{'$sort': {'count_int': -1}}
		],
		allowDiskUse=True)

	print('DONE')

	names_lst  = []
	counts_lst = []
	for r in results:
		names_lst.append(r['_id'])
		counts_lst.append(r['count_int'])


	df = pd.DataFrame({
	    "name":         names_lst,
	    "links_counts": counts_lst
	})

	df.set_index("name",drop=True,inplace=True)
	print(df)

	#casting subject_alt_names_counts_lst to list() first because its a "multiprocessing.managers.ListProxy"
	#count_s = pd.Series(results_lst)
	#l = count_s.value_counts(sort=False)
	#l.sort_index(inplace=True)
	#l.plot.bar(figsize=(10,6))

	df.plot.bar(figsize=(10,6),alpha=0.75, rot=0)

	plt.title("crawler_page_outgoing_link's per crawler_name",fontsize=18)
	plt.xlabel("crawler_name's",                              fontsize=14)
	plt.ylabel('number of links',                             fontsize=14)
	plt.xticks(size = 6)
	plt.axes().yaxis.grid() #horizontal-grid

	plt.savefig(p_output_img_str) #save figure to file
	