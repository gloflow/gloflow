


import envoy


output_file_str = './bin/gf_landing_page.js'
files_lst = [
	'./src/gf_calc.ts',
	'./src/gf_email_registration.ts',
	'./src/gf_landing_page.ts',
	'./src/gf_procedural_art.ts',
	'./src/gf_images.ts'
]

print 'files_lst - %s'%(files_lst)

print 'RUNNING COMPILE...'
r = envoy.run('tsc --out %s %s'%(output_file_str,
								' '.join(files_lst)))

print r.std_out
print r.std_err