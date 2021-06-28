
gf_picker__main();

//---------------------------------------------------
function gf_picker__main() {

    console.log("gf_page_picker");

    // import jquery if its not defined.
    // testing both condition, because on some sites window.jQuery is defined 
    // but $ is not defined and vice-versa.
    if (!window.jQuery || typeof $ === 'undefined') {
        console.log("GF - jquery not defined - inserting");

        let s = document.createElement("script");
        s.setAttribute('crossorigin', 'anonymous');
        s.setAttribute('integrity',   'sha256-cCueBR6CsyA4/9szpPfrX3s49M9vUU5BgtiJj06wt/s=');
        s.setAttribute('src',         'https://code.jquery.com/jquery-3.1.0.min.js');
        s.setAttribute('type',        'text/javascript');
        document.body.appendChild(s);
        
        s.onload = () => {
            gf_picker__create_ui();
        }
    } else {
        gf_picker__create_ui();
    }
}

//---------------------------------------------------
function gf_picker__create_ui() {
    $("body").append(`
    <style>
    
    div#gf_page_picker {
        
        position: absolute;
        width: 100%;
        height: 100%;
        top: 0px;
        left: 0px;
        
    }

    div#background {
        position: absolute;
        width: 100%;
        height: 100%;
        background-color: gray;
        opacity: 80%;
        top: 0px;
        left: 0px;
    }

    div#gf_url {
        top:      30px;
        left:     19%;
        position: relative;
        background-color: yellow;
        width:   800px;
        padding: 10px;
        padding-bottom: 7px;
    }

    div#gf_url div#description input {
        height:       60px;
        width:        99%;
        border-width: 2px;
        border-color: #808080;
        border-style: solid;
        font-size:    16px;
    }

    div#gf_url div#close_btn {
        position: absolute;
        top: 0px;
        right: -49px;
        width: 50px;
        height: 34px;
        background-color: #fcfd30;
        text-align: center;
        padding-top: 18px;
        color: #7d5d20;

        cursor: pointer;
    }

    div#gf_url div#close_btn:hover {
        background-color: white;
    }

    div#gf_url div#submit {
        background-color: gray;
        text-align: center;
        padding-top: 11px;
        padding-bottom: 10px;
        width: 100%;
        cursor: pointer;
        color: white;
        opacity: 1;
    }

    </style>`);

    $("body").append(`
    <div id="gf_page_picker">
        <div id="background"></div>
    </div>`);

    



    const current_url_str = window.location.href;
    $("#gf_page_picker").append(`

    <div id="gf_url">
        <div id="url">${current_url_str}</div>
        <div id="description">
            <input value="url description"></input>
        </div>
        <div id='close_btn'>x</div>
        <div id='submit'>ok</div> 
    </div>`);


    $("div#gf_url div#close_btn").on('click', function() {

        $("body").find("#gf_page_picker").remove();
    })
}

//---------------------------------------------------