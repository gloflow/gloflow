


import delegator

#----------------------------------------
# TS
print('RUNNING TS COMPILE...')

output_file_str = './build/gf_lang.js'
ts_files_lst = [
	'./../ts/gf_lang.ts',
    './gf_lang_test.ts'
]

for f in ts_files_lst:
    print(f)


r = delegator.run(f"tsc --module system --target es2016 --out {output_file_str} {' '.join(ts_files_lst)}")

print(r.out)
print(r.err)



r = delegator.run(f"cp ./build/gf_lang.js ./../build/gf_lang.js")
r = delegator.run(f"cp ./../go/build/gf_lang_web.wasm ./../build/gf_lang_web.wasm")