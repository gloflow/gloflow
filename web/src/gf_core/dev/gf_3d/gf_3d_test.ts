







///<reference path="./../../../d/jquery.d.ts" />

import * as gf_3d  from "./../../ts/gf_3d";

//-------------------------------------------------
$(document).ready(()=>{



    gf_3d.div_follow_mouse($("#target")[0], document, 20);
    gf_3d.div_follow_mouse($("#target2")[0], document, 20);
    gf_3d.div_follow_mouse($("#target3")[0], document, 20);


});