package gf_images_lib

//-------------------------------------------------
type Config struct {
	Flow_to_s3bucket_map map[string]string
}
//-------------------------------------------------
func Config__get() Config {

	flow_to_s3bucket_map := map[string]string{
		"general"   :"gf--img",
		"discovered":"gf--img--discover",
	}

	config := Config{
		Flow_to_s3bucket_map:flow_to_s3bucket_map,
	}

	return config
}