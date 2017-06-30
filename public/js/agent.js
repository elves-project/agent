  var dashboard = {};
  dashboard.getGen = function () {
    $.get("/stat/general", function (data) {
        if (data.msg == "success") {
            $("#agent-mode").text(data.data['Mode']);
            if(data.data['Mode']=="DEVELOP"){
                $("#agent-mode").addClass("label label-warning");
            }else{
                $("#agent-mode").addClass("label label-success");
            }
            $("#agent-asset").text(data.data['Asset']);
            $("#agent-ip").text(data.data['Ip']);
            $("#agent-uptime").text(data.data['Uptime']);
            $("#agent-heartbeat").text(data.data['Hbtime']);
            $("#agent-version").text(data.data['Ver']);
        }
    }, "json");
 }

 dashboard.getApps = function () {
    $.get("/stat/apps", function (data) {
      var tr = ""
      for(var key in data.data)
      {
          tr = tr+"<tr><td>"+key+"</td><td>"+data.data[key]+"</td></tr>"
      }
      $("#apps_table").html(tr);
    }, "json");
 }

 dashboard.getCrons = function () {
    $.get("/stat/crons", function (data) {
      var tr = ""
      for(var key in data.data)
      {
          tr = tr+"<tr><td>"+key+"</td><td>"+data.data[key].App+"</td><td>"+data.data[key].Func+"</td><td>"+data.data[key].Mode+"</td><td>"+data.data[key].Rule+"</td><td>"+data.data[key].Comment+"</td><td>"+data.data[key].Lastexec+"</td></tr>"
      }
      $("#crons_table").html(tr);
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
    all: dashboard.getAll,
    gen: dashboard.getGen,
    app: dashboard.getApps,
    cron: dashboard.getCrons
};

$(document).ready(function() {
    dashboard.getAll();
});