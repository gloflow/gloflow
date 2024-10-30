




from gf_core import gf_core_sql_db
from gf_apps.gf_images.gf_images_client import gf_images_client



print("gloflow...")


#-------------------------
# EXPORT_API
# simplified unified GF Py API with flattened namespace
db_init_client  = gf_core_sql_db.init_db_client
db_table_exists = gf_core_sql_db.table_exists
add_image       = gf_images_client.add_image

#-------------------------
def run():
    print("gloflow.run()")
    True

