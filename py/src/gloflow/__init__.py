from gf_core import gf_core_id, gf_core_sql_db
from gf_apps.gf_images.gf_images_client import gf_images_client
from gf_observe import gf_extern_load
from gf_ml import gf_llm_core

print("gloflow...")
version = "0.1.18"

#-------------------------
# EXPORT_API
# simplified unified GF Py API with flattened namespace

# DB
class db():
    init         = gf_core_sql_db.init_db_client
    table_exists = gf_core_sql_db.table_exists

# IMAGES
class images():
    add = gf_images_client.add_image

# CORE
class core():
    create_id = gf_core_id.create

# OBSERVE
class observe():
    init           = gf_extern_load.init
    ext_load       = gf_extern_load.observe
    get_cached     = gf_extern_load.get_cached
    get_cached_group = gf_extern_load.get_cached_group
    relate           = gf_extern_load.relate

# ML
class ml():
    run_llm = gf_llm_core.run_llm