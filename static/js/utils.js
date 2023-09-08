function HideTag(tag) {
    var id=document.getElementById(tag);
    if (id) {
        id.className="hide";
    }
}
function ShowTag(tag) {
    var id=document.getElementById(tag);
    if (id) {
        id.className="show";
    }
}
function FlipTag(tag) {
    var id=document.getElementById(tag);
    if (id) {
        if  (id.className == "show") {
            id.className="hide";
        } else {
            id.className="show";
        }
    }
}
function AddWorkflow(wflow) {
    var id=document.getElementById("workflow");
    if (id) {
        id.className="show";
    }
    id.innerHTML = wflow;
    id.innerHTML += "&nbsp; workflow";
    href = "\"javascript:ajaxWorkflowConfig('" + wflow + "')\"";
    id.innerHTML += "<br/>show <a href=" + href + "id=\"getconfig\">config</a>"
    // update hidden chap input
    var cid=document.getElementById("chapworkflow");
    if (cid) {
        cid.value = wflow;
    }
}
function RunCHAP(profile) {
    var bid=document.getElementById("base");
    var tid=document.getElementById("token");
    rurl = bid.value+"/chap/run?token="+tid.value;
    if (profile == "profile") {
        rurl = bid.value+"/chap/profile?token="+tid.value;
    } else if (profile == "batch") {
        rurl = bid.value+"/chap/batch?token="+tid.value;
    }
    var id=document.getElementById("chapworkflow");
    if (id) {
        workflow = "&chapworkflow="+id.value;
        rurl += workflow;
    }
    // replace notebook and buttons with run please wait message
    HideTag("notebook");
    HideTag("chap-buttons");
    ShowTag("please-wait");
    console.log("will call "+rurl);
    // execute rurl call to our server
    window.onbeforeunload = null;
    window.location.href = rurl;
}
function DocResponse(doc) {
    // replace notebook and buttons with run please wait message
    HideTag("notebook");
    HideTag("chap-buttons");
    HideTag("please-wait");
    var arr = ["doc-response", "workflowconfig", "workflow"];
    for (var i = 0; i < arr.length; i++) {
        var id=document.getElementById(arr[i]);
        if (id) {
            id.innerHTML = "";
        }
    }
    ShowTag("doc-response");
    ajaxDocResponse(doc)
}
function ShowNotebook() {
    HideTag("please-wait");
    var arr = ["doc-response", "workflowconfig"];
    for (var i = 0; i < arr.length; i++) {
        var id=document.getElementById(arr[i]);
        if (id) {
            id.innerHTML = "";
        }
    }
    ShowTag("notebook");
    ShowTag("chap-buttons");
}
function HideMenu() {
    HideTag('workflows');
    HideTag('readers');
    HideTag('writers');
    HideTag('processors');
}
