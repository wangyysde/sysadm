// calling this function when check or uncheck passiveMode checkbox on add host form page
function changeAgentMode(){
    var passiveMode = document.getElementById("passiveMode");
    var divAgentPort = document.getElementById("divAgentPort");
    var divCommandUri = document.getElementById("divCommandUri");
    var agentIsTls = document.getElementById("agentIsTls");
    var divisTLS = document.getElementById("divisTLS");
    var divinsecureSkipVerify = document.getElementById("divinsecureSkipVerify");
    var divcommandStatusUri = document.getElementById("divcommandStatusUri");
    var divcommandLogsUri = document.getElementById("divcommandLogsUri");
    var caInput = document.getElementById("caInput");
    var certInput = document.getElementById("certInput");
    var keyInput = document.getElementById("keyInput");
    if(passiveMode.checked){
        divAgentPort.style.display="none";
        divCommandUri.style.display = "none";
        divisTLS.style.display = "none";
        agentIsTls.checked = false;
        divinsecureSkipVerify.style.display = "none";
        divcommandStatusUri.style.display = "none";
        divcommandLogsUri.style.display = "none";
        caInput.style.display = "none";
        certInput.style.display = "none";
        keyInput.style.display = "none";
    }else {
        divAgentPort.style.display="block";
        divCommandUri.style.display = "block";
        divisTLS.style.display = "inline";
        divinsecureSkipVerify.style.display = "inline";
        divcommandStatusUri.style.display = "block";
        divcommandLogsUri.style.display = "block";

    }
}

function enableISTLS(){
    var caInput = document.getElementById("caInput");
    var certInput = document.getElementById("certInput");
    var keyInput = document.getElementById("keyInput");
    var insecureSkipVerify = document.getElementById("insecureSkipVerify");
    var divinsecureSkipVerify = document.getElementById("divinsecureSkipVerify");
    var agentIsTls = document.getElementById("agentIsTls");

    if (agentIsTls.checked){
        caInput.style.display = "block";
        certInput.style.display = "block";
        keyInput.style.display = "block";
        divinsecureSkipVerify.style.display = "inline";
        insecureSkipVerify.Checked = true;
        return true;
    }

    caInput.style.display = "none";
    certInput.style.display = "none";
    keyInput.style.display = "none";
    divinsecureSkipVerify.style.display = "none";
    insecureSkipVerify.Checked = false;

    return ;
}

function infrastructurlistcheckData(){

    var hostname = document.getElementById("hostname");
    if(hostname.length < 1 || hostname.length >255){
        alert("you input hostname " + hostname + "is not valid.");
        return false;
    }

    var ipAdd = document.getElementById("ip").value;
    if(ipAdd.length < 7){
        alert("You input ip address " + ipAdd + "is not valid");
        return false;
    }

    // checking agent parameters if passive mode is not choiced
    var passiveMode = document.getElementById("passiveMode");
    if(passiveMode.checked){
        var commandUri = document.getElementById("commandUri");
        commandUri.value = "";

        var commandStatusUri = document.getElementById("commandStatusUri");
        commandStatusUri.value = "";

        var commandLogsUri = document.getElementById("commandLogsUri");
        commandLogsUri.value = "";

        var agentPort = document.getElementById("agentPort");
        agentPort.value = "0";

        infrastructureResetTLS();

    }
    else {
        var commandUri = document.getElementById("commandUri").value;
        if (commandUri.length < 1){
            if(!window.confirm("你没有输入指令接收地址，确认使用默认的指令接收地址？")){
                return false;
            }
        }

        var commandStatusUri = document.getElementById("commandStatusUri").value;
        if (commandStatusUri.length < 1){
            if(!window.confirm("指令状态查询地址为空，确认使用默认的指令接收地址？")){
                return false;
            }
        }

        var commandLogsUri = document.getElementById("commandLogsUri").value;
        if (commandLogsUri.length < 1){
            if(!window.confirm("指令日志查询地址为空，确认使用默认的指令接收地址？")){
                return false;
            }
        }

        var agentIsTls = document.getElementById("agentIsTls");
        if (agentIsTls.checked) {
            if(!infrastructureCheckTLS()){
                return false;
            }
        } else {
            infrastructureResetTLS();
        }

        var agentPort = document.getElementById("agentPort");
        var portValue = agentPort.value;
        if (!isPort(portValue)) {
            alert("Agent端口不合法");
            return false;
        }

    }

    // checking OS distribution
    var osID = document.getElementById("osID");
    if (osID[osID.selectedIndex].value == 0){
        alert("Please select the OS for this host");
        return false;
    }

    // checking OS version
    var osVersionSelected = document.getElementById("osversionid");
    var index = osVersionSelected.selectedIndex;
    var osVersionID = osVersionSelected.options[index].value;
    if(osVersionID == 0){
        alert("Please select OS and its version");
        return false;
    }

    return true;
}

