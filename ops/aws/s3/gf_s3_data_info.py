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
import boto3
from colored import fg,bg,attr

import gf_s3_utils
#---------------------------------------------------
def stats__image_buckets_general(p_aws_access_key_id_str,
	p_aws_secret_access_key_str):

	gf_s3_info = gf_s3_utils.s3_connect(p_aws_access_key_id_str, p_aws_secret_access_key_str)
	assert isinstance(gf_s3_info, gf_s3_utils.Gf_s3_info)

	#LIST_BUCKETS
	for bucket in gf_s3_info.s3_resource.buckets.all():
		print(bucket.name)

	#-------------------
	main_images__bucket_info_map    = process_bucket('gf--img',             gf_s3_info.imgs__bucket)
	crawler_images__bucket_info_map = process_bucket('gf--discovered--img', gf_s3_info.discovered_imgs__bucket)
	#-------------------
	
	view_bucket_info([crawler_images__bucket_info_map, main_images__bucket_info_map])

#---------------------------------------------------
def process_bucket(p_name_str, p_bucket):

	i                        = 0
	total_size_kb_int        = 0
	counts_per_day_map       = {}
	counts_per_file_type_map = {
		'thumbnails/':     0, #used??
		'thumbnails':      0,
		'thumbnails_jpeg': 0,
		'thumbnails_png':  0
	}

	for o in p_bucket.objects.all():
		file_name_str       = o.key
		dt                  = o.last_modified
		file_size_bytes_int = o.size
		file_size_kb_int    = file_size_bytes_int/1024

		#IMPORTANT!! - AWS access logs are stored in a /logs dir in some of the buckets
		if 'logs' in file_name_str:
			continue

		#-------------------
		#COUNTS_PER_EXTENSION

		file_ext_str = file_name_str.split('.')[-1:][0].lower()

		if file_ext_str == 'jpg':
			file_ext_str = 'jpeg'

		#FILE_EXTENSION_COUNTS
		if counts_per_file_type_map.has_key(file_ext_str):
			counts_per_file_type_map[file_ext_str] += 1
		else:
			counts_per_file_type_map[file_ext_str] = 1

		#THUMBNAILS
		if 'thumbnails/' in file_name_str:
			counts_per_file_type_map['thumbnails'] += 1

			if file_ext_str == 'jpeg':  counts_per_file_type_map['thumbnails_jpeg'] += 1
			elif file_ext_str == 'png': counts_per_file_type_map['thumbnails_png'] += 1
		#-------------------
		#TOTAL_FILE_COUNTS_PER_DAY

		day_str = '%s-%02d-%02d'%(dt.year,dt.month,dt.day)
		if counts_per_day_map.has_key(day_str):
			counts_per_day_map[day_str] += 1
		else:
			counts_per_day_map[day_str] = 0
		#-------------------

		print '%s - %s - %sKB'%(i,file_name_str,file_size_kb_int)

		total_size_kb_int += file_size_kb_int
		i+=1

	return {
		'name_str':                 p_name_str,
		'total_size_kb_int':        total_size_kb_int,
		'counts_per_day_map':       counts_per_day_map,
		'counts_per_file_type_map': counts_per_file_type_map,
	}

#---------------------------------------------------
def view_bucket_info(p_infos_lst):

	#-------------------------
	for info_map in p_infos_lst:
		print ''
		print ''
		print ''
		print 'bucket_name_str - %s'%(info_map['name_str'])
		print ''

		sorted_lst = info_map['counts_per_day_map'].items()
		sorted_lst.sort()
		
		for k,v in sorted_lst:
			print '%s -- %s'%(k, v) 

		print ''
		print ''
		print ''
		for k,v in info_map['counts_per_file_type_map'].items():
			print '%s%s%s -- %s'%(bg('blue'), k, attr(0), v)

		total_size_mb_int = info_map['total_size_kb_int']/1024
		print 'total size - %s%s%smb'%(bg('green'), total_size_mb_int, attr(0))
		print ''
		print ''
		print ''
	#-------------------------

	df_lst = []
	for bucket_info_map in p_infos_lst:

		day_counts_sorted_lst = sorted([(d_str,c_int) for d_str,c_int in bucket_info_map['counts_per_day_map'].items()])

		#days_lst   = []
		#counts_lst = []
		for d,c in day_counts_sorted_lst:

			#all_buckets__days_set.add(d)
			#days_lst.append(d)
			#counts_lst.append(c)

			bname_str = bucket_info_map['name_str']
			d = {
				'day':                   d,
				'bucket':                bname_str,
				'%s__count'%(bname_str): c
			}
			df_lst.append(d)


	#df_columns_map = {
	#	"days":        list(all_buckets__days_set),
	#	"total_counts":counts_lst,
	#}
	#df = pd.DataFrame(data=df_columns_map)
	df = pd.DataFrame(df_lst)
	df.set_index("day", drop=True, inplace=True)
	df = df.fillna(0)
	df = df.sort_values(by=['day'])
	#print df

	#-------------------------
	#PLOT

	df.plot.line(figsize=(10, 6), alpha=0.75)

	#fig = plt.figure(figsize=(30,10))

	plt.title("S3 buckets info per day", fontsize=18)
	plt.xlabel("day",                    fontsize=14)
	plt.ylabel('# of images',            fontsize=14)
	plt.xticks(size = 6)
	plt.axes().yaxis.grid() #horizontal-grid

	plt.show()
	#-------------------------