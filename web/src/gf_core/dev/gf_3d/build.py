


import delegator


output_file_str = './bin/gf_3d.js'
files_lst = [
	'gf_3d_test.ts',
	'./../../ts/gf_3d.ts',
]

print('files_lst - %s'%(files_lst))

print('RUNNING COMPILE...')
r = delegator.run(f"tsc --module system --out {output_file_str} {' '.join(files_lst)}")

print(r.out)
print(r.err)