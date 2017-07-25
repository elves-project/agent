$(document).ready(function() {
    $("#dev_app").val(read_cookie("dev_app"))
    $("#dev_func").val(read_cookie("dev_func"))
    if(read_cookie("dev_proxy")){
        $("#dev_proxy").val(read_cookie("dev_proxy").replace("%7C","|"))
    }
    $("#dev_timeout").val(read_cookie("dev_timeout"))
    if(read_cookie("dev_param")){
      $("#dev_param").val(decodeURI(read_cookie("dev_param")).replace("%3A",":"))
    }
    //console.log(read_cookie("dev_mode"))
    if(read_cookie("dev_mode")=="P"){
        $("#dev_mode option").eq(1).attr('selected', 'true');
    }
});

function getapiuri(){
  if($("#dev_app").val()=="" || $("#dev_func").val()=="" ){
      alert("param error!!");
      return;
  }
  $.get("/api/gettesturl?app="+$("#dev_app").val()+"&func="+$("#dev_func").val()+"&param="+$("#dev_param").val()+"&proxy="+$("#dev_proxy").val()+"&timeout="+$("#dev_timeout").val(), function (data) {
        if (data.msg == "success") {
            if(data.data.status=="true"){
                write_cookie("dev_app",$("#dev_app").val());
                write_cookie("dev_func",$("#dev_func").val());
                write_cookie("dev_proxy",$("#dev_proxy").val());
                write_cookie("dev_timeout",$("#dev_timeout").val());
                write_cookie("dev_param",$("#dev_param").val());
                url = window.location.protocol+"//"+window.location.host+"/";
                $("#dev_uri").val(url+"api/v2/rt/exec")
                $("#dev_post").val(data.data.uri)
            }else{
                alert("Only Dev Mode Can Use This Service.")
            }
        }
    }, "json");
}

function gandr(){
    if($("#dev_app").val()=="" || $("#dev_func").val()=="" ){
        alert("param error!!");
        return;
    }
    $.get("/api/gettesturl?app="+$("#dev_app").val()+"&func="+$("#dev_func").val()+"&param="+$("#dev_param").val()+"&proxy="+$("#dev_proxy").val()+"&timeout="+$("#dev_timeout").val(), function (data) {
        if (data.msg == "success") {
            if(data.data.status=="true"){
                write_cookie("dev_app",$("#dev_app").val());
                write_cookie("dev_func",$("#dev_func").val());
                write_cookie("dev_proxy",$("#dev_proxy").val());
                write_cookie("dev_timeout",$("#dev_timeout").val());
                write_cookie("dev_param",$("#dev_param").val());
                $("#dev_uri").val(window.location.protocol+"//"+window.location.host+"/"+"api/v2/rt/exec");
                $("#dev_post").val(data.data.uri);
                $.get($("#dev_uri").val()+"?"+$("#dev_post").val(), function(result){
                    //console.log(result);
                    $(".rst").html(JSON.stringify(result, null, 4));
                });
            }else{
                alert("Only Dev Mode Can Use This Service.")
            }
        }
    }, "json");
}

function runtest(){
  $.get($("#dev_uri").val()+"?"+$("#dev_post").val(), function(result){
      $(".rst").html(JSON.stringify(result, null, 4));
  });
}

function runtestnew(){
    window.open($("#dev_uri").val()+"?"+$("#dev_post").val());
}

function read_cookie(key){
            var str,ary;
            str=document.cookie;
            ary=str.replace(/ *; */g,";").split(";");
            key=escape(key)+"=";
            for(var i=0;i<ary.length;i++){
                if(ary[i].indexOf(key)==0){
                     return ((ary[i].split("=")[1]));
                }
            }
}

function  write_cookie(key,value,cookieDomain,cookiePath,expireTime,targetWindow){
            var strAppendix="";
            strAppendix+=cookieDomain?";domain="+cookieDomain:"";
            strAppendix+=cookiePath?";path="+cookiePath:"";
            strAppendix+=expireTime?";expires="+expireTime:"";
            targetWindow=targetWindow?targetWindow:top;
            targetWindow.document.cookie=escape(key)+"="+escape(value)+strAppendix;
}