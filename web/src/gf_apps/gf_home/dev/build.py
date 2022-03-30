


import delegator


output_file_str = './bin/gf_home_main.js'
files_lst = [
	'./../ts/gf_home_main.ts',
]

print('files_lst - %s'%(files_lst))

print('RUNNING COMPILE...')
r = delegator.run(f"tsc --module system --target es2017 --out {output_file_str} {' '.join(files_lst)}")

print(r.out)
print(r.err)