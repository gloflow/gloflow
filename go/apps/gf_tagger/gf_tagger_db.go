package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"apps/gf_publisher_lib"
)
//---------------------------------------------------
func db__get_objects_with_tag_count(p_tag_str string,
						p_object_type_str string,
						p_mongodb_coll    *mgo.Collection,
						p_log_fun         func(string,string)) (int,error) {
	p_log_fun("FUN_ENTER","gf_tagger_db.db__get_objects_with_tag_count()")

	switch p_object_type_str {
		case "post":
			count_int,err := p_mongodb_coll.Find(bson.M{
												"t"       :"post",
												"tags_lst":bson.M{"$in":[]string{p_tag_str,}},
											}).
											Count()
			if err != nil {
				return 0,err
			}

			return count_int,nil
	}

	return 0,nil
}
//---------------------------------------------------
//POSTS
//---------------------------------------------------
func db__get_post_notes(p_post_title_str *string,
			p_mongodb_coll *mgo.Collection,
			p_log_fun      func(string,string)) ([]*Note,error) { //([]*gf_publisher_lib.Post_snippet,error) {
	p_log_fun("FUN_ENTER","gf_tagger_db.db__get_post_notes()")

	post,err := gf_publisher_lib.DB__get_post(p_post_title_str,
										p_mongodb_coll,
										p_log_fun)
	if err != nil {
		return nil,err
	}

	post_notes_lst := post.Notes_lst
	notes_lst      := []*Note{}
	for _,s := range post_notes_lst {

		note := &Note{
			User_id_str          :s.User_id_str,
			Body_str             :s.Body_str,
			Target_obj_id_str    :post.Title_str,
			Target_obj_type_str  :"post",
			Creation_datetime_str:s.Creation_datetime_str,
		}
		notes_lst = append(notes_lst,note)
	}
	p_log_fun("INFO","got # notes - "+fmt.Sprint(len(notes_lst)))
	return notes_lst,nil //post_notes_lst,nil
}
//---------------------------------------------------
func db__add_post_note(p_note *Note,
			p_post_title_str *string,
			p_mongodb_coll   *mgo.Collection,
			p_log_fun        func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_tagger_db.db__add_post_note()")

	//--------------------
	post_note := &gf_publisher_lib.Post_note{
		User_id_str          :p_note.User_id_str,
		Body_str             :p_note.Body_str,
		Creation_datetime_str:p_note.Creation_datetime_str,
	}
	//--------------------
	
	err := p_mongodb_coll.Update(bson.M{"t":"post","title_str":*p_post_title_str,},
							bson.M{"$push":bson.M{"notes_lst":post_note},})
	if err != nil {
		return err
	}
	return nil
}
//---------------------------------------------------
func db__get_posts_with_tag(p_tag_str string,
					p_page_index_int int,
					p_page_size_int  int,
					p_mongodb_coll   *mgo.Collection,
					p_log_fun        func(string,string)) ([]*gf_publisher_lib.Post,error) {
	p_log_fun("FUN_ENTER","gf_tagger_db.db__get_posts_with_tag()")
	p_log_fun("INFO"     ,"p_tag_str - "+p_tag_str)

	//FIX!! - potentially DOESNT SCALE. if there is a huge number of posts
	//        with a tag, toList() will accumulate a large collection in memory. 
	//        instead use a Stream-oriented way where results are streamed lazily
	//        in some fashion
		
	var posts_lst []*gf_publisher_lib.Post
	err := p_mongodb_coll.Find(bson.M{
								"t"       :"post",
								"tags_lst":bson.M{"$in":[]string{p_tag_str,}},
							}).
							Sort("-creation_datetime"). //descending:true
							Skip(p_page_index_int).
							Limit(p_page_size_int).
							All(&posts_lst)

	if err != nil {
		return nil,err
	}

	return posts_lst,nil
}
//---------------------------------------------------
func db__add_tags_to_post(p_post_title_str *string,
					p_tags_lst     []string,
					p_mongodb_coll *mgo.Collection,
					p_log_fun      func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_tagger_db.db__add_tags_to_post()")

	err := p_mongodb_coll.Update(bson.M{"t":"post","title_str":*p_post_title_str,},
							bson.M{"$push":bson.M{"tags_lst":p_tags_lst},})
	if err != nil {
		return err
	}

	return nil
}
//---------------------------------------------------
//IMAGES
//---------------------------------------------------
func db__add_tags_to_image(p_image_id_str *string,
						p_tags_lst     []string,
						p_mongodb_coll *mgo.Collection,
						p_log_fun      func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_tagger_db.db__add_tags_to_image()")

	err := p_mongodb_coll.Update(bson.M{"t":"img","id_str":*p_image_id_str,},
							bson.M{"$push":bson.M{"tags_lst":p_tags_lst},})
	if err != nil {
		return err
	}

	return nil
}