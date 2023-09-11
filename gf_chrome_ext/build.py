


import delegator


print('COPY GF WEB FILES...')

js_files_lst = [
	("./../web/src/gf_apps/gf_tagger/js/gf_tagger_ui.js",   "./js/build"),
    ("./../web/src/gf_apps/gf_tagger/js/gf_tagger_http.js", "./js/build"),
    ("./../web/src/gf_apps/gf_tagger/css/gf_tagger_ui.css", "./css/build")
]

for f_tpl in js_files_lst:
    src_path_str, target_path_str = f_tpl

    c_str = f"cp {src_path_str} {target_path_str}"
    print(c_str)
    r = delegator.run(c_str)

    if not r.out == "": print(r.out)
    if not r.err == "": print(r.err)