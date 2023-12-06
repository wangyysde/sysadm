window.addEventListener("click", closePopMenu);
function showPopMenu(event,itemids,objid) {
    var items = itemids.split(",");
    var menuContent = "";
    for(i = 0; i < items.length; i++){
        var popmenustr = popMenuItems[items[i]];
        if(popmenustr) {
            var popmenuArray = popmenustr.split(",");
            if (i > 0) {
                menuContent = menuContent + '<br><li><a href="#" onclick=\'doMenuItem("' + objid + '","' + popmenuArray[1] + '","' + popmenuArray[2] + '","' + popmenuArray[3] + '")\'>' + popmenuArray[0] + '</a></li>';
            } else {
                menuContent = menuContent + '<li><a href="#" onclick=\'doMenuItem("' + objid + '","' + popmenuArray[1] + '","' + popmenuArray[2] + '","' + popmenuArray[3] + '")\'>' + popmenuArray[0] + '</a></li>';
            }
        }
    }
    var popMenu = document.getElementById("popmenu");
    popMenu.innerHTML = menuContent;
    var posX = event.clientX;
    var posY = event.clientY;
    popMenu.style.display = "block";
    popMenu.style.zIndex = 2100;
    popMenu.style.left = (posX - 100) + "px";
    popMenu.style.top = posY + "px";
}

function closePopMenu(e){
    if (e.target.id != "popmenuid" ) {
        var popMenu = document.getElementById("popmenu");
        if(popMenu && popMenu.style.display == "block"){
            popMenu.style.display = "none";
            popMenu.style.zIndex = 2000;
        }
    }
}

// this function be called when user click a item of popmenu
function doMenuItem(objID,action,actionType,method){
    var popMenu = document.getElementById("popmenu");
    popMenu.style.display = "none";
    popMenu.style.zIndex = 2000;

    var doUrl = pageUrl + action+"?objID=" + objID;

    if(method == "poppage"){
        doAjaxForPoppage(doUrl,actionType);
        return;
    }
    if(method == "tip") {
        doAjax(doUrl, actionType);
        return;
    }

    if(method == "page"){
        $('#container').load(doUrl);
        return;
    }

    if(method == "window"){
        window.open(doUrl,"_blank");
        return;
    }

}

function closeDetailPage(){
    var poppage = document.getElementById("detailPoppage");
    poppage.style.display = "none";
    poppage.style.zIndex = 0;
    var maskLayer = document.getElementById("maskLayer");
    maskLayer.style.zIndex = 0;
}

function openPoppage(content){
    var maskLayer = document.getElementById("maskLayer");
    maskLayer.style.zIndex = 2000;
    var poppage = document.getElementById("detailPoppage");
    poppage.style.display = "block";
    poppage.style.zIndex = 2100;
    poppage.innerHTML = content;
    setAceEditorValueForDetails();

    poppage.focus();

    return;
}
// check or uncheck all all object checkbox on object list page
function selectAllObjectCheckbox(isChecked) {
    if(isChecked) {
        var chklist = document.getElementsByName('objectid[]');
        for (var i = 0; i < chklist.length; i++) {
            chklist[i].checked = true;
        }
    } else {
        var chklist = document.getElementsByName('objectid[]');
        for (var i = 0; i < chklist.length; i++) {
            chklist[i].checked = false;
        }
    }
}

function selectObjectCheckbox(isChecked){
    if(isChecked){

        var chklist = document.getElementsByName('objectid[]');
        var allChecked = true;
        for (var i = 0; i < chklist.length; i++) {
            if(!chklist[i].checked){
                allChecked = false;
            }
        }
        if(allChecked){
            var hostth = document.getElementById("objectListTH");
            hostth.checked = true;
        }
        return;
    }

    var allNotChecked = true;
    var chklist = document.getElementsByName('objectid[]');
    for (var i = 0; i < chklist.length; i++) {
        if(chklist[i].checked){
            allNotChecked = false;
        }
    }

    if(allNotChecked){
        var hostth = document.getElementById("objectListTH");
        hostth.checked = false;
    }
}

function doAjax(actionUrl,actionType) {
    $.ajax({
        type: actionType,
        // dataType: "json",
        // crossDomain: true,
        contentType: false,
        cache: false,
        processData: false,
        url: actionUrl,
        data: "",
        //async: false,
        error: function(xmlObj, request) {
            var errMsg = "";
            if (xmlObj.responseText == null) {
                errMsg = "Connection error";
            } else {
                errMsg = xmlObj.responseText;
            }
            var tip = document.getElementById("tip");
            tip.innerHTML = errMsg;
            tip.style.display = "block";
            setTimeout(function() {
                var tip = document.getElementById("tip");
                tip.style.display = "none";
            }, 5000);
        },
        success: function(result) {
            var tip = document.getElementById("tip");
            if (result.errorCode == 0) {
                var message = result.message;
                var msg = "未知错误";
                if (message != "") {
                   msg = message;
                }
                tip.style.display = "block";
                tip.style.backgroundColor = "#1f6f4a";
                tip.style.display = "block";
                tip.innerHTML = msg;
                setTimeout(function() {
                    var tip = document.getElementById("tip");
                    tip.style.display = "none";
                }, 5000);

            } else {
                var message = result.message;
                var msg = "未知错误";
                if (message != "") {
                    msg = message;
                }
                tip.innerHTML = msg;
                tip.style.display = "block";
                setTimeout(function() {
                    var tip = document.getElementById("tip");
                    tip.style.display = "none";
                }, 5000);
            }
        }
    });
}

function doAjaxForPoppage(actionUrl,actionType) {
    $.ajax({
        type: actionType,
        contentType: false,
        cache: false,
        processData: false,
        url: actionUrl,
        data: "",
        //async: false,
        error: function(xmlObj, request) {
            var errMsg = "";
            if (xmlObj.responseText == null) {
                errMsg = "连接服务器出错";
            } else {
                errMsg = xmlObj.responseText;
            }
            openPoppage("<span>" + errMsg + "</span>");
        },
        success: function (result){
    //        msg = xmlObj.responseText;
            openPoppage(result);
        }
       // success: function(result) {
       //     openPoppage(result);
       // }
    });
}

function listContenChanged(urlparas) {
    var url = pageUrl + "list?" + urlparas;
    $('#container').load(url);
}

function GroupSelectChanged(groupID){
    var urlParas ="groupSelectID=" + groupID;
    listContenChanged(urlParas);
}

function doSearch(searchContent) {
    var urlParas = "searchContent=" + searchContent;
    listContenChanged(urlParas);
}

function displayAddObjectForm(module){
    var dcIDObj = document.getElementById("dcID");
    var dcIDValue = dcIDObj.options[dcIDObj.options.selectedIndex].value;
    var clusterID = document.getElementById("clusterID");
    var clusterValue =  clusterID.options[clusterID.options.selectedIndex].value;
    var nsValue = "";
    var nsObj = document.getElementById("namespace");
    if(nsObj){
        nsValue = nsObj.options[nsObj.options.selectedIndex].value;
    }
    var url = pageUrl + module + "/addform" ;
    url = url + "?dcID=" + dcIDValue + "&clusterID=" + clusterValue + "&namespace=" + nsValue;
    $('#container').load(url);
}

function addCluster(){
    $('#container').load("/k8scluster/addform");
}