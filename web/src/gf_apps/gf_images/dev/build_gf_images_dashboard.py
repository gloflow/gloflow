


import envoy


output_file_str = './bin/gf_images_dashboard.js'
files_lst = [
	'./../ts/gf_images_dashboard/gf_images_dashboard.ts',
	'./../ts/gf_images_dashboard/gf_images_stats.ts',
]

print 'files_lst - %s'%(files_lst)

print 'RUNNING COMPILE...'
c = 'tsc --out %s %s'%(output_file_str, ' '.join(files_lst))
print c

r = envoy.run(c)

print r.std_out
print r.std_err