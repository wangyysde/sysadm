function nsAddLabel(formID,lineID,itemID,uri,obj){
    var labelObjLine  = document.getElementById("linenewNsLabel");
    var labelContainer = document.getElementById("containernewNsLabel");
    var newEle = document.createElement("div");
    newEle.className = labelObjLine.className;
    newEle.innerHTML = labelObjLine.innerHTML;
    labelContainer.appendChild(newEle);
    return;
}


function nsDelLabel(formID,lineID,itemID,uri,obj){
    var labelContainer = document.getElementById("containernewNsLabel");
    var spanElement = obj.parentNode;
    var lineElement = spanElement.parentNode;
    labelContainer.removeChild(lineElement);

    return;
}

function addNewNamespace(formID,module,dcID,clusterID,namespace,actionType){
    var addTypeObj = document.getElementById("addType");
    var addTypeValue = addTypeObj.value;
    if(addTypeValue == "0"){
        var editor = ace.edit("cardeditor");
        var code = editor.getValue();
        var objContent = document.getElementById("objContent");
        objContent.value = code;
        if(code.length < 10 ){
            return formDataShowTip("yaml内容长度不合法","RED",0);
        }
    }

    if(addTypeValue == "1"){
        var objFile = document.getElementById("objFile");
        var objFileValue = objFile.value;
        if(objFileValue.length < 3 ){
            return formDataShowTip("未选择上传的文件","RED",0);
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

function addNewObjFunc(formID,module,dcID,clusterID,namespace,actionType,postUri){
    var addTypeObj = document.getElementById("addType");
    var addTypeValue = addTypeObj.value;
    if(addTypeValue == "0"){
        var editor = ace.edit("cardeditor");
        var code = editor.getValue();
        var objContent = document.getElementById("objContent");
        objContent.value = code;
        if(code.length < 10 ){
            return formDataShowTip("yaml内容长度不合法","RED",0);
        }
    }

    if(addTypeValue == "1"){
        var objFile = document.getElementById("objFile");
        var objFileValue = objFile.value;
        if(objFileValue.length < 3 ){
            return formDataShowTip("未选择上传的文件","RED",0);
        }
    }

    postUri = postUri + "?dcID=" + dcID + "&clusterID=" + clusterID + "&namespace=" + namespace;
    var respValue = formAjax(formID,postUri,actionType,true);

    $('#container').load(lastUrl);

    if(respValue["error"]){
        return formDataShowTip(respValue["msg"],"RED",0);
    }

    return formDataShowTip(respValue["msg"],"GREEN",0);
}

function addNewQuotaFunc(formID,module,dcID,clusterID,namespace,actionType) {
    var postUri = pageUrl + module + "/addNewQuota";

    return addNewObjFunc(formID,module,dcID,clusterID,namespace,actionType,postUri);
}

function addNewLimitRangeFunc(formID,module,dcID,clusterID,namespace,actionType){
    var postUri = pageUrl + module + "/addNewLimitRange";

    return addNewObjFunc(formID,module,dcID,clusterID,namespace,actionType,postUri);
}