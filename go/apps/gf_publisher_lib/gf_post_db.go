package gf_publisher_lib

import (
	"time"
	"fmt"
	"math/rand"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)
//---------------------------------------------------
func DB__get_post(p_post_title_str *string,
	p_mongodb_coll *mgo.Collection,
	p_log_fun      func(string,string)) (*Post,error) {
	p_log_fun("FUN_ENTER","gf_post_db.DB__get_post()")

	var post Post
	//final mongo.Cursor c = posts_coll.find(mongo.where.eq("title_str",p_post_title_str));
	err := p_mongodb_coll.Find(bson.M{"t":"post","title_str":*p_post_title_str}).One(&post)
	if err != nil {
		return nil,err
	}

	return &post,nil

		/*if (result_map != null) {
			
			post_info_map = gf_post_serialize.deserialize(result_map,
																	p_log_fun);
			final gf_post.Post_ADT post_adt = gf_post.create(post_info_map,
															 p_log_fun);
			return post_adt;
		}
		else {
			p_log_fun("INFO","POST WITH TITLE ($p_post_title_str) NOT FOUND");
			return null;
		}
	}
	//---------------------
	//ERROR_HANDLING
	catch(p_exc,
		  p_stack_trace) {
		final Map error_map = {
			"msg_str"    :"failed to get post from DB with title [$p_post_title_str]",
			"error"      :p_exc,
			"stack_trace":p_stack_trace
		};
		p_log_fun("ERROR",error_map);
		throw error_map;
	};
	//---------------------*/
}
//---------------------------------------------------
func DB__get_posts_page(p_cursor_start_position_int int, //0
	p_elements_num_int int, //50
	p_mongodb_coll     *mgo.Collection,
	p_log_fun          func(string,string)) ([]*Post,error) {
	p_log_fun("FUN_ENTER","gf_post_db.DB__get_posts_page()")

	posts_lst := []*Post{}
	//descending - true - sort the latest items first
	err := p_mongodb_coll.Find(bson.M{"t":"post"}).
		Sort("-creation_datetime_str"). //descending:true
		Skip(p_cursor_start_position_int).
		Limit(p_elements_num_int).
		All(&posts_lst)
	if err != nil {
		return nil,err
	}

	return posts_lst,err
}
//---------------------------------------------------
func DB__create_post(p_post *Post,
	p_mongodb_coll *mgo.Collection,
	p_log_fun      func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_post_db.DB__create_post()")

	err := p_mongodb_coll.Insert(p_post) //writeConcern: mongo.WriteConcern.ACKNOWLEDGED);
	if err != nil {
		return err
	}

	return nil
}
//---------------------------------------------------
func DB__update_post(p_post *Post, 
	p_mongodb_coll *mgo.Collection,
	p_log_fun      func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_post_db.DB__update_post()")

	err := p_mongodb_coll.Update(bson.M{
			"t"        :"post",
			"title_str":p_post.Title_str,
		},
		p_post)
	if err != nil {
		return err
	}

	return nil
}
/*//---------------------------------------------------
func DB__delete_post(p_post_title_str *string,
			p_mongodb_coll *mgo.Collection,
			p_log_fun      func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_post_db.DB__delete_post()")

	posts_coll.remove({"title_str":p_post_title_str},
							mongo.WriteConcern.ACKNOWLEDGED);
	return;
}*/
//---------------------------------------------------
func DB__get_random_posts_range(p_posts_num_to_get_int int, //5
	p_max_random_cursor_position_int int, //500
	p_mongodb_coll                   *mgo.Collection,
	p_log_fun                        func(string,string)) ([]*Post,error) {
	p_log_fun("FUN_ENTER","gf_post_db.DB__get_random_posts_range()")

	rand.Seed(time.Now().Unix())
	random_cursor_position_int := rand.Intn(p_max_random_cursor_position_int) //new Random().nextInt(p_max_random_cursor_position_int)
	p_log_fun("INFO","random_cursor_position_int - "+fmt.Sprint(random_cursor_position_int))

	posts_in_random_range_lst,err := DB__get_posts_from_offset(random_cursor_position_int,
		p_posts_num_to_get_int,
		p_mongodb_coll,
		p_log_fun)
	if err != nil {
		return nil,err
	}

	return posts_in_random_range_lst,nil
}
//---------------------------------------------------
func DB__get_posts_from_offset(p_cursor_position_int int,
	p_posts_num_to_get_int int,
	p_mongodb_coll         *mgo.Collection,
	p_log_fun              func(string,string)) ([]*Post,error) {
	p_log_fun("FUN_ENTER","gf_post_db.DB__get_posts_from_offset()")

	//----------------
	//IMPORTANT!! - because mongo"s skip() scans the collections number of docs
	//				to skip before it returns the results, p_posts_num_to_get_int should
	//				not be large, otherwise the performance will be very bad
	//assert(p_cursor_position_int < 500);
	//----------------

	var posts_lst []*Post
	err := p_mongodb_coll.Find(bson.M{"t":"post"}).
		Skip(p_cursor_position_int).
		Limit(p_posts_num_to_get_int).
		All(&posts_lst)
	if err != nil {
		return nil,err
	}

	return posts_lst,nil
}
//---------------------------------------------------
func DB__check_post_exists(p_post_title_str *string,
	p_mongodb_coll *mgo.Collection,
	p_log_fun      func(string,string)) (bool,error) {
	p_log_fun("FUN_ENTER","gf_post_db.DB__check_post_exists()")
	
	count_int,err := p_mongodb_coll.Find(bson.M{"t":"post","title_str":*p_post_title_str,}).Count()
	if err != nil {
		return false,err
	}
	if count_int > 0 {
		return true,nil
	} else {
		return false,nil
	}
}