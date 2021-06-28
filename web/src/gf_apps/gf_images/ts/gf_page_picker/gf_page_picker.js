














function gf_picker__main() {




    console.log("gf works");


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




gf_picker__main();