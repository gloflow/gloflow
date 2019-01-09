


import envoy


output_file_str = './bin/gf_domains_browser.js'
files_lst = [
	'./src/domains_browser/gf_domains_browser.ts',
	'./src/domains_browser/gf_domain.ts',
	'./src/domains_browser/gf_domains_conn.ts',
	'./src/domains_browser/gf_domains_infos.ts',
	'./src/domains_browser/gf_domains_search.ts',
	'../../../gf_core/client/src/gf_color.ts',
]

print 'files_lst - %s'%(files_lst)

print 'RUNNING COMPILE...'
c = 'tsc --out %s %s'%(output_file_str,
					' '.join(files_lst))

print c
r = envoy.run(c)

print r.std_out
print r.std_err