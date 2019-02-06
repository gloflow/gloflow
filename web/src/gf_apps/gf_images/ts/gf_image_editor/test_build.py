







import envoy


print '---------------------'
print ''
print 'RUN FOR TESTING >>>> -- "python -m SimpleHTTPServer 2222"'
print ''
print '---------------------'


output_file_str = './test/gf_image_editor__test_build.js'
files_lst = [
	'gf_image_editor.ts',
]

print 'files_lst - %s'%(files_lst)

print 'RUNNING COMPILE...'
r = envoy.run('tsc --out %s %s'%(output_file_str,
								' '.join(files_lst)))

print r.std_out
print r.std_err