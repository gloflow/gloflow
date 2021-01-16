







import delegator


print("---------------------")
print("")
print("RUN FOR TESTING >>>> -- 'python -m SimpleHTTPServer 2222'")
print("")
print("---------------------")


output_file_str = "./test/build/gf_image_editor__test_build.js"
files_lst = [
	"gf_image_editor.ts",
]

print("files_lst - %s"%(files_lst))

print("RUNNING COMPILE...")

# "--module system" - needed with the "--out" option
r = delegator.run(f'tsc --module system --out {output_file_str} {" ".join(files_lst)}')

print(r.out)
print(r.err)