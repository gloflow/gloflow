


import delegator

#----------------------------------------
# TS
print('RUNNING TS COMPILE...')

output_file_str = './bin/gf_image_colors_test.js'
ts_files_lst = [
    './gf_image_colors_test.ts'
]

for f in ts_files_lst:
    print(f)


r = delegator.run(f"tsc --module system --target es2017 --out {output_file_str} {' '.join(ts_files_lst)}")

print(r.out)
print(r.err)