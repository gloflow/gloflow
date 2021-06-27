



// save just the JS in a browser bookmark
javascript:void((function(d){
    
    d.addEventListener('securitypolicyviolation', function(r){
        alert('ContentSecurityPolicyError!');
    });
    e=d.createElement('script');
    e.setAttribute('type','text/javascript');
    e.setAttribute('charset','UTF-8');
    e.setAttribute('debug','true');
    e.setAttribute('src','//gloflow.com/images/static/js/gf_bookmarklet.js?r='+(Math.random()*99999999));
    d.body.appendChild(e);
}(document)));