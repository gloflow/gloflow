


import delegator

#----------------------------------------
# TS
print('RUNNING TS COMPILE...')

output_file_str = 'gf_viz_group.js'
ts_files_lst = [
	'./gf_viz_group_test.ts',
    './../ts/gf_viz_group_paged.ts'
]

for f in ts_files_lst:
    print(f)


r = delegator.run(f"tsc --module system --target es6 --out {output_file_str} {' '.join(ts_files_lst)}")

print(r.out)
print(r.err)