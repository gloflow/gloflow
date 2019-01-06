package gf_core
//-------------------------------------------------------------
func Str_in_lst(p_str string,
			p_lst []string) bool {
	for _,s := range p_lst {
		if p_str == s {
			return true
		}
	}
	return false
}
//-------------------------------------------------------------
func Map_has_key(p_map interface{},
			p_key_str string) bool {

	if _,ok := p_map.(map[string]interface{}); ok {
		_,ok := p_map.(map[string]interface{})[p_key_str]
		return ok
	} else if _,ok := p_map.(map[string]string); ok {
		_,ok := p_map.(map[string]string)[p_key_str]
		return ok
	}

	panic("no handler for p_map type")
	return false
}