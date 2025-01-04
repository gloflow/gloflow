


import envoy


output_file_str = './bin/gf_post.js'
files_lst = [
	'./src/gf_post/gf_post.ts',
	'./src/gf_post/gf_post_image_view.ts',
	'../../gf_tagger/client/src/gf_tagger_client/gf_tagger_client.ts',
	'../../gf_tagger/client/src/gf_tagger_client/gf_tagger_input_ui.ts',
	'../../../gf_core/client/src/gf_sys_panel.ts'
]

print('files_lst - %s'%(files_lst))

print('RUNNING COMPILE...')
c = 'tsc --out %s %s'%(output_file_str, ' '.join(files_lst))
print(c)

r = envoy.run(c)

print(r.std_out)
print(r.std_err)