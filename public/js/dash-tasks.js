 var dashboard = {};
 dashboard.getTasks = function () {
    $.get("/stat/tasks", function (data) {
        if (data.msg == "success") {
          var tr = ""
          $.each(data.data, function(){
                   if (this.ID!=""){
                       if(this.Flag=="failure"){
                           cls = "label label-info"
                       }else if(this.Flag=="error"){
                           cls = "label label-warning"
                      }else{
                           cls = "label label-success"
                      }
                      tr = tr+"<tr><td>"+this.Time+"</td><td>"+this.ID+"</td><td>"+this.Type+"</td><td>"+this.Mode+"</td><td>"+this.App+"</td><td>"+this.Func+"</td><td>"+this.Proxy+"</td><td>"+this.Costtime+"ms</td><td><span class='"+cls+"'>"+this.Flag+"<span></td></tr>"
                   }
          });
          $("#tasks_table").html(tr);
        }
    }, "json");
 }

dashboard.getAll = function () {
    for (var item in dashboard.fnMap) {
        if (dashboard.fnMap.hasOwnProperty(item) && item !== "all") {
            dashboard.fnMap[item]();
        }
    }
}

dashboard.fnMap = {
    tasks: dashboard.getTasks
};

$(document).ready(function() {
    dashboard.getAll();
});