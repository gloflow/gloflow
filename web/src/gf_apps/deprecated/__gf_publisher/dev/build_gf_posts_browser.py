


import envoy


output_file_str = './bin/gf_posts_browser.js'
files_lst = [
	'./src/gf_posts_browser/gf_posts_browser.ts',
	'./src/gf_posts_browser/gf_posts_browser_view.ts',
	'./src/gf_posts_browser/gf_posts_browser_client.ts',
	'../../gf_tagger/client/src/gf_tagger_client/gf_tagger_client.ts',
	'../../gf_tagger/client/src/gf_tagger_client/gf_tagger_input_ui.ts',
	'../../gf_tagger/client/src/gf_tagger_client/gf_tagger_notes_ui.ts',
	'../../../gf_core/client/src/gf_sys_panel.ts'
]

print('files_lst - %s'%(files_lst))

#FIX!! - dont use the "--t es6" flag, since minifierjs doesnt support its features yet
print('RUNNING COMPILE...')
c = 'tsc --t es6 --out %s %s'%(output_file_str, ' '.join(files_lst))
print(c)
r = envoy.run(c)

print(r.std_out)
print(r.std_err)