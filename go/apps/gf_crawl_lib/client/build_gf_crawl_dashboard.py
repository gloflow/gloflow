







import envoy


output_file_str = './bin/gf_crawl_dashboard.js'
files_lst = [
	'./src/dashboard/gf_crawl__img_preview_tooltip.ts',
	'./src/dashboard/gf_crawl_dashboard.ts',
	'./src/dashboard/gf_crawl_events.ts',
	'./src/dashboard/gf_crawl_images_browser.ts',
]

print 'files_lst - %s'%(files_lst)

print 'RUNNING COMPILE...'
r = envoy.run('tsc --out %s %s'%(output_file_str,
								' '.join(files_lst)))

print r.std_out
print r.std_err