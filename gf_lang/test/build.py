
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


c_str = f"tsc --module system --target es2016 --outFile {output_file_str} {' '.join(ts_files_lst)}"
r = delegator.run(c_str)

print(r.out)
print(r.err)


print("--------------------------------------------------")
print("COPYING FILES...")

copy_lst = [
    (f"{modd_str}/build/gf_lang.js",  f"{modd_str}/../build/gf_lang.js"),
    (f"{modd_str}/gf_lang_test.html", f"{modd_str}/../build/gf_lang_test.html"),
    (f"{modd_str}/../css/gf_ide.css", f"{modd_str}/../build/gf_ide.css"),

    # (f"{modd_str}/../go/build/gf_lang_web.wasm", f"{modd_str}/../build/gf_lang_web.wasm"),
]

for f in copy_lst:
    print(f[0], "  ", f[1])

    r = delegator.run(f"cp {f[0]} {f[1]}")
    if not r.out == "": print(r.out)
    if not r.err == "": print(r.err)

    if r.return_code != 0:
        print("ERROR")
        sys.exit(1)