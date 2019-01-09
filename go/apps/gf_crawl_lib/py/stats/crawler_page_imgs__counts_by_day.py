


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
	p_log_fun,
	p_output_img_str = '../plots/crawler_page_imgs__counts_by_day.png'):
		
	#-------------------------------------------------------------
	def query(p_obj_type_str):

		coll    = p_mongo_client['prod_db']['data_symphony']
		results = coll.aggregate([
				{'$match':{
					't':p_obj_type_str}},

				{"$project":{
					"creation_unix_time_f":1,

					#IMPORTANT!! - not using the "domain_str" here because for images thats usually the url 
					#              of the image on some CDN. a lot of domains might have images on the same
					#              CDN domain, which will hide too much meaningful information.
					#              instead we want the domain of the page where the image was discovered.
					#           
					"origin_page_url_domain_str":1
				}},
				{"$group":{

					#group by day
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
							        ]}},},
						'domain_str':'$origin_page_url_domain_str'
			        },
			        "count_int":{"$sum":1},
			    }},

			    {'$sort':{'count_int':-1}},

			    {'$group':{
			    	'_id'        :"$_id.date_str",
			    	"domains_lst":{"$push":{
			    		"domain_str":"$_id.domain_str",
			    		"count_int" :"$count_int" #count for the individual domain in that day
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
	def query__img_attribute(p_img_attribute_str):
		coll    = p_mongo_client['prod_db']['data_symphony']
		results = coll.aggregate([
				{'$match':{'t':'crawler_page_img'}},
				{"$project":{
					"creation_unix_time_f":1,       
					p_img_attribute_str   :1}},
				{"$group":{
					#group by day
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
							        ]}},},
						p_img_attribute_str:'$%s'%(p_img_attribute_str)
			        },
			        "count_int":{"$sum":1},
			    }},
			    {'$sort':{'count_int':-1}},
			    {'$group':{
			    	'_id'      :"$_id.date_str",
			    	"count_int":{"$sum":"$count_int"}, #add the count for each domain on that date
			    }},
			    
			    #IMPORTANT!! - for plotting, oldest records need to be first, timeseries plots go from left to right
			    {'$sort':{'_id':1}},
			],
			allowDiskUse=True)

		print 'DONE - allowDiskUse'
		return results
	#-------------------------------------------------------------

	imgs__results            = query('crawler_page_img')
	imgs_downloaded__results = query__img_attribute('downloaded_bool')
	imgs_s3_stored__results  = query__img_attribute('s3_stored_bool')
	imgs_refs__results       = query('crawler_page_img_ref')
	
	fig = plt.figure(figsize=(30,10))

	#-------------------------------------------------------------

	def process__page_imgs(p_result):
		days_lst  = []
		counts_lst = []
		top_domains_counts_per_day_lst = []
		for r in p_result:

			#pprint.pprint(r)

			days_lst.append(r['_id'])
			counts_lst.append(r['count_int'])

			#------------------
			#IMPORTANT!! - this is relevant because on some days there is a huge number
			#              domains that are crawled, in which case plots become visually
			#              hard to interpret. 

			#IMPORTANT!! - domains are sorted by number of links from that domain,
			#              so first N are the top N most numerous domains.
			domains_lst     = r['domains_lst']
			top_domains_lst = None

			if len(domains_lst) > 10:
				top_domains_lst = domains_lst[:10]
			else:
				top_domains_lst = domains_lst


			top_domains__count_int = 0
			for domain_stats_map in top_domains_lst:

				domain_count_int        = domain_stats_map['count_int']
				top_domains__count_int += domain_count_int
			#------------------

			top_domains_counts_per_day_lst.append(top_domains__count_int)
		return days_lst,counts_lst,top_domains_counts_per_day_lst
	#-------------------------------------------------------------
	def process_page_imgs_attribute(p_result):
		counts_lst = []
		for r in p_result:
			counts_lst.append(r['count_int'])
		return counts_lst
	#-------------------------------------------------------------
	def process_page_imgs_refs(p_result):
		counts_lst = []
		for r in p_result:
			counts_lst.append(r['count_int'])
		return counts_lst
	#-------------------------------------------------------------

	imgs__days_lst,imgs__counts_lst,imgs__top_domains_counts_per_day_lst = process__page_imgs(imgs__results)
	imgs__downloaded__counts_lst                                         = process_page_imgs_attribute(imgs_downloaded__results)
	imgs__s3_stored__counts_lst                                          = process_page_imgs_attribute(imgs_s3_stored__results)
	imgs_refs__counts_lst                                                = process_page_imgs_refs(imgs_refs__results)

	df = pd.DataFrame({
	    "days"                         : imgs__days_lst,
	    "total_counts"                 : imgs__counts_lst,
	    "top_10_domains_counts_per_day": imgs__top_domains_counts_per_day_lst,

	    #CAUTION!! - this assumes that these lists are of the same length as the imgs__days_lst,
	    #            same number of days. if thats not the case Pandas will complain.
	    "imgs__downloaded__counts":imgs__downloaded__counts_lst,
	    "imgs__s3_stored__counts" :imgs__s3_stored__counts_lst,
	    "imgs_refs__counts"       : imgs_refs__counts_lst
	})

	df.set_index("days",drop=True,inplace=True)
	print df

	#casting subject_alt_names_counts_lst to list() first because its a "multiprocessing.managers.ListProxy"
	#count_s = pd.Series(results_lst)
	#l = count_s.value_counts(sort=False)
	#l.sort_index(inplace=True)
	#l.plot.bar(figsize=(10,6))

	df.plot.line(figsize=(10,6),alpha=0.75) #,rot=0)

	plt.title("crawler_page_img's and crawler_page_img_ref's counts per day",fontsize=18)
	plt.xlabel("day",                                                        fontsize=14)
	plt.ylabel('number of imgs/img_refs',                                    fontsize=14)
	plt.xticks(size = 6)
	plt.axes().yaxis.grid() #horizontal-grid

	plt.savefig(p_output_img_str) #save figure to file