function infrastructureCheckTLS(){
    var agentCa = document.getElementById("agentCa").value;
    if (agentCa.length < 1){
        if(!window.confirm("未选择根证书，确认使用系统根证书？")){
            return false;
        }
    }

    var agentCert = document.getElementById("agentCert").value;
    if (agentCert.length < 1){
        alert("启用TLS连接时，需要上传所需要的证书文件");
        return false;
    }

    var agentKey = document.getElementById("agentKey").value;
    if (agentKey.length < 1){
        alert("启用TLS连接时，需要上传所需要的密钥文件");
        return false;
    }

    return true;
}

function infrastructureResetTLS(){
    var agentCa = document.getElementById("agentCa");
    agentCa.value = "";

    var agentCert = document.getElementById("agentCert");
    agentCert.value = "";

    var agentKey = document.getElementById("agentKey");
    agentKey.value = "";

    var insecureSkipVerify = document.getElementById("insecureSkipVerify");
    insecureSkipVerify.checked = false;

    return true;
}

// display add host form page
function infrastructureAddHost() {
    var ownerid = document.getElementById("userid");
    if (ownerid.value == "") {
        alert("添加节点需要登录到系统!");
        return
    }
    var maskLayer = document.getElementById("maskLayer");
    maskLayer.style.zIndex = 2000
    var addProjectForm = document.getElementById("addform");
    addProjectForm.style.display = "block";
    addProjectForm.style.zIndex = 2100;
    addProjectForm.focus();
}

