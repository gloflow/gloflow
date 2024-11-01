




from gf_core import gf_core_sql_db
from gf_apps.gf_images.gf_images_client import gf_images_client
from gf_observe import gf_extern_load


print("gloflow...")
version = "0.1.17"

#-------------------------
# EXPORT_API
# simplified unified GF Py API with flattened namespace

# DB
db_init_client  = gf_core_sql_db.init_db_client
db_table_exists = gf_core_sql_db.table_exists

# IMAGES
add_image = gf_images_client.add_image

# OBSERVE
observe_init     = gf_extern_load.init
observe_ext_load = gf_extern_load.observe

#-------------------------
def run():
    print("gloflow.run()")
    True

