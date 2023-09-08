function ajaxWorkflowConfig(wflow) {
    $(document).ready(function(){
        rurl = "/chap/config/"+wflow;
        $.get(rurl, function(data, status){
            var id=document.getElementById("workflowconfig");
            if (id) {
                id.className="show";
            }
            id.innerHTML = "<a href=\"javascript:HideTag('workflowconfig')\" class=\"button button-small button-round\">Hide</a>";
            id.innerHTML += "<a href=\"javascript:ajaxSaveWorkflowConfig()\" class=\"button button-small button-round\">Save</a>";
            id.innerHTML += "<div><textarea id=\"config-textarea\" class=\"wflowconfig\">" + data + "</textarea></div><br/>";
        });
    });
}
// helper function to get doc response
function ajaxDocResponse(doc) {
    $(document).ready(function(){
        rurl = "/chap/doc/"+doc;
        $.get(rurl, function(data, status){
            var id=document.getElementById("doc-response");
            if (id) {
                id.className="show";
            }
            id.innerHTML = data;
        });
    });
}
function ajaxSaveWorkflowConfig() {
    HideTag('workflowconfig')
    var bid=document.getElementById("base");
    var id=document.getElementById('config-textarea');
    if (!id) {
        return
    }
    wflow = id.value;
    console.log("send POST request with worfklow config:\n"+wflow);
    var pid=document.getElementById("chapworkflow");
    if (!pid) {
        return
    }
    // make POST request to backend server to save config content
    workflow = pid.value;
    rurl = bid.value+"/chap/config/" + workflow;
    $(document).ready(function(){
        $.post(rurl,
            wflow,
            function(data, status){
            console.log("HTTP POST responst: "+status);
        });
        /*
        $.ajax({
            contentType: 'application/json',
            data: {
                "content": id.value
            },
            dataType: 'json',
            success: function(data){
                console.log("device control succeeded");
            },
            error: function(){
                console.log("Device control failed");
            },
            processData: false,
            type: 'POST',
            url: rurl
        });
        */
    });
}
// helper function to generate tar ball for given workflow
function ajaxGenTarball(user, wflow) {
    $(document).ready(function(){
        // send request to /chap/tar/:workflow
        rurl = "/chap/tar/"+wflow;
        $.get(rurl, function(data, status){
            var id=document.getElementById("tarball");
            if (id) {
                id.className="show";
            }
            // generate HTML link
            tball = wflow+'.tar.gz';
            href = '/usrs/'+user+'/'+tball;
            id.innerHTML = 'tar-ball '+tball+' has been generated: '+'<a href="'+href+'">download</a>';
            // http://localhost:8182/users/dev-user/basic.tar.gz
        });
    });
}