function submitAddHost() {
    if(!infrastructurlistcheckData()){
        return false;
    }

    // var data = new FormData($('#infrastructureAddHost'));
    $.ajax({
        type: "POST",
    //    dataType: "json",
    //    crossDomain: true,
        contentType: false,
        cache: false,
        processData: false,
    //    mimeType: "multipart/form-data",
        url: "/api/1.0/infrastructure/add",
        data: new FormData($('#infrastructureAddHostForm')[0]),
    //    data: $('#infrastructureAddHostForm').serialize(), // 你的formid
        //async: false,
        error: function(xmlObj, request) {
            var errMsg = "";
            if (xmlObj.responseText == null) {
                errMsg = "Connection error";
            } else {
                errMsg = xmlObj.responseText;
            }
            cancelAddHost();
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
            cancelAddHost();
            if (result.errorCode == 0) {
                var messageArray = result.message;
                var msg = "未知错误";
                if (messageArray[0]) {
                    var msgLine = messageArray[0];
                    msg = window.atob(msgLine["msg"]);

                }
                tip.style.display = "block";
                tip.style.backgroundColor = "#1f6f4a";
                tip.style.display = "block";
                tip.innerHTML = msg;
                refreshPage();
                setTimeout(function() {
                    var tip = document.getElementById("tip");
                    tip.style.display = "none";
                }, 5000);

            } else {
                var messageArray = result.message;
                var msg = "未知错误";
                if (messageArray[0]) {
                    var msgLine = messageArray[0];
                    msg = window.atob(msgLine["msg"]);

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

// check or uncheck all yum checkbox on add host form page
function selectAllYumCheckbox(isChecked) {
    if(isChecked) {
        var chklist = document.getElementsByName('yumid[]');
        for (var i = 0; i < chklist.length; i++) {
            chklist[i].checked = true;
        }
    } else {
        var chklist = document.getElementsByName('yumid[]');
        for (var i = 0; i < chklist.length; i++) {
            chklist[i].checked = false;
        }
    }
}

function clickagentCaButton(objectID){
    var agentCa = document.getElementById(objectID);
    agentCa.click();
}

function changefilevalue(divID,inputID){
    var inputObj = document.getElementById(inputID);
    var divObj = document.getElementById(divID);
    var objValue = inputObj.value;
    var pos = objValue.lastIndexOf("\\");
    divObj.innerHTML = objValue.substring(pos+1);
}

function cancelAddHost() {
    var addProjectForm = document.getElementById("addform");
    addProjectForm.style.display = "none";
    addProjectForm.style.zIndex = 0;
    var maskLayer = document.getElementById("maskLayer");
    maskLayer.style.zIndex = 0
}

function changePage(urlparas) {
    var url = "/infrastructure/list" + urlparas;
    $('#container').load(url);
}

function projectChanged(projectid) {
    var urlParam = "?projectid="+projectid;
    changePage(urlParam);
}

function doSearch(searchKey) {
    var urlParam = "?searchKey=" + searchKey;
    changePage(urlParam);
}

// check or uncheck all host checkbox on host list page
function selectAllHostCheckbox(isChecked) {
    if(isChecked) {
        var chklist = document.getElementsByName('hostid[]');
        for (var i = 0; i < chklist.length; i++) {
            chklist[i].checked = true;
        }
        var buttonDel = document.getElementById("buttonDel");
        buttonDel.style.color = "#ffffff";
        buttonDel.style.background = "#3c8dbc";
        buttonDel.disabled = false;
    } else {
        var chklist = document.getElementsByName('hostid[]');
        for (var i = 0; i < chklist.length; i++) {
            chklist[i].checked = false;
        }
        var buttonDel = document.getElementById("buttonDel")
        buttonDel.style.color = "#3c8dbc";
        buttonDel.style.background = "#cbd8df";
        buttonDel.disabled = true;
    }
}

function selectHostCheckbox(isChecked){
    if(isChecked){
        var buttonDel = document.getElementById("buttonDel")
        buttonDel.style.color = "#ffffff";
        buttonDel.style.background = "#3c8dbc";
        buttonDel.disabled = false;

        var chklist = document.getElementsByName('hostid[]');
        var allChecked = true;
        for (var i = 0; i < chklist.length; i++) {
            if(!chklist[i].checked){
                allChecked = false;
            }
        }
        if(allChecked){
            var hostth = document.getElementById("hostidth");
            hostth.checked = true;
        }
        return;
    }

    var allNotChecked = true;
    var chklist = document.getElementsByName('hostid[]');
    for (var i = 0; i < chklist.length; i++) {
        if(chklist[i].checked){
            allNotChecked = false;
        }
    }

    if(allNotChecked){
        var buttonDel = document.getElementById("buttonDel")
        buttonDel.style.color = "#3c8dbc";
        buttonDel.style.background = "#cbd8df";
        buttonDel.disabled = true;
        var hostth = document.getElementById("hostidth");
        hostth.checked = false;
    }
}

function delHostJs() {

    var chklist = document.getElementsByName('hostid[]');
    var checkedItem = 0;
    for (var i = 0; i < chklist.length; i++) {
        if (chklist[i].checked) {
            checkedItem = checkedItem + 1;
        }
    }

    if (checkedItem == 0) {
        alert("没有主机信息可以删除!");
        return
    }
    var ok = confirm("确认需要删除这些主机信息吗？");
    if (!ok) {
        return
    }
    $.ajax({
        statusCode: {
            500: function () {
                //	refreshPage();
                var tip = document.getElementById("tip");
                tip.innerHTML = "出现服务器端错误，请稍后再试";
                tip.style.display = "block";
                setTimeout(function () {
                    var tip = document.getElementById("tip");
                    tip.style.display = "none";
                }, 5000);
            },
            501: function () {
                //	refreshPage();
                var tip = document.getElementById("tip");
                tip.innerHTML = "出现服务器端错误，请稍后再试";
                tip.style.display = "block";
                setTimeout(function () {
                    var tip = document.getElementById("tip");
                    tip.style.display = "none";
                }, 5000);
            },
            502: function() {
                //	refreshPage();
                var tip = document.getElementById("tip");
                tip.innerHTML = "出现服务器端错误，请稍后再试";
                tip.style.display = "block";
                setTimeout(function () {
                    var tip = document.getElementById("tip");
                    tip.style.display = "none";
                }, 5000);
            },

        },
        type: "POST",
        dataType: "json",
        url: "/infrastructure/delhost",
        data: $('#hostList').serialize(), // 你的formid
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
            if (result.status) {
                tip.style.display = "block";
                tip.style.backgroundColor = "#1f6f4a";
                tip.style.display = "block";
                var messageArray = result.message;
                var msg = "未知错误";
                if (messageArray[0]) {
                    var msgLine = messageArray[0];
                    msg = window.atob(msgLine["msg"]);

                }
                tip.innerHTML = msg;

                refreshPage();
                setTimeout(function() {
                    var tip = document.getElementById("tip");
                    tip.style.display = "none";
                }, 5000);
            } else {
                var messageArray = result.message;
                var msg = "未知错误";
                if (messageArray[0]) {
                    var msgLine = messageArray[0];
                    msg = window.atob(msgLine["msg"]);

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

function displayHostDetails(hostid){
    var url = "/infrastructure/hostdetails?hostid=" + hostid;
    var detailshostid = document.getElementById("detailshostid");
    detailshostid.value = hostid;

    $('#detailHost').load(url);

}

function closeDetailsPage() {
    var detailHost = document.getElementById("detailHost");
    detailHost.style.display = "none";
    detailHost.style.zIndex = 0;
    var maskLayer = document.getElementById("maskLayer");
    maskLayer.style.zIndex = 0
}