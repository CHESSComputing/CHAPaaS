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
    id.innerHTML += wflow;
    id.innerHTML += "&nbsp; workflow";
    href = "\"javascript:ajaxWorkflowConfig('" + wflow + "')\"";
    id.innerHTML += "&nbsp; <a href=" + href + "id=\"getconfig\">config</a>"
    // update hidden chap input
    document.getElementById("chapworkflow").value = wflow;
}
function AddReader(tag) {
    var id=document.getElementById("workflow");
    if (id) {
        id.className="show";
    }
    id.innerHTML += tag;
    id.innerHTML += "&nbsp; reader &nbsp; &rarr; &nbsp; Processor &nbsp; &rarr; &nbsp;";
    // update hidden chap input
    document.getElementById("reader").value = tag;
}
function AddWriter(tag) {
    var id=document.getElementById("workflow");
    if (id) {
        id.className="show";
    }
    id.innerHTML += tag;
    id.innerHTML += "&nbsp; writer";
    // update hidden chap input
    document.getElementById("writer").value += tag;
}
function RunCHAP(profile) {
    var bid=document.getElementById("base");
    var tid=document.getElementById("token");
    rurl = bid.value+"/chap/run?token="+tid.value;
    if (profile == "profile") {
        rurl = bid.value+"/chap/profile?token="+tid.value;
    }
    var id=document.getElementById("reader");
    if (id) {
        reader = "&reader="+id.value;
        rurl += reader;
    }
    var id=document.getElementById("writer");
    if (id) {
        writer = "&writer="+id.value;
        rurl += writer;
    }
    var id=document.getElementById("chapworkflow");
    if (id) {
        workflow = "&chapworkflow="+id.value;
        rurl += workflow;
    }
    // replace notebook and buttons with run CHAP message
    var did=document.getElementById("chap-notebook")
    if (did) {
        did.innerHTML = "<div class=\"img-center\">";
        did.innerHTML += "<img src=\"/images/wait.gif\" alt=\"please wait for its completion\" width=\"48\">";
        did.innerHTML += "Your new CHAP workflow is scheduled and run right now";
        did.innerHTML += "</div>";
    }
    var xid=document.getElementById("chap-buttons")
    if (xid) {
        xid.innerHTML = ""
    }
    console.log("will call "+rurl);
    // execute rurl call to our server
    window.onbeforeunload = null;
    window.location.href = rurl;
}
