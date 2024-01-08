
function formDataShowTip(msg,color,timeout){
    var  colorStr = "#1f6f4a";
    if(color.toUpperCase() == "RED"){
        colorStr = "#ff0000";
    }

    if(msg == "" ){
        msg = "发生了未知错误，请稍后再试或联系系统管理员";
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

function formAjax(formID, actionUrl,actionType,isMultiPart) {
    var data;
    if(isMultiPart){
        data = new FormData($("#"+formID)[0]);
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
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            var msg = "提交错误";
            if(textStatus == "error"){
                msg = msg + ":" + "服务器端处理错误";
            }
            if(textStatus == "timeout") {
                msg = msg + ":" + "网络超时";
            }
            if(textStatus == "parseerror"){
                msg = msg + ":" + "数据解析错误";
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
                    msg = formDatab64DecodeUnicode(msgLine["msg"]);

                }
                ajaxRet.msg = msg;
                ajaxRet.error = false;
            } else {
                var messageArray = result.message;
                var msg = "未知错误";
                if (messageArray[0]) {
                    var msgLine = messageArray[0];
                    msg = formDatab64DecodeUnicode(msgLine["msg"]);
                }
                ajaxRet.msg = msg;
                ajaxRet.error = true;
            }
        }
    });

    return ajaxRet;
}

