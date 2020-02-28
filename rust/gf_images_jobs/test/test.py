













import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir




print("%s/../../build"%(modd_str))


import delegator
print(delegator.run("ls -al %s"%("%s/../../build"%(modd_str))).out)


sys.path.append("%s/../../build"%(modd_str))
import gf_images_jobs_py as gf_images_jobs




print("TEST ------------------------------")
print(dir(gf_images_jobs))






collage__files_lst       = []
collage__output_file_str = "test__collage.jpeg"
for i in range(0, 300):
	collage__files_lst.extend([
		"50b230e8933860a01cd5da61d082887a.jpeg",
		"49a180a9ab8548b69f50e0bb2c96b4d0_thumb_small.jpeg",
		"test_output__contrast_2.jpeg"
	])
gf_images_jobs.create_collage(collage__files_lst,
    collage__output_file_str,
    800,
    80,
    1,
    10)









img_source_file_path_str = "50b230e8933860a01cd5da61d082887a.jpeg"

# NOISE
output_f_str = "test_output__noise.jpeg"
gf_images_jobs.apply_transforms(["noise"],
    img_source_file_path_str,
    output_f_str)
    
# CONTRAST
for i in range(0, 3):

    factor_f     = i * 100.0
    output_f_str = "test_output__contrast_%s.jpeg"%(i)
    gf_images_jobs.apply_transforms(["contrast:%s"%(factor_f)],
        img_source_file_path_str,
        output_f_str)


# SATURATE

saturation_img_source_file_path_str = "49a180a9ab8548b69f50e0bb2c96b4d0_thumb_small.jpeg"
for i in range(0, 3):

    factor_f     = i * 0.5
    output_f_str = "test_output__saturate_%s.jpeg"%(i)
    gf_images_jobs.apply_transforms(["saturate:%s"%(factor_f)],
        saturation_img_source_file_path_str,
        output_f_str)