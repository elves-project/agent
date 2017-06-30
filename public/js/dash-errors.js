 var dashboard = {};
  dashboard.getErrors = function () {
    $.get("/stat/errors", function (data) {
        if (data.msg == "success") {
            var tr = ""
            $.each(data.data, function(){
                  if (this!=""){
                    tr = tr+"<li>"+this+"</li>"
                  }
             });
             $("#errors_table").html(tr);
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
    errors: dashboard.getErrors
};

$(document).ready(function() {
    dashboard.getAll();
});