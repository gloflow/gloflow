
import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import delegator

#----------------------------------------
# TS
print('RUNNING TS COMPILE...')

output_file_str = f'{modd_str}/build/gf_lang.js'
ts_files_lst = [
	f'{modd_str}/../ts/gf_lang.ts',
    f'{modd_str}/gf_lang_test.ts'
]

for f in ts_files_lst:
    print(f)


c_str = f"tsc --module system --target es2016 --out {output_file_str} {' '.join(ts_files_lst)}"
r = delegator.run(c_str)

print(r.out)
print(r.err)



r = delegator.run(f"cp {modd_str}/build/gf_lang.js {modd_str}/../build/gf_lang.js")
r = delegator.run(f"cp {modd_str}/../go/build/gf_lang_web.wasm {modd_str}/../build/gf_lang_web.wasm")