
function displayAddObjectForm(module){
    var url = pageUrl + module + "/addform" ;
    $('#container').load(url);
}

function changeCard(cardNo){
    var cardHead = document.getElementsByClassName("cardheadline");
    var headSpans = cardHead[0].getElementsByTagName("span");
    var cardcontentdiv = document.getElementsByName("cardcontentdiv");

    for(i = 0; i < headSpans.length; i++ ){
        headSpans[i].className = "";
        cardcontentdiv[i].style.display = "none";
        if(i == cardNo){
            headSpans[i].className = "activecard";
            cardcontentdiv[i].style.display = "block";
            var addType = document.getElementById("addType");
            addType.value = cardNo;
        }
    }
    return;
}

function doAjax(actionUrl,actionType,isMultiPart) {
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

function addObjvalidTextValue(module,action,obj){
    var postUri = pageUrl + "/validate/" + module + "/" + action;
    if(!obj) {
        return displayTip("操作出错","RED",0);
    }
    var objValue = obj.value;
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;
    uri = postUri + "?objvalue=" + objValue +"&dcID=" + dcID + "&clusterID=" + clusterID + "&namespace=" + namespace;
    var respValue = doAjax(uri,"POST",false);
    if(respValue["error"]){
        obj.focus();
        return displayTip(respValue["msg"],"RED",0);
    }

    return;
}

function addObjSelectChanged(module,action,subObjId,objValue) {
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;
    action = action.replace(/^\s+|\s+$/g," ");
    subObjId = subObjId.replace(/^\s+|\s+$/g," ");
    objValue = objValue.replace(/^\s+|\s+$/g," ");
    if(action == "" || subObjId == "" || objValue == ""){
        return;
    }

    var postUri = pageUrl + "/validate/" + module + "/" + action;
    var uri = postUri + "?objvalue=" + objValue +"&dcID=" + dcID + "&clusterID=" + clusterID + "&namespace=" + namespace;
    var respValue = doAjax(uri,"POST",false);
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

function addObjRadioClick(module,action,actionKind,objID,relatedIsDisplay, subObjID,radioOption){
    // actionKind 所定义的函定义参见GO语言编写的objectsUI包中 objectsUIDefined.go中的定义
    if(actionKind == "1"){
        subObjStr = "span" + subObjID;
        var obj = document.getElementById(subObjStr);
        if(obj){
            if(relatedIsDisplay){
                obj.style.display = "block";
            } else {
                obj.style.display = "none";
            }
            return;
        }
        alert("批定的子对象" + subObjID + "不存在");
        return;
    }

    if(actionKind == "2"){
        if(typeof  addObjRadioClick == "function"){
            return addObjRadioClick(module,action,objID, relatedIsDisplay,subObjID,radioOption);
        }
        alert("JS的动作类型为自定义函数，但是页面内没有定义addObjRadioClick函数");
        return ;
    }
    alert("所指定的动作类型不正确");
    return;
}