// 将经过base64编码的utf-8字符串解码为普通字符串，支持中文
function formDatab64DecodeUnicode(str) {
    return decodeURIComponent(Array.prototype.map.call(atob(str), function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
}

function formDataTextInputValueChange(formID,formMethod,module, uri,obj,fn){
    formID = formID.replace(/^\s+|\s+$/g," ");
    uri = uri.replace(/^\s+|\s+$/g," ");
    fn = fn.replace(/^\s+|\s+$/g," ");
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;

    if(fn.length > 1 && typeof window[fn] === "function"){
        return window[fn](formID,module,dcID,clusterID,namespace,uri,obj);
    }

    if(!obj) {
        return formDataShowTip("操作出错","RED",0);
    }

    var objValue = obj.value;
    var requestUri = pageUrl + module + "/";
    if(uri.length > 0){
        requestUri = requestUri + uri + "/";
    }

    requestUri = requestUri + "?dcID=" + dcID + "&clusterID=" + clusterID + "&namespace=" + namespace;
    requestUri = requestUri + "&objValue=" + objValue;

    var respValue = formAjax(formID,requestUri,formMethod,false);
    if(respValue["error"]){
        obj.focus();
        return formDataShowTip(respValue["msg"],"RED",0);
    }

    return;
}

function formDataSelectInputValueChange(formID,module, uri,obj,fn){
    formID = formID.replace(/^\s+|\s+$/g," ");
    uri = uri.replace(/^\s+|\s+$/g," ");
    obj = obj.replace(/^\s+|\s+$/g," ");
    fn = fn.replace(/^\s+|\s+$/g," ");
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;

    if(fn.length > 1 && typeof window[fn] === "function"){
        return window[fn](formID,module,dcID,clusterID,namespace,uri,obj);
    }

    if(!obj) {
        return formDataShowTip("操作出错","RED",0);
    }

    var requestUri = pageUrl + "/";
    if(uri.length > 0){
        requestUri = requestUri + uri + "/";
    }

    requestUri = requestUri + "?dcID=" + dcID + "&clusterID=" + clusterID + "&namespace=" + namespace;
    requestUri = requestUri + "&objvalue=" + objValue;
    $('#container').load(requestUri);

    return
}

function formDataCheckBoxClick(formID,module,groupID,option,fn) {
    formID = formID.replace(/^\s+|\s+$/g, " ");
    groupID = groupID.replace(/^\s+|\s+$/g, " ");
    option = option.replace(/^\s+|\s+$/g, " ");
    fn = fn.replace(/^\s+|\s+$/g, " ");
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;

    if (fn.length > 1 && typeof window[fn] === "function") {
        return window[fn](formID,module,dcID,clusterID,namespace,groupID,option);
    }

    alert("JS函数" + fn + "未定义");
    return;
}

function formDataRadioClick(formID,module,groupID,option,fn) {
    formID = formID.replace(/^\s+|\s+$/g, " ");
    groupID = groupID.replace(/^\s+|\s+$/g, " ");
    fn = fn.replace(/^\s+|\s+$/g, " ");
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;

    if (fn.length > 1 && typeof window[fn] === "function") {
        return window[fn](formID,module,dcID,clusterID,namespace,groupID,option);
    }

    alert("JS函数" + fn + "未定义");
    return;
}

function formDataFileInputValueChange(divID,inputID){
    var inputObj = document.getElementById(inputID);
    var divObj = document.getElementById(divID);
    var objValue = inputObj.value;
    var pos = objValue.lastIndexOf("\\");
    divObj.innerHTML = objValue.substring(pos+1);
}

function formDataFileInputButtonClick(objectID){
    var obj = document.getElementById(objectID);
    obj.click();
}

function formDataTextareaInputValueChange(formID, module, uri,obj,fn){
    formID = formID.replace(/^\s+|\s+$/g," ");
    uri = uri.replace(/^\s+|\s+$/g," ");
    obj = obj.replace(/^\s+|\s+$/g," ");
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;

    if(fn.length > 1 && typeof window[fn] === "function"){
        return window[fn](formID,module,dcID,clusterID,namespace,uri,obj);
    }

    if(!obj) {
        return formDataShowTip("操作出错","RED",0);
    }

    var objValue = obj.Value;
    var requestUri = pageUrl + "/";
    if(uri.length > 0){
        requestUri = requestUri + uri + "/";
    }

    var formObj = document.getElementById(formID);
    if(!formObj) {
        return formDataShowTip("操作出错","RED",0);
    }
    var formMethod = formObj.method;
    var enctype = formObj.enctype;
    var isMultiPart = false;
    if(enctype == 'multipart/form-data'){
        isMultiPart = true;
    }

    requestUri = requestUri + "?dcID=" + dcID + "&clusterID=" + clusterID + "&namespace=" + namespace;
    requestUri = requestUri + "&objvalue=" + objValue;

    var respValue = formAjax(formID,requestUri,formMethod,isMultiPart);
    if(respValue["error"]){
        obj.focus();
        return formDataShowTip(respValue["msg"],"RED",0);
    }

    return;

}

function formDataWordsInputClick(formID,module,lineID,itemID,uri,fn,obj){
    formID = formID.replace(/^\s+|\s+$/g," ");
    lineID = lineID.replace(/^\s+|\s+$/g," ");
    itemID = itemID.replace(/^\s+|\s+$/g," ");
    uri = uri.replace(/^\s+|\s+$/g," ");
    fn = fn.replace(/^\s+|\s+$/g," ");
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;


    if(fn.length > 1 && typeof window[fn] === "function"){
        return window[fn](formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj);
    }

    alert("函数" + fn + "没有定义");

    return;
}

function forumDataSubmitForm(formID,module,actionType,fn){
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;

    if(fn.length > 1 ){
        if(typeof window[fn] === "function") {
            return window[fn](formID, module, dcID, clusterID, namespace, actionType);
        } else {
            return formDataShowTip("函数" + fn +"未被定义","RED",0);
        }
    }

    var postUri = pageUrl + module + "/add";
    postUri = postUri + "?dcID=" + dcID + "&clusterID=" + clusterID + "&namespace=" + namespace;
    var respValue = formAjax(formID,postUri,actionType,true);

    $('#container').load(lastUrl);

    if(respValue["error"]){
        return formDataShowTip(respValue["msg"],"RED",0);
    }

    return formDataShowTip("内容已经添加成功","GREEN",0);
}

function forumDataCancelForm(formID,module,actionType,fn){
    var dcID = document.getElementById("dcID").value;
    var clusterID = document.getElementById("clusterID").value;
    var namespace = document.getElementById("namespace").value;

    if(fn.length > 1 && typeof window[fn] === "function"){
        return window[fn](formID,module,dcID,clusterID,namespace,actionType);
    }

    $('#container').load(lastUrl);
    return;
}