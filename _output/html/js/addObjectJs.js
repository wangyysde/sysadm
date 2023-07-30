
function addObjvalidTextValue(actionUri,obj){
    var uri = actionUri;
    if(!obj) {
        return displayTip("操作出错","RED",0);
    }
    var objValue = obj.value;
    uri = uri + "?objvalue=" + objValue;
    var respValue = doAjax(uri,"POST",false);
    if(respValue["error"]){
        obj.focus();
        return displayTip(respValue["msg"],"RED",0);
    }

    return;
}

function addObjSelectChanged(actionUri,subObjId,objId) {
    actionUri = actionUri.replace(/^\s+|\s+$/g," ");
    subObjId = subObjId.replace(/^\s+|\s+$/g," ");
    objId = objId.replace(/^\s+|\s+$/g," ");
    if(actionUri == "" || subObjId == "" || objId == ""){
        return;
    }

    var uri = pageUrl + actionUri;
    uri = uri + "?objID=" + objId;
    var respValue = doAjax(uri,"GET",false);
    if(respValue["error"]){
        return displayTip(respValue["msg"],"RED",0);
    }

    var subObjDataStr = respValue["msg"];
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

    return;
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

function doAjax(actionUrl,actionType,isMultiPart) {
    var data;
    if(isMultiPart){
        data = new FormData($("#addObjectForm")[0]);
    } else {
        data = "";
    }
    var ajaxRet = {"msg": "未知错误aaa", "error": true};
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

function addObjClickFileButton(objectID){
    var obj = document.getElementById(objectID);
    obj.click();
}

function addObjChangeFilevalue(divID,inputID){
    var inputObj = document.getElementById(inputID);
    var divObj = document.getElementById(divID);
    var objValue = inputObj.value;
    var pos = objValue.lastIndexOf("\\");
    divObj.innerHTML = objValue.substring(pos+1);
}

function addObjSubmit(postUri, redirectUri) {
    var respValue = doAjax(postUri,"POST",true);
    if(respValue["error"]){
        alert(respValue["msg"]);
        return displayTip(respValue["msg"],"RED",0);
    }

    alert(respValue["msg"]);
    var url = pageUrl + redirectUri;
    $('#container').load(url);
    displayTip(respValue["msg"],"",0);

}

function addObjCancel(cancelRedirect) {
    if (confirm("确认取消本次集群的添加？") == true){
        var url = pageUrl + cancelRedirect;
        $('#container').load(url);
    }

    return;
}