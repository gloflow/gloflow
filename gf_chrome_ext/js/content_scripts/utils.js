




function hash_code(p_str,
			p_log_fun) {
	var hash = 0;
	if (p_str.length == 0) return hash;
	for (i = 0; i < p_str.length; i++) {
		char = p_str.charCodeAt(i);
		hash = ((hash<<5)-hash)+char;
		hash = hash & hash; // Convert to 32bit integer
	}
	return hash;
}
//---------------------------------------------------
function get_image_histogram(p_on_complete_fun,
						p_log_fun) {
	p_log_fun('FUN_ENTER','utils.get_image_histogram()');

	const hist_map = {};	
	const img       = new Image();
	img.onload      = () => {
		Pixastic.process(img,"colorhistogram", {
			paint      :true,
			returnValue:hist_map
		});
		hist_map.rvals; // <- array[255] red channel
		hist_map.gvals; // <- array[255] green channel
		hist_map.bvals; // <- array[255] blue channel

		p_on_complete_fun(hist_map);
	}
	document.body.appendChild(img);
	img.src = "myimage.jpg";
}