function ajaxWorkflowConfig(wflow) {
    $(document).ready(function(){
//      $("#getconfig").click(function(){
        rurl = "/chap/config/"+wflow;
        $.get(rurl, function(data, status){
            var id=document.getElementById("workflowconfig");
            if (id) {
                id.className="show";
            }
            id.innerHTML = "<a href=\"javascript:HideTag('workflowconfig')\" class=\"button button-small button-round\">Hide</a>";
            id.innerHTML += "<a href=\"javascript:HideTag('workflowconfig')\" class=\"button button-small button-round\">Save</a>";
            id.innerHTML += "<div><textarea class=\"wflowconfig\">" + data + "</textarea></div><br/>";
        });
//      });
    });
}