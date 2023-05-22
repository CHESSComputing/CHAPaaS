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
function AddReader(tag) {
    var id=document.getElementById("workflow");
    if (id) {
        id.className="show";
    }
    id.innerHTML += tag;
    id.innerHTML += "&nbsp; reader &nbsp; &rarr; &nbsp; Processor &nbsp; &rarr; &nbsp;";
}
function AddWriter(tag) {
    var id=document.getElementById("workflow");
    if (id) {
        id.className="show";
    }
    id.innerHTML += tag;
    id.innerHTML += "&nbsp; writer";
}
