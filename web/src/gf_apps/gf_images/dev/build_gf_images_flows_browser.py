


import envoy


output_file_str = './bin/gf_images_flows_browser.js'
files_lst = [
	'./src/flows_browser/gf_images_flows_browser.ts',
	'./src/flows_browser/gf_images_viewer.ts',
	'../../../gf_core/client/src/gf_sys_panel.ts'
]

print 'files_lst - %s'%(files_lst)

print 'RUNNING COMPILE...'
c = 'tsc --out %s %s'%(output_file_str, ' '.join(files_lst))
print c

r = envoy.run(c)

print r.std_out
print r.std_err