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

import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import argparse
import pymongo

sys.path.append('%s/../stats'%(cwd_str))
import crawler_page_imgs__counts_by_day
import crawler_page_outgoing_links__counts_by_day
import crawler_page_outgoing_links__null_breakdown
import crawler_page_outgoing_links__per_crawler
import crawler_url_fetches__counts_by_day
#-------------------------------------------------------------
def main(p_log_fun):

	args_map = parse_args()

	plots_dir_str    = args_map['plots_dir']
	mongodb_host_str = args_map['mongodb_host']
	mongo_client     = get_mongodb_client(mongodb_host_str,p_log_fun)


	stats__config_map = {

		'crawler_page_imgs__counts_by_day':{
			'plot_path_str':'%s/crawler_page_imgs__counts_by_day.png'%(plots_dir_str)
		},

		'crawler_page_outgoing_links__counts_by_day':{
			'plot_path_str':'%s/crawler_page_outgoing_links__counts_by_day.png'%(plots_dir_str)
		},

		'crawler_page_outgoing_links__per_crawler':{
			'plot_path_str':'%s/crawler_page_outgoing_links__per_crawler.png'%(plots_dir_str)
		},

		'crawler_page_outgoing_links__null_breakdown':{
			'plot_path_str':'%s/crawler_page_outgoing_links__null_breakdown.png'%(plots_dir_str)
		},
		'crawler_url_fetches__counts_by_day':{
			'plot_path_str':'%s/crawler_url_fetches__counts_by_day.png'%(plots_dir_str)
		}
	}

	#---------------------
	if args_map['run'] == 'run_batch_sequence':

		run_batch_sequence(stats__config_map, mongo_client, p_log_fun)
	#---------------------
	elif args_map['run'] == 'stat__crawler_page_imgs__counts_by_day':		
		crawler_page_imgs__counts_by_day.run(mongo_client,
			p_log_fun,
			p_output_img_str = stats__config_map['crawler_page_imgs__counts_by_day']['plot_path_str'])
	#---------------------
	elif args_map['run'] == 'stat__crawler_page_outgoing_links__counts_by_day':		
		crawler_page_outgoing_links__counts_by_day.run(mongo_client,
			p_log_fun,
			p_output_img_str = stats__config_map['crawler_page_outgoing_links__counts_by_day']['plot_path_str'])
	#---------------------
	if args_map['run'] == 'stat__crawler_page_outgoing_links__per_crawler':		
		crawler_page_outgoing_links__per_crawler.run(mongo_client,
			p_log_fun,
			p_output_img_str = stats__config_map['crawler_page_outgoing_links__per_crawler']['plot_path_str'])
	#---------------------
	elif args_map['run'] == 'stat__crawler_page_outgoing_links__null_breakdown':
		crawler_page_outgoing_links__null_breakdown.run(mongo_client,
			p_log_fun,
			p_output_img_str = stats__config_map['crawler_page_outgoing_links__null_breakdown']['plot_path_str'])
	#---------------------
	elif args_map['run'] == 'stat__crawler_url_fetches__counts_by_day':
		crawler_url_fetches__counts_by_day.run(mongo_client,
			p_log_fun,
			p_output_img_str = stats__config_map['crawler_url_fetches__counts_by_day']['plot_path_str'])
	#---------------------
#----------------------------------------------
def run_batch_sequence(p_stats__config_map,
	p_mongo_client,
	p_log_fun):
	p_log_fun('FUN_ENTER','cli__gf_crawl__stats.run_batch_sequence()')


	print ''
	print '         ------ RUN_BATCH >>>>'
	print ''

	crawler_page_outgoing_links__counts_by_day.run(p_mongo_client,
		p_log_fun,
		p_output_img_str = p_stats__config_map['crawler_page_outgoing_links__counts_by_day']['plot_path_str'])

	crawler_page_outgoing_links__per_crawler.run(p_mongo_client,
		p_log_fun,
		p_output_img_str = p_stats__config_map['crawler_page_outgoing_links__per_crawler']['plot_path_str'])

	crawler_page_outgoing_links__null_breakdown.run(p_mongo_client,
		p_log_fun,
		p_output_img_str = p_stats__config_map['crawler_page_outgoing_links__null_breakdown']['plot_path_str'])
#----------------------------------------------
#ADD!! - figure out a smarter way to pick the right hostport from p_host_port_lst,
#        instead of just picking the first element

def get_mongodb_client(p_host_str, p_log_fun):
	p_log_fun('FUN_ENTER','cli__gf_crawl__stats.get_mongodb_client()')

	mongo_client = pymongo.MongoClient(p_host_str,27017)
	return mongo_client
#-------------------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

	#---------------------------------
	arg_parser.add_argument('-run', 
		action  = "store",
		default = None,
		help    = '''
- run_batch_sequence                                - process all stats in a sequence
- stat__crawler_page_imgs__counts_by_day            - various counts of new discovered page images per day 
- stat__crawler_page_outgoing_links__counts_by_day  - various counts of new discovered page links per day 
- stat__crawler_page_outgoing_links__null_breakdown - stats on discovered page links that dont have a crawler_name
- stat__crawler_page_outgoing_links__per_crawler    - number of discovered page links per crawler_name (histogram)
- stat__crawler_url_fetches__counts_by_day          - various counts of url_fetch's per day 

					''')
	#---------------------------------
	arg_parser.add_argument('-mongodb_host', 
		action  = "store",
		default = '127.0.0.1',
		help    = '''
host of the Mongodb server
					''')
	#---------------------------------
	arg_parser.add_argument('-plots_dir', 
		action  = "store",
		default = '../plots/',
		help    = '''
dir in which to place generated plots
					''')
	#---------------------------------
	arg_parser.add_argument('-env_var_args', 
		action  = "store",
		default = 'false',
		help    = '''
if arguments should be read from the ENV variables as well the CLI 
					''')
	#---------------------------------

	passed_in_args_lst = sys.argv[1:]
	args_namespace     = arg_parser.parse_args(passed_in_args_lst)


	if args_namespace.env_var_args == 'true':

		mongodb_host_str = os.environ['GF_MONGODB_HOST']
		plots_dir_str    = os.environ['GF_PLOTS_DIR']

		print 'mongodb_host_str - %s'%(mongodb_host_str)
		print 'plots_dir_str    - %s'%(plots_dir_str)

		return {
			'run':         args_namespace.run,
			'mongodb_host':mongodb_host_str,
			'plots_dir':   plots_dir_str
		}
	else:
		return {
			'run':         args_namespace.run,
			'mongodb_host':args_namespace.mongodb_host,
			'plots_dir':   args_namespace.plots_dir
		}
#-------------------------------------------------------------
if __name__ == '__main__':
	def log_fun(g,m):print '%s:%s'%(g,m)
	main(log_fun)