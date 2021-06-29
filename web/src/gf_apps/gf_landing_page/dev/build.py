


import delegator


output_file_str = './bin/gf_landing_page.js'
files_lst = [
	'./../ts/gf_calc.ts',
	'./../ts/gf_email_registration.ts',
	'./../ts/gf_landing_page.ts',
	'./../ts/gf_images.ts',
	'./../ts/procedural_art/gf_procedural_art.ts',
]

print('files_lst - %s'%(files_lst))

print('RUNNING COMPILE...')
r = delegator.run(f"tsc --module system --out {output_file_str} {' '.join(files_lst)}")

print(r.out)
print(r.err)