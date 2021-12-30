

/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

function create(trigger_y_position, name, screen_height, activate_fn, deactivate_fn) {

    console.log("SCROLL_TRIGGER CREATE")

    $("body").append("<div id='"+name+"'></div>"); // place between divs div<!--"+name+"--> if you want letters on triggers
    $("body").find("#"+name).css({
        position: "absolute",
        right:    "0px",
        top:      trigger_y_position+"px",
        //width:    "10px",
        //height:   "2px",
        //"background-color": "yellow",
        "z-index": 20

    })

    var active = false;
    $(window).scroll(function(e){
        // console.log(e);

        // console.log(window.scrollY)

        // checking if the bottom of the screen has passed the trigger, in order to activate it.
        // only activate the trigger when its not active (active == false)
        var bottom_scroll_y = window.scrollY + screen_height;

        if (bottom_scroll_y > trigger_y_position && active == false) {
            active = true

            activate_fn();
        }

        if (bottom_scroll_y < trigger_y_position && active == true) {
            active = false

            deactivate_fn();
        }

    });

}


function remove_triggers(name){
    var target_trigger = $("body").find("#"+name)
    $(target_trigger).remove();
}