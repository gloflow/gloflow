




# image storage:
#   - this represents storage of image that is expected to be durrable in some way
#   - does not include temporary storage of images on the local FS (workspace)

# ------------------------------------------------------
# functions that do persistent image storage:

# GF_GIF_LIB
gf_apps/gf_images_lib/gf_gif_lib/gf_gif.go
    - Process_and_upload()
    - gif__s3_upload_preview_frames()
    - Gif__frames__save_to_fs()

# GF_IMAGES_CORE
gf_apps/gf_images_lib/gf_images_core/gf_images_s3.go
    - S3__store_gf_image()        - uploads both main GF transformed image
    - S3__store_gf_image_thumbs() - uploads thumbs

# GF_IMAGES_JOBS_CORE
gf_apps/gf_images_lib/gf_images_jobs_core/gf_jobs_pipeline.go
    - job__pipeline__process_image_uploaded()
    - job__pipeline__process_image_extern()

# ------------------------------------------------------
# functions that use local temporary workspace FS storage:

# GF_IMAGE_EDITOR
gf_apps/gf_images_lib/gf_image_editor/gf_image_editor.go
    - save_edited_image()