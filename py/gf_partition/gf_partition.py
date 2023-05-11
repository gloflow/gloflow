# GloFlow application and media management/publishing platform
# Copyright (C) 2023 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

sys.path.append(f"{modd_str}/../gf_core")
import gf_core_sql_db

#---------------------------------------------------------------------------------
# DECORATORS
#---------------------------------------------------------------------------------
def gf_materialize(p_partition_map,
    p_partition_sql_id_int,
    p_dagster_run_id,
    p_db_client):
    def decorator(p_func):
        def wrapper(*args, **kwargs):
            func_name_str = p_func.__name__
            print(f"gf_materialize decorator ---> {p_func.__name__}")
            
            # CREATE
            materialization_sql_id_int = db_create_materialization(p_partition_map,
                p_partition_sql_id_int,
                p_dagster_run_id,
                p_db_client)
            
            # FUNC
            try:
                result = p_func(*args, **kwargs)

            except Exception as e:
                
                db_update_materialization("failed", )
                
                # rethrow the exception further up the chain
                raise
                    
            # COMPLETE
            db_update_materialization("completed", materialization_sql_id_int, p_db_client)

            return result
        return wrapper
    return decorator

#---------------------------------------------------------------------------------
def db_create_materialization(p_partition_map,
    p_partition_sql_id_int,
    p_dagster_run_id_str,
    p_db_client):
    assert isinstance(p_partition_map, dict)

    print("create materialization...")

    table_name_str = "gf_data_partitions_materilize"
    status_str     = "started"

    query_str = f"""
        INSERT INTO {table_name_str} (
            set_name,
            group_name,
            status,
            dagster_run_id,
            partition_id) 
        VALUES (%s, %s, %s, %s, %s)
        RETURNING id;"""

    values_tpl = (p_partition_map["set_name_str"],
        p_partition_map["group_name_str"],
        status_str,
        p_dagster_run_id_str,
        p_partition_sql_id_int)

    cur = p_db_client.cursor()
    cur.execute(query_str, values_tpl)

    new_sql_id_int = cur.fetchone()[0]


    print("DDDDDDDDDDDDDDDDDDDD")
    print(new_sql_id_int)



    p_db_client.commit()
    cur.close()

    return new_sql_id_int

#---------------------------------------------------------------------------------
def db_update_materialization(p_status_str,
    p_materialization_id_int,
    p_db_client):
    assert p_status_str == "failed" or p_status_str == "completed"

    table_name_str = "gf_data_partitions_materilize"
    status_str     = "completed"

    query_str = f"""
        UPDATE {table_name_str}
        SET status = '{p_status_str}'
        WHERE id={p_materialization_id_int};
    """

    cur = p_db_client.cursor()
    cur.execute(query_str)
    p_db_client.commit()
    cur.close()

#---------------------------------------------------------------------------------
# MAIN
#---------------------------------------------------------------------------------
def db_create(p_partition_map,
    p_db_client):
    assert isinstance(p_partition_map, dict)

    print("create partition...")

    table_name_str = "gf_data_partitions"

    query_str = f"""
        INSERT INTO {table_name_str} (
            set_name,
            group_name,
            dagster_asset_id,
            partition_i,
            partition_size,
            materialized_actively,
            materialized) 
        VALUES (%s, %s, %s, %s, %s, %s, %s)
        RETURNING id;"""

    values_tpl = (p_partition_map["set_name_str"],
        p_partition_map["group_name_str"],
        p_partition_map["dagster_asset_id_str"],
        p_partition_map["partition_i_int"],
        p_partition_map["partition_size_int"],
        False, # materialized_actively
        True)  # materialized

    cur = p_db_client.cursor()
    cur.execute(query_str, values_tpl)

    new_partition_sql_id_int = cur.fetchone()[0]

    p_db_client.commit()
    cur.close()

    return new_partition_sql_id_int
    
#---------------------------------------------------------------------------------
def db_check_exists(p_set_name_str,
    p_group_name_str,
    p_partition_i_int,
    p_partition_size_int,
    p_db_client):
    
    table_name_str = "gf_data_partitions"
    partition_str = "%s:%s:%s:%s"%(p_set_name_str, p_group_name_str, p_partition_i_int, p_partition_size_int)

    try:
        result = db_get_one(p_set_name_str,
            p_group_name_str,
            p_partition_i_int,
            p_partition_size_int,
            p_db_client)

    except Exception as e:
        if isinstance(e, errors.UndefinedTable):
            # whole table doesnt exist, not just the partition, so report that
            return False
        else:
            # rethrow the exception further up the chain
            raise
    
    if result is not None:
        
        materialized_bool = result[0]

        print(f"partition exists - {partition_str}")
        return True, materialized_bool
    else:
        print(f"partition doesnt exists - {partition_str}")
        return False, False

    cur.close()
    p_db_client.close()

#---------------------------------------------------------------------------------
def db_get_one(p_set_name_str,
    p_group_name_str,
    p_partition_i_int,
    p_partition_size_int,
    p_db_client):

    table_name_str = "gf_data_partitions"
    cur = p_db_client.cursor()
    query_str = f"""
        SELECT *
        FROM {table_name_str}
        WHERE set_name = %s AND group_name = %s AND partition_i = %s AND partition_size = %s
    """

    # QUERY
    cur.execute(query_str,
        (p_set_name_str, p_group_name_str, p_partition_i_int, p_partition_size_int))

    result = cur.fetchone()
    return result

#---------------------------------------------------------------------------------
def db_init(p_db_client):
    
    cur = p_db_client.cursor()

    table_name_str = "gf_data_partitions"
    table_partition_materialize_name_str = "gf_data_partitions_materilize"


    #---------------------------------------------------------------------------------
    def create_table():

        cur.execute(f"DROP TABLE {table_name_str} CASCADE;")
        cur.execute(f"DROP TABLE {table_partition_materialize_name_str} CASCADE;")

        # GF_DATA_PARTITIONS
        if not gf_core_sql_db.table_exists(table_name_str, cur):

            # dagster_asset_id - rereference to the appropriate Dagster Asset that
            #                    materialized/computed this partition.
            # actively_processed - flag indicating if anyone is activelly processing this 
            # materialized_bool - has this partition been ever materialized?
            create_table_sql_str = f"""
                CREATE TABLE {table_name_str} (
                    id               BIGSERIAL PRIMARY KEY,
                    created_at       TIMESTAMP DEFAULT NOW(),
                    set_name         VARCHAR(255) NOT NULL,
                    group_name       VARCHAR(255) NOT NULL,
                    dagster_asset_id VARCHAR(255) NOT NULL,

                    partition_i    INT NOT NULL,
                    partition_size INT NOT NULL,

                    materialized_started_at TIMESTAMP,
                    materialized_actively   BOOLEAN NOT NULL,
                    materialized            BOOLEAN NOT NULL,

                    UNIQUE(set_name, group_name, partition_i, partition_size)
                );"""
            cur.execute(create_table_sql_str)
            p_db_client.commit()

        # GF_DATA_PARTITIONS_MATERIALIZE
        if not gf_core_sql_db.table_exists(table_partition_materialize_name_str, cur):

            # status         - "started" | "failed" | "completed"
            # dagster_run_id - ID of the Dagster Run within which this materialization is executing
            create_table_sql_str = f"""
                CREATE TABLE {table_partition_materialize_name_str} (
                    id             BIGSERIAL PRIMARY KEY,
                    created_at     TIMESTAMP DEFAULT NOW(),
                    set_name       VARCHAR(255) NOT NULL,
                    group_name     VARCHAR(255) NOT NULL,
                    status         VARCHAR(255) NOT NULL,
                    dagster_run_id VARCHAR(255) NOT NULL,
                    partition_id   BIGINT REFERENCES {table_name_str}(id)
                );"""
            cur.execute(create_table_sql_str)
            p_db_client.commit()
    
    #---------------------------------------------------------------------------------
    create_table()