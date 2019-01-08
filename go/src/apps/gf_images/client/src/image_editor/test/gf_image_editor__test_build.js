///<reference path="../d/jquery.d.ts" />
var gf_image_editor;
(function (gf_image_editor) {
    //-------------------------------------------------
    function init(p_target_image_div_element, p_log_fun) {
        p_log_fun('FUN_ENTER', 'gf_image_editor.init()');
        var target_image = $(p_target_image_div_element).find('img')[0];
        var width_int    = target_image.clientWidth;
        var height_int   = target_image.clientHeight;
        console.log('img width  - ' + width_int);
        console.log('img height - ' + height_int);
        var container = $("\n\t\t<div class='gf_image_editor'>\n\t\t\t<div class='open_editor_btn'>editor</div>\n\t\t</div>");
        $(p_target_image_div_element).append(container);
        //-------------------------------------------------
        function create_pane() {
            p_log_fun('FUN_ENTER', 'gf_image_editor.init().create_pane()');
            var editor_pane = $("\n\t\t\t<div class='editor_pane'>\n\t\t\t\t<div class='close_btn'>x</div>\n\t\t\t\t<div class='save_btn'>save</div>\n\n\t\t\t\t<canvas class='modified_image_canvas' width=\"" + width_int + "\" height=\"" + height_int + "\"></canvas>\n\n\t\t\t\t<div class=\"slider_input\">\n\t\t\t\t\t<form>\n\t\t\t\t\t\t<div>\n\t\t\t\t\t\t\t<input id=\"contrast\" name=\"contrast\" type=\"range\" min=\"-100\" max=\"100\" value=\"0\">\n\t\t\t\t\t\t\t<label for=\"contrast\">contrast</label>\n\t\t\t\t\t\t</div>\n\t\t\t\t\t\t\n\t\t\t\t\t\t<div>\n\t\t\t\t\t\t\t<input id=\"brightness\" name=\"brightness\" type=\"range\" min=\"-100\" max=\"100\" value=\"0\">\n\t\t\t\t\t\t\t<label for=\"brightness\">brightness</label>\n\t\t\t\t\t\t</div>\n\t\t\t\t\t\t\n\t\t\t\t\t\t<div>\n\t\t\t\t\t\t\t<input id=\"saturation\" name=\"saturation\" type=\"range\" min=\"-100\" max=\"100\" value=\"0\">\n\t\t\t\t\t\t\t<label for=\"saturation\">saturation</label>\n\t\t\t\t\t\t</div>\n\n\t\t\t\t\t\t<div>\n\t\t\t\t\t\t\t<input id=\"sharpen\" name=\"sharpen\" type=\"range\" min=\"0\" max=\"100\" value=\"0\">\n\t\t\t\t\t\t\t<label for=\"sharpen\">sharpen</label>\n\t\t\t\t\t\t</div>\n\n\t\t\t\t\t\t<div>\n\t\t\t\t\t\t\t<input id=\"sepia\" name=\"sepia\" type=\"range\" min=\"0\" max=\"100\" value=\"0\">\n\t\t\t\t\t\t\t<label for=\"sepia\">sepia</label>\n\t\t\t\t\t\t</div>\n\n\t\t\t\t\t\t<div>\n\t\t\t\t\t\t\t<input id=\"noise\" name=\"noise\" type=\"range\" min=\"0\" max=\"100\" value=\"0\">\n\t\t\t\t\t\t\t<label for=\"noise\">noise</label>\n\t\t\t\t\t\t</div>\n\n\t\t\t\t\t\t<div>\n\t\t\t\t\t\t\t<input id=\"hue\" name=\"hue\" type=\"range\" min=\"0\" max=\"100\" value=\"0\">\n\t\t\t\t\t\t\t<label for=\"hue\">hue</label>\n\t\t\t\t\t\t</div>\n\t\t\t\t\t</form>\n\t\t\t\t</div>\n\t\t\t</div");
            $(editor_pane).find('input[type=range]').change(apply_filters_fun);
            //-------------------------------------------------
            function apply_filters_fun() {
                var contrast   = parseInt($('#contrast').val());
                var brightness = parseInt($('#brightness').val());
                var saturation = parseInt($('#saturation').val());
                var sharpen    = parseInt($('#sharpen').val());
                var sepia      = parseInt($('#sepia').val());
                var noise      = parseInt($('#noise').val());
                var hue        = parseInt($('#hue').val());
                Caman('.editor_pane canvas', target_image, function () {
                    this.revert(false);
                    this.contrast(contrast);
                    this.brightness(brightness);
                    this.saturation(saturation);
                    this.sharpen(sharpen);
                    this.sepia(sepia);
                    this.noise(noise);
                    this.hue(hue);
                    this.render(function () { return console.log('filter applied'); });
                });
            }
            //-------------------------------------------------
            var canvas = $(editor_pane).find('canvas')[0];
            Caman(canvas, $(target_image).attr('src'), function () {
                this.render();
            });
            //-------------
            //SAVE_MODIFIED_IMAGE
            $(editor_pane).find('.save_btn').on('click', function () {
                save_modified_image(editor_pane);
            });
            //-------------
            return editor_pane;
        }
        //-------------------------------------------------
        function save_modified_image(p_editor_pane) {
            p_log_fun('FUN_ENTER', 'gf_image_editor.init().save_modified_image()');
            var canvas            = $(p_editor_pane).find('.modified_image_canvas')[0];
            var canvas_base64_str = canvas.toDataURL();
            console.log(canvas_base64_str);
        }
        //-------------------------------------------------
        var opened_bool = false;
        $(container).find('.open_editor_btn').on('click', function () {
            if (opened_bool) {
                return;
            }
            var editor_pane = create_pane();
            $(editor_pane).find('.close_btn').on('click', function () {
                $(editor_pane).remove();
                opened_bool = false;
            });
            $(container).append(editor_pane);
            opened_bool = true;
        });
        return container;
    }
    gf_image_editor.init = init;
    //-------------------------------------------------
})(gf_image_editor || (gf_image_editor = {}));
