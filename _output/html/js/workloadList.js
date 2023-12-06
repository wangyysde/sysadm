function workloadDoAjax(actionUrl,actionType,isMultiPart) {
    var data;
    if(isMultiPart){
        data = new FormData($("#addObjectForm")[0]);
    } else {
        data = "";
    }
    var ajaxRet = {"msg": "未知错误", "error": true};
    $.ajax({
        type: actionType,
        // dataType: "json",
        // crossDomain: true,
        contentType: false,
        cache: false,
        processData: false,
        url: actionUrl,
        data: data,
        async: false,
        error: function(xmlObj, request) {
            var msg = "";
            if (xmlObj.responseText == null) {
                msg = "网络出错";
            } else {
                msg = xmlObj.responseText;
            }
            ajaxRet.msg =msg;
            ajaxRet.error = true;
        },
        success: function(result) {
            if (result.errorCode == 0) {
                var messageArray = result.message;
                var msg = "未知错误";
                if (messageArray[0]) {
                    var msgLine = messageArray[0];
                    msg = b64DecodeUnicode(msgLine["msg"]);

                }
                ajaxRet.msg = msg;
                ajaxRet.error = false;
            } else {
                var messageArray = result.message;
                var msg = "未知错误";
                if (messageArray[0]) {
                    var msgLine = messageArray[0];
                    msg = b64DecodeUnicode(msgLine["msg"]);
                }
                ajaxRet.msg = msg;
                ajaxRet.error = true;
            }
        }
    });

    return ajaxRet;
}

function displayTip(msg,color,timeout){
    var  colorStr = "#1f6f4a";
    if(color.toUpperCase() == "RED"){
        colorStr = "#ff0000";
    }

    msg=msg.replace(/^\s+|\s+$/g," ");
    if(msg == ""){
        return;
    }

    if(timeout == 0){
        timeout = 5000;
    }
    var tip = document.getElementById("tip");
    tip.style.display = "block";
    tip.style.backgroundColor = colorStr;
    tip.innerHTML = msg;
    setTimeout(function() {
        var tip = document.getElementById("tip");
        tip.style.display = "none";
    }, timeout);

    return;
}

function ChangeSubOptions(actionUri,subObjId,objId,objName) {
    actionUri = actionUri.replace(/^\s+|\s+$/g," ");
    subObjId = subObjId.replace(/^\s+|\s+$/g," ");
    objId = objId.replace(/^\s+|\s+$/g," ");
    if(actionUri == "" || subObjId == "" || objId == ""){
        return;
    }

    var uri = pageUrl + actionUri;
    uri = uri + "?objID=" + objId;
    var respValue = workloadDoAjax(uri,"GET",false);
    if(respValue.error){
        return displayTip(respValue.msg,"RED",0);
    }

    var subObjDataStr = respValue.msg;
    var subObjLines = subObjDataStr.split(",");
    var subObj = document.getElementById(subObjId);
    if(!subObj){
        return ;
    }

    subObj.options.length = 0;
    var newOption;
    for(i =0; i < subObjLines.length; i++ ){
        var lineData  = subObjLines[i].split(":");
        if(lineData.length < 2){
            continue;
        }
        newOption = new Option(lineData[1], lineData[0]);
        subObj.options.add(newOption);
    }

    if(objId == "0"){
        var nsObj = document.getElementById("namespace");
        if(nsObj){
            nsObj.options.length = 0;
            newOption = new Option("==选择命名空间===", "0");
            nsObj.options.add(newOption);
        }
    }

    reloadListPage("0","","","",objName);
    return;
}

function ChangeNSOptions(actionUri,subObjId,objValue,objName) {
    if(objValue == "0"){
        var subObj = document.getElementById(subObjId);
        subObj.options.length = 0;
        newOption = new Option("==选择命名空间===", "0");
        subObj.options.add(newOption);
    }
    reloadListPage("0","","","",objName);
    return;
}
function reloadListPage(start,orderfield,direction,searchContent,objName) {
    var dcObj = document.getElementById("dcID");
    var clusterObj = document.getElementById("clusterID");
    var nsObj = document.getElementById("namespace");

    var paras = objName + "/list?dcID=" + dcObj.options[dcObj.selectedIndex].value + "&clusterID=" + clusterObj.options[clusterObj.selectedIndex].value;
    if(nsObj){
        paras = paras + "&namespace=" +nsObj.options[nsObj.selectedIndex].value;
    }
    paras = paras +"&start=" + start + "&orderfield=" + orderfield + "&direction=" + direction;
    var url = pageUrl + paras;
    $('#container').load(url);

    return;
}

function workloadListShowPopMenu(event,itemids,objid,objName) {
    var items = itemids.split(",");
    var menuContent = "";

    for(i = 0; i < items.length; i++){
        var popmenustr = popMenuItems[items[i]];
        if(popmenustr) {
            var popmenuArray = popmenustr.split(",");
            if (i > 0) {
                menuContent = menuContent + '<br><li><a href="#" onclick=\'doMenuItemForWorkloadList("' + objid + '","' + popmenuArray[1] + '","' + popmenuArray[2] + '","' + popmenuArray[3] + '","' + objName + '")\'>' + popmenuArray[0] + '</a></li>';
            } else {
                menuContent = menuContent + '<li><a href="#" onclick=\'doMenuItemForWorkloadList("' + objid + '","' + popmenuArray[1] + '","' + popmenuArray[2] + '","' + popmenuArray[3] + '","' + objName + '")\'>' + popmenuArray[0] + '</a></li>';
            }
        }
    }
    var popMenu = document.getElementById("popmenu");
    popMenu.innerHTML = menuContent;
    var posX = event.clientX;
    var posY = event.clientY;
    popMenu.style.display = "block";
    popMenu.style.zIndex = 2100;
    //popMenu.style.left = (posX - 100) + "px";
    popMenu.style.left = (posX - 50) + "px";
    popMenu.style.top = posY + "px";
}

// this function be called when user click a item of popmenu
function doMenuItemForWorkloadList(objID,action,actionType,method,objName){
    var dcObj = document.getElementById("dcID");
    var clusterObj = document.getElementById("clusterID");
    var nsObj = document.getElementById("namespace");

    var paras = "?dcID=" + dcObj.options[dcObj.selectedIndex].value + "&clusterID=" + clusterObj.options[clusterObj.selectedIndex].value;
    if(nsObj){
        paras = paras + "&namespace=" +nsObj.options[nsObj.selectedIndex].value;
    }


    var popMenu = document.getElementById("popmenu");
    popMenu.style.display = "none";
    popMenu.style.zIndex = 2000;

    var doUrl = pageUrl + objName + "/" + action + paras + "&objID=" + objID;

    if(method == "poppage"){
        doAjaxForPoppage(doUrl,actionType);
        return;
    }
    if(method == "tip") {
        var respValue =  workloadDoAjax(doUrl, actionType,false);
        var refreshUri = pageUrl + objName + "/list"  + paras;
        $('#container').load(refreshUri);
        if(respValue["error"]){
            return displayTip(respValue["msg"],"RED",5000);
        }

        return displayTip(respValue["msg"],"GREEN",5000);
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
