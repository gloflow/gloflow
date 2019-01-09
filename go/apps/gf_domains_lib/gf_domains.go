package gf_domains_lib

import (
	"fmt"
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/fatih/color"
	"gf_core"
)
//--------------------------------------------------
type Domain struct {
	Id                    bson.ObjectId `bson:"_id,omitempty"`
	Id_str                string        `bson:"id_str"`
	T_str                 string        `bson:"t"` //"domain"

	Name_str      string        `bson:"name_str"`
	Count_int     int           `bson:"count_int"`
	Domain_posts  Domain_Posts  `bson:"posts_domain"`
	Domain_images Domain_Images `bson:"images_domain"`
}
//--------------------------------------------------
func Init_domains_aggregation(p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains.Init_domains_aggregation()")

	go func() {
		for ;; {

			//--------------------
			//IMPORTANT!! - RUN AGGREGATION EVERY Xs (since this is a demanding aggregation)
			//              this is run first, in the loop, so that initialy when this is
			//              initialized it doesnt run, and only later when service active 
			//              for a while it will run for its first iteration.
			time_to_sleep := time.Second*time.Duration(60*5) //5min
			time.Sleep(time_to_sleep)
			//--------------------
			
			gf_err := Discover_domains_in_db(p_runtime_sys)
			if gf_err != nil {
				continue
			}
		}
	}()
}
//--------------------------------------------------
func Discover_domains_in_db(p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains.Discover_domains_in_db()")

	//ADD!! - issue the posts/images queries in parallel via their own go-routines
	//---------------
	//POSTS
	posts_domains_lst,gf_err := Get_domains_posts__mongo(p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	//---------------
	//IMAGES
	images_domains_lst,gf_err := Get_domains_images__mongo(p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	//---------------
	//APP_LEVEL_JOIN
	domains_map := accumulate_domains(posts_domains_lst,
								images_domains_lst,
								p_runtime_sys)
	//DB PERSIST
	gf_err = db__persist_domains(domains_map,p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	//--------------------

	return nil
}
//--------------------------------------------------
func accumulate_domains(p_posts_domains_lst []Domain_Posts,
				p_images_domains_lst []Domain_Images,
				p_runtime_sys        *gf_core.Runtime_sys) map[string]Domain {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains.accumulate_domains()")

	domains_map := map[string]Domain{}

	//--------------------------------------------------
	//POSTS DOMAINS
	//IMPORTANT!! - these run first so they just create a Domain struct without checks

	for _,domain_posts := range p_posts_domains_lst {

		domain_name_str := domain_posts.Name_str

		//IMPORTANT!! - no existing domain with this domain_str has been found
		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		id_str               := fmt.Sprintf("domain:%f",creation_unix_time_f)
		new_domain           := Domain{
			Id_str:      id_str,
			T_str:       "domain",
			Name_str:    domain_name_str,
			Count_int:   domain_posts.Count_int,
			Domain_posts:domain_posts,
		}
		domains_map[domain_name_str] = new_domain
	}
	//--------------------------------------------------
	//IMAGES DOMAINS
	for _,images_domain := range p_images_domains_lst {

		domain_name_str := images_domain.Name_str

		if domain,ok := domains_map[domain_name_str]; ok {
			domain.Domain_images = images_domain
			domain.Count_int     = domain.Count_int + images_domain.Count_int
		} else {
			//IMPORTANT!! - no existing domain with this domain_str has been found
			creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
			id_str               := fmt.Sprintf("domain:%f",creation_unix_time_f)
			new_domain := Domain{
				Id_str:       id_str,
				T_str:        "domain",
				Name_str:     domain_name_str,
				Count_int:    images_domain.Count_int,
				Domain_images:images_domain,
			}
			domains_map[domain_name_str] = new_domain
		}
	}
	//--------------------------------------------------

	return domains_map
}
//--------------------------------------------------
func db__persist_domains(p_domains_map map[string]Domain,
				p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains.db__persist_domains()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	white  := color.New(color.FgWhite).SprintFunc()

	i := 0
	for _,d := range p_domains_map {


		p_runtime_sys.Log_fun("INFO",yellow("persisting ")+white("domain")+yellow(" "+fmt.Sprint(i)+" >---------------- ")+cyan(d.Name_str))

		//IMPORTANT!! -  finds a single document matching the provided selector document 
		//               and modifies it according to the update document. If no document 
		//               matching the selector is found, the update document is applied 
		//               to the selector document and the result is inserted in the collection
		_,err := p_runtime_sys.Mongodb_coll.Upsert(bson.M{"t":"domain","name_str":d.Name_str,},d)
		if err != nil {
			gf_err := gf_core.Error__create("failed to persist a domain in mongodb",
				"mongodb_update_error",
				&map[string]interface{}{"domain_name_str":d.Name_str,},
				err,"gf_domains_lib",p_runtime_sys)
			return gf_err
		}

		i+=1
	}
	return nil
}
//--------------------------------------------------
func db__get_domains(p_runtime_sys *gf_core.Runtime_sys) ([]Domain,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains.db__get_domains()")

	var results_lst []Domain
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{
					"t":        "domain",
					"count_int":bson.M{"$exists":true}, //"count_int" is a new required field, and we want those records, not the old ones
				}).
				Sort("-count_int").
				All(&results_lst)
	if err != nil {
		gf_err := gf_core.Error__create("failed to get all domains",
			"mongodb_find_error",
			nil,err,"gf_domains_lib",p_runtime_sys)
		return nil,gf_err
	}

	return results_lst,nil
}