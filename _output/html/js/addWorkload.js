
function addWorkloadChangeType(cardNo){
    var cardHead = document.getElementsByClassName("addWorkloadNavigationBlock");
    var headSpans = cardHead[0].getElementsByTagName("span");
    var cardcontentdiv = document.getElementsByName("addWorkloadFormTypeContent");

    for(i = 0; i < headSpans.length; i++ ){
        headSpans[i].className = "";
        cardcontentdiv[i].style.display = "none";
        if(i == cardNo){
            headSpans[i].className = "addWorkloadNavigationActiveCard";
            cardcontentdiv[i].style.display = "block";
            var addType = document.getElementById("addType");
            addType.value = cardNo;
        }
    }
    return;
}

function addWorkloadChangeBlockCatelory(cardNo){
    var cardHeads = document.getElementsByName("addWorkloadFormPartCardHead");
    var cardcontents = document.getElementsByName("addWorkloadFormPartContent");

    for(i = 0; i < cardHeads.length; i++ ){
        cardHeads[i].className = "";
        cardcontents[i].style.display = "none";
        if(i == cardNo){
            cardHeads[i].className = "formPartCardHeadSpansActived";
            cardcontents[i].style.display = "block";
        }
    }
    return;
}

function addNewContainerFormData(){
    // 基本信息部分
    var containerNameObj = document.getElementById("containerNameID");
    var containerName = containerNameObj.value;
    if(!validateContainerName(containerName)){
        containerNameObj.focus();
        return;
    }

    var containerTypeObjs = document.getElementsByName("containerType");
    var containerType = "";
    for(i = 0; i < containerTypeObjs.length; i++){
        if(containerTypeObjs[i].checked){
            var containerTypeTemp = containerTypeObjs[i].value;
            if(containerType == "") {
                containerType = containerTypeTemp;
            } else {
                containerType = containerType + "," + containerTypeTemp;
            }
        }
    }
    if(containerType == ""){
        containerType = "0";
    }

    var containerImageObj = document.getElementById("containerImageID");
    var containerImage = containerImageObj.value;
    if(!validateImagePath(containerImage)){
        containerImageObj.focus();
        return;
    }

    var imagePullPolicyObj = document.getElementById("imagePullPolicyID");
    var imagePullPolicy = imagePullPolicyObj.options[imagePullPolicyObj.selectedIndex].value;

    var startCommandObj = document.getElementById("startCommandID");
    var startCommand = startCommandObj.value;

    var environmentObj = document.getElementById("environmentID");
    var environment = environmentObj.value;
    var basicInfo = createContainerBasicInfo(containerName,containerType,containerImage,imagePullPolicy,startCommand,environment);

    // Quota部分
    var cpuRequestObj = document.getElementById("cpuRequestID");
    var cpuRequest = cpuRequestObj.value;

    var cpuLimitObj = document.getElementById("cpuLimitID");
    var cpuLimit = cpuLimitObj.value;

    var memRequestObj = document.getElementById("memRequestID");
    var memRequest = memRequestObj.value;

    var memLimitObj = document.getElementById("memLimitID");
    var memLimit = memLimitObj.value;
    var quota = createContainerQuota(cpuRequest,cpuLimit,memRequest,memLimit);

    // 健康探测部分
    var startup = setHealthyCheck("startup");
    var readiness = setHealthyCheck("readiness");
    var liveness= setHealthyCheck("liveness");
    if(startup == null || readiness == null || liveness == null){
        return;
    }

    var data = createContainerData(basicInfo,quota,startup,readiness,liveness)
    var jsonData = JSON.stringify(data);
    var encodeData = Base64.encode(jsonData);

    addContainerToContainerList(containerType,containerName,encodeData);

    resetFormInputItemValue();

    return;
}

function validateContainerName(containerName){
    var containerNameObjs = document.getElementsByName("containerName[]");
    for(i = 0; i < containerNameObjs.length; i++){
        var existContainerName = containerNameObjs[i].value;
        if(containerName == existContainerName){
            formDataShowTip(("容器名" + containerName + "已经存在"),"READ",0);
            return false;
        }
    }

    const dnsLabel = /[a-z0-9]([-a-z0-9]*[a-z0-9]){1,63}?/;
    if(dnsLabel.test(containerName)){
        return true;
    }
    formDataShowTip(("容器名" + containerName + "不合法"),"RED",0);
    return false;
}

function validateImagePath(imagePath){
    const imagePathreg = /[a-zA-Z]([-a-zA-Z:0-9]){1,128}?/;
    if(imagePathreg.test(imagePath)){
        return true;
    }
    formDataShowTip(("容器" + imagePath + "地址不合法"),"RED",0);
    return false;
}

function  setHealthyCheck(kind){
    var probeTypeObj = document.getElementsByName((kind + "ProbeType"));
    var probeType = "0";
    for(i = 0; i < probeTypeObj.length; i++){
        if(probeTypeObj[i].checked){
            probeType = probeTypeObj[i].value;
            break;
        }
    }

    var httpProtocolObj = document.getElementById((kind + "HttpProtocolID"));
    var httpProtocol = httpProtocolObj.options[httpProtocolObj.selectedIndex].value;

    var httpPortObj = document.getElementById((kind + "HttpPortID"));
    var httpPort = httpPortObj.value;

    var httpPathObj = document.getElementById((kind + "HttpPathID"));
    var httpPath = httpPathObj.value;

    var headerKeyObj = document.getElementsByName((kind + "HttpHeaderKey"));
    var headerValueObj = document.getElementsByName((kind + "HttpHeaderValue"));
    var httpHeaderKey = new Array();
    var httpHeaderValue = new Array();
    for(i = 0; i < headerKeyObj.length; i++){
        if(headerKeyObj[i].value == "" && headerValueObj[i].value != "" ){
            formDataShowTip(("值为" + headerValueObj[i].value + "HTTP Header的Key不能为空"),"RED",0);
            return null;
        }
        if(headerKeyObj[i].value != "" && headerValueObj[i].value == "" ){
            formDataShowTip((headerKeyObj[i].value + "项HTTP Header的值不能为空"),"RED",0);
            return null;
        }
        httpHeaderKey.push(headerKeyObj[i].value);
        httpHeaderValue.push(headerValueObj[i].value);
    }

    var tcpPortObj = document.getElementById((kind + "TcpPortID"));
    var tcpPort = tcpPortObj.value;

    var commandObj = document.getElementById((kind + "CommandID"));
    var command = commandObj.value;

    var initialDelaySecondsObj = document.getElementById((kind + "InitialDelaySecondsID"));
    var initialDelaySeconds = initialDelaySecondsObj.value;

    var periodSecondsObj = document.getElementById((kind + "PeriodSecondsID"));
    var periodSeconds = periodSecondsObj.value;

    var timeoutSecondsObj = document.getElementById((kind + "TimeoutSecondsID"));
    var timeoutSeconds = timeoutSecondsObj.value;

    var failureThresholdObj = document.getElementById((kind + "FailureThresholdID"));
    var failureThreshold = failureThresholdObj.value;

    var successThresholdObj = document.getElementById((kind + "SuccessThresholdID"));
    var successThreshold = successThresholdObj.value;

    var probeData = createContainerHealthyCheck(probeType,httpProtocol,httpPort,httpPath,httpHeaderKey,httpHeaderValue,tcpPort,command,initialDelaySeconds,
        periodSeconds,timeoutSeconds,failureThreshold,successThreshold)

    return probeData;
}

function resetFormInputItemValue(){
    var containerNameObj = document.getElementById("containerNameID");
    containerNameObj.value = "";

    var containerTypeObj = document.getElementsByName("containerType");
    for(i = 0; i < containerTypeObj.length; i++){
        containerTypeObj[i].checked = false;
    }

    var containerImageObj = document.getElementById("containerImageID");
    containerImageObj.value = "";

    var imagePullPolicyObj = document.getElementById("imagePullPolicyID");
    for(i = 0; i < imagePullPolicyObj.options.length; i++){
        if(imagePullPolicyObj.options[i].value == "IfNotPresent"){
            imagePullPolicyObj.options[i].selected = true;
            break;
        }
    }

    var startCommandObj = document.getElementById("startCommandID");
    startCommandObj.value = "";

    var environmentObjs = document.getElementById("environmentID");
    environmentObjs.value = "";

    var cpuRequestObj = document.getElementById("cpuRequestID");
    cpuRequestObj.value = "";

    var cpuLimitObj = document.getElementById("cpuLimitID");
    cpuLimitObj.value = "";

    var memRequestObj = document.getElementById("memRequestID");
    memRequestObj.value = "";

    var memLimitObj = document.getElementById("memLimitID");
    memLimitObj.value = "";

    resetFormItemValueForHealthyCheck("startup");
    resetFormItemValueForHealthyCheck("readiness");
    resetFormItemValueForHealthyCheck("liveness");

    return;
}

function resetFormItemValueForHealthyCheck(kind){
    var probeTypeObj = document.getElementsByName((kind + "ProbeType"));
    for(i = 0; i < probeTypeObj.length; i++){
        probeTypeObj[i].checked = false;
        if(probeTypeObj[i].value == "0"){
            probeTypeObj[i].checked = true;
        }
    }

    var httpProtocolObj = document.getElementById((kind + "HttpProtocolID"));
    for(i = 0; i < httpProtocolObj.options.length; i++ ){
        httpProtocolObj.options[i].selected = false;
        if(httpProtocolObj.options[i].value =="0"){
            httpProtocolObj.options[i].selected = true;
        }
    }
    var httpPortObj = document.getElementById((kind + "HttpPortID"));
    httpPortObj.value = "";

    var httpPathObj = document.getElementById((kind + "HttpPathID"));
    httpPathObj.value = "";

    var headerLineObj = document.getElementById(("container" + kind + "HttpHeaderLineID"));
    for(i = 0; i < (headerLineObj.childNodes.length -1 ); i++ ){
        headerLineObj.removeChild(headerLineObj.childNodes[i]);
    }

    var httpHeaderKeyObj = document.getElementsByName((kind + "HttpHeaderKey"));
    for(i = 0; i < httpHeaderKeyObj.length; i++){
        httpHeaderKeyObj[i].value = "";
    }

    var httpHeaderValueObj = document.getElementsByName((kind + "HttpHeaderValue"));
    for(i = 0; i < httpHeaderValueObj.length; i++ ){
        httpHeaderValueObj[i].value = "";
    }

    var tcpPortObj = document.getElementById((kind + "TcpPortID"));
    tcpPortObj.value = "";

    var commandObj = document.getElementById((kind + "CommandID"));
    commandObj.value = "";

    var initialDelaySecondsObj = document.getElementById((kind + "InitialDelaySecondsID"));
    initialDelaySecondsObj.value = "";

    var periodSecondsObj = document.getElementById((kind + "PeriodSecondsID"));
    periodSecondsObj.value = "";

    var timeoutSecondsObj = document.getElementById((kind + "TimeoutSecondsID"));
    timeoutSecondsObj.value = "";

    var failureThresholdObj = document.getElementById((kind + "FailureThresholdID"));
    failureThresholdObj.value = "";

    var successThresholdObj = document.getElementById((kind + "SuccessThresholdID"));
    successThresholdObj.value = "";

}

function addContainerToContainerList(containerType,containerName,data){
    var blockClass = "";
    var titleClass = "";
    var nameClass = "";
    var title = "";

    switch (containerType){
        case "1":
            blockClass = "newContainerFormForPrivilege";
            titleClass = "privilegeContainerTitle";
            nameClass = "privilegeContainerName";
            title = "特权容器";
            break;
        case "2":
            blockClass = "newContainerFormForInit";
            titleClass = "initContainerTitle";
            nameClass = "initContainerName";
            title = "初始化容器";
            break;
        default:
            blockClass = "newContainerFormForWork";
            titleClass = "workContainerTitle";
            nameClass = "workContainerName";
            title = "工作容器";
    }

    var newContainerHtml = "<div>";
    newContainerHtml =newContainerHtml + "<div class=\"" + blockClass + "\"  >";
    newContainerHtml =newContainerHtml + "<div> <span class=\"" + titleClass + "\" onClick=\"displayContainerFormData(this,\'0\')\" title=\"单击显示容器详情\">" + title + "</span>";
    newContainerHtml =newContainerHtml + "<span class=\"containerDelFlag\" title=\"删除本容器\" onClick=\"removeContainter(this,\'1\')\">X</span>";
    newContainerHtml =newContainerHtml + "</div> <div class=\"" + nameClass + "\" onClick=\"displayContainerFormData(this,\'1\')\" title=\"单击显示容器详情\">" + containerName + " </div>";
    newContainerHtml =newContainerHtml + "<input type=\"hidden\" name=\"containerData[]\" value=\"" + data + "\">";
    newContainerHtml =newContainerHtml + "</div></div>";

    var newContainerList = document.getElementById("newContainerLists");
    var oldContainerListHtml = newContainerList.innerHTML;
    var  containerListHtml = oldContainerListHtml + newContainerHtml;
    newContainerList.innerHTML = containerListHtml;

    return;
}

function displayContainerFormData(obj,level){
    var parentObj = null;
    if(level == "0"){
        var spanParent = obj.parentNode;
        if(spanParent == null){
            formDataShowTip(("操作错误，请刷新页面再试或联系系统管理员"),"RED",0);
            return;
        }
        parentObj = spanParent.parentNode;
    }
    if(level == "1"){
        parentObj = obj.parentNode;
    }
    if(parentObj == null){
        formDataShowTip(("操作错误，请刷新页面再试或联系系统管理员2222"),"RED",0);
        return;
    }
    var containerInputObj = parentObj.querySelector('input');
    var containerData = containerInputObj.value;
    var decodeData = Base64.decode(containerData);
    var jsonObj = JSON.parse(decodeData);

    // 回显基本信息
    var containerNameObj = document.getElementById("containerNameID");
    containerNameObj.value = jsonObj.basicInfo.containerName;

    var containerType = jsonObj.basicInfo.containerType;
    var containerTypeObj = document.getElementsByName("containerType");
    for(i = 0; i<containerTypeObj.length; i++){
        containerTypeObj[i].checked = false;
        if(containerTypeObj[i].value == containerType){
            containerTypeObj[i].checked = true;
            break;
        }
    }

    var containerImageObj = document.getElementById("containerImageID");
    containerImageObj.value = jsonObj.basicInfo.containerImage;

    var imagePullPolicyObj = document.getElementById("imagePullPolicyID");
    for(i = 0; i < imagePullPolicyObj.options.length; i++ ){
        imagePullPolicyObj.options[i].selected = false;
        if(imagePullPolicyObj.options[i].value == jsonObj.basicInfo.imagePullPolicy ){
            imagePullPolicyObj.options[i].selected = true;
        }
    }

    var startCommandObj = document.getElementById("startCommandID");
    startCommandObj.value = jsonObj.basicInfo.startCommand;

    var environmentObj = document.getElementById("environmentID");
    environmentObj.value = jsonObj.basicInfo.environment;


    // 回显quota信息
    var cpuRequestObj = document.getElementById("cpuRequestID");
    cpuRequestObj.value = jsonObj.quota.cpuRequest;

    var cpuLimitObj = document.getElementById("cpuLimitID");
    cpuLimitObj.value = jsonObj.quota.cpuLimit;

    var memRequestObj = document.getElementById("memRequestID");
    memRequestObj.value = jsonObj.quota.memRequest;

    var memLimitObj = document.getElementById("memLimitID");
    memLimitObj.value = jsonObj.quota.memLimit;

    displayHealthyFormData("startup",jsonObj.startup);
    displayHealthyFormData("readiness",jsonObj.readiness);
    displayHealthyFormData("liveness",jsonObj.liveness);
    removeContainter(parentObj,0);
    return;
}

function displayHealthyFormData(kind,data){
    var ProbeTypeObj = document.getElementsByName((kind + "ProbeType"));
    for(i = 0; i < ProbeTypeObj.length; i++){
        ProbeTypeObj[i].checked = false;
        if(ProbeTypeObj[i].value == data.probeType ){
            ProbeTypeObj[i].checked = true;
        }
    }

    var httpProtocolObj = document.getElementById((kind + "HttpProtocolID"));
    for(i = 0; i < httpProtocolObj.options.length; i++){
        httpProtocolObj.options[i].selected = false;
        if(httpProtocolObj.options[i].value == data.httpProtocol){
            httpProtocolObj.options[i].selected = true;
            break;
        }
    }

    var httpPortObj = document.getElementById((kind + "HttpPortID"));
    httpPortObj.value = data.httpPort;

    var httpPathObj = document.getElementById((kind + "HttpPathID"));
    httpPathObj.value = data.httpPath;

    httpHeaderKeyValue = data.httpHeaderKey;
    httpHeaderValueValue = data.httpHeaderValue;
    var containerHeaderBlockObj = document.getElementById("container" + kind + "HttpHeaderLineID");
    var insertHtml = "";
    for(i = 0; i < httpHeaderKeyValue.length; i++){
        insertHtml = insertHtml + "<div class=\"formDataInputline\"><span> HttpHeader Key </span>";
        insertHtml = insertHtml + "<input type=text name=\"" + kind + "HttpHeaderKey\" size=\"30\" value=\"" + httpHeaderKeyValue[i] + "\" title=\" Http Header Key\">";
        insertHtml = insertHtml + "</span><span><span>值</span><input type='text' name=\"" + kind + "HttpHeaderValue\" size='30' value=\"" + httpHeaderValueValue[i] + "\" title=值 >";
        insertHtml = insertHtml + "</span><span><a href='#' onclick='formDataWordsInputClick(\"addDeployment\",\"deployment\",\"" + kind + "HttpHeaderLineID\",\"" + kind + "HttpHeaderDelID\",\",\",\"HttpHeaderDel\",this)\">";
        insertHtml = insertHtml + "</span></div>";
    }
    containerHeaderBlockObj.innerHTML = insertHtml;

    var tcpPortObj = document.getElementById((kind + "TcpPortID"));
    tcpPortObj.value = data.tcpPort;

    var commandObj = document.getElementById((kind + "CommandID"));
    commandObj.value = data.command;

    var initialDelaySecondsObj = document.getElementById((kind + "InitialDelaySecondsID"));
    initialDelaySecondsObj.value = data.initialDelaySeconds;

    var periodSecondsObj = document.getElementById((kind + "PeriodSecondsID"));
    periodSecondsObj.value = data.periodSeconds;

    var timeoutSecondsObj = document.getElementById((kind + "TimeoutSecondsID"));
    timeoutSecondsObj.value = data.timeoutSeconds;

    var failureThresholdObj = document.getElementById((kind + "FailureThresholdID"));
    failureThresholdObj.value = data.failureThreshold;

    var successThresholdObj = document.getElementById((kind + "SuccessThresholdID"));
    successThresholdObj.value = data.successThreshold;

    return;
}

function removeContainter(obj,level){
    var parentObj = null;
    if(level == "1"){
        var spanParent = obj.parentNode;
        if(spanParent == null){
            formDataShowTip(("操作错误，请刷新页面再试或联系系统管理员"),"RED",0);
            return;
        }
        parentObj = obj.parentNode;
    }
    if(level == "0"){
        parentObj = obj;
    }
    if(parentObj == null){
        formDataShowTip(("操作错误，请刷新页面再试或联系系统管理员"),"RED",0);
        return;
    }
    parentObj = parentObj.parentNode;
    parentObj.remove();

    return;
}

function probeSwitch(kind, status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }
    var lineObj = document.getElementById("container" + kind + "ProbeHTTPGetLineID");
    lineObj.style.display = statusStr;

    return;
}

function httpProbeSwitch(kind, status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }

    var containerHttpGetLineObj = document.getElementById("line" + kind + "HttpHeaderLineID");
    containerHttpGetLineObj.style.display = statusStr;
    var httpHeaderLineObj = document.getElementById("container" + kind + "HttpHeaderLineID");
    httpHeaderLineObj.style.display = statusStr;
    var httpHeaderAnchorLineObj = document.getElementById("line" + kind + "HttpHeaderAnchorLine");
    httpHeaderAnchorLineObj.style.display = statusStr;
    var probeHttpGetLineObj = document.getElementById("line" + kind + "ProbeHTTPGetLineID");
    probeHttpGetLineObj.style.display = statusStr;

    return;
}

function tcpProbeSwitch(kind, status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }
    var probeTcpPortLineObj = document.getElementById("line" + kind + "TcpPortLineID");
    probeTcpPortLineObj.style.display = statusStr;

    return;
}

function commandProbeSwitch(kind, status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }

    var probeCommandLineObj = document.getElementById("line" + kind + "CommandLineID");
    probeCommandLineObj.style.display = statusStr;

    return;
}

function probeTypeChanged(obj,kind){
    var typeValue = obj.value;
    switch (typeValue){
        case "0":
            probeSwitch(kind,"0");
            break;
        case "1":
            probeSwitch(kind,"1");
            httpProbeSwitch(kind, "1");
            tcpProbeSwitch(kind, "0");
            commandProbeSwitch(kind, "0");
            break;
        case "2":
            probeSwitch(kind,"1");
            httpProbeSwitch(kind, "0");
            tcpProbeSwitch(kind, "1");
            commandProbeSwitch(kind, "0");
            break;
        case "3":
            probeSwitch(kind,"1");
            httpProbeSwitch(kind, "0");
            tcpProbeSwitch(kind, "0");
            commandProbeSwitch(kind, "1");
            break;
        default:
            probeSwitch(kind,"0");
            break;
    }

    return;
}


function startupProbeTypeChanged(formID,module,dcID,clusterID,namespace,groupID,option){
    probeTypeChanged(option,"startup");

    return;
}

function readinessProbeTypeChanged(formID,module,dcID,clusterID,namespace,groupID,option){
    probeTypeChanged(option,"readiness");
    return;
}

function livenessProbeTypeChanged(formID,module,dcID,clusterID,namespace,groupID,option){
    probeTypeChanged(option,"liveness");

    return;
}

function startupHttpHeaderAdd(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var startupHttpHeaderContainerObj  = document.getElementById("containerstartupHttpHeaderLineID");

    var htmlContent = "<span> <span> HttpHeader Key</span>";
    htmlContent = htmlContent + "<input  type=\"text\" name=\"startupHttpHeaderKey\" size=\"30\" value=\"\" title=\"HttpHeader Key\">";
    htmlContent = htmlContent + "</span><span>";
    htmlContent = htmlContent + "<span> 值</span>";
    htmlContent = htmlContent + "<input  type=\"text\" name=\"startupHttpHeaderValue\" size=\"30\" value=\"\" title=\"值\">";
    htmlContent = htmlContent + "</span>";
    htmlContent = htmlContent + "<span name=\"startupHttpHeaderDel[]\">";
    htmlContent = htmlContent + "<a href=\"#\" onClick=\"formDataWordsInputClick(&quot;addDeployment&quot;,&quot;deployment&quot;,&quot;startupHttpHeaderLineID&quot;,&quot;startupHttpHeaderDelID&quot;,&quot;#&quot;,&quot;workloadDelHttpHeader&quot;,this)\">";
    htmlContent = htmlContent + "<span class=\"awesomeFont\"> <b class=\"fa-trash\"></b> </span>";
    htmlContent = htmlContent + "</a></span>";

    var newDiv = document.createElement("div");
    newDiv.className = "formDataInputline";
    newDiv.innerHTML = htmlContent;
    startupHttpHeaderContainerObj.appendChild(newDiv);

    return;
}

function readinessHttpHeaderAdd(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var readinessHttpHeaderContainerObj  = document.getElementById("containerreadinessHttpHeaderLineID");

    var htmlContent = "<span> <span> HttpHeader Key</span>";
    htmlContent = htmlContent + "<input  type=\"text\" name=\"readinessHttpHeaderKey\" size=\"30\" value=\"\" title=\"HttpHeader Key\">";
    htmlContent = htmlContent + "</span><span>";
    htmlContent = htmlContent + "<span> 值</span>";
    htmlContent = htmlContent + "<input  type=\"text\" name=\"readinessHttpHeaderValue\" size=\"30\" value=\"\" title=\"值\">";
    htmlContent = htmlContent + "</span>";
    htmlContent = htmlContent + "<span name=\"readinessHttpHeaderDel[]\">";
    htmlContent = htmlContent + "<a href=\"#\" onClick=\"formDataWordsInputClick(&quot;addDeployment&quot;,&quot;deployment&quot;,&quot;readinessHttpHeaderLineID&quot;,&quot;readinessHttpHeaderDelID&quot;,&quot;#&quot;,&quot;workloadDelHttpHeader&quot;,this)\">";
    htmlContent = htmlContent + "<span class=\"awesomeFont\"> <b class=\"fa-trash\"></b> </span>";
    htmlContent = htmlContent + "</a></span>";

    var newDiv = document.createElement("div");
    newDiv.className = "formDataInputline";
    newDiv.innerHTML = htmlContent;
    readinessHttpHeaderContainerObj.appendChild(newDiv);

    return;
}

function livenessHttpHeaderAdd(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var livenessHttpHeaderContainerObj  = document.getElementById("containerlivenessHttpHeaderLineID");

    var htmlContent = "<span> <span> HttpHeader Key</span>";
    htmlContent = htmlContent + "<input  type=\"text\" name=\"livenessHttpHeaderKey\" size=\"30\" value=\"\" title=\"HttpHeader Key\">";
    htmlContent = htmlContent + "</span><span>";
    htmlContent = htmlContent + "<span> 值</span>";
    htmlContent = htmlContent + "<input  type=\"text\" name=\"livenessHttpHeaderValue\" size=\"30\" value=\"\" title=\"值\">";
    htmlContent = htmlContent + "</span>";
    htmlContent = htmlContent + "<span name=\"livenessHttpHeaderDel[]\">";
    htmlContent = htmlContent + "<a href=\"#\" onClick=\"formDataWordsInputClick(&quot;addDeployment&quot;,&quot;deployment&quot;,&quot;livenessHttpHeaderLineID&quot;,&quot;livenessHttpHeaderDelID&quot;,&quot;#&quot;,&quot;workloadDelHttpHeader&quot;,this)\">";
    htmlContent = htmlContent + "<span class=\"awesomeFont\"> <b class=\"fa-trash\"></b> </span>";
    htmlContent = htmlContent + "</a></span>";

    var newDiv = document.createElement("div");
    newDiv.className = "formDataInputline";
    newDiv.innerHTML = htmlContent;
    livenessHttpHeaderContainerObj.appendChild(newDiv);

    return;
}

function createContainerBasicInfo(containerName,containerType,containerImage,imagePullPolicy,startCommand,environment){
    var basicInfo = new Object();
    basicInfo.containerName = containerName;
    basicInfo.containerType = containerType;
    basicInfo.containerImage = containerImage;
    basicInfo.imagePullPolicy = imagePullPolicy;
    basicInfo.startCommand = startCommand;
    basicInfo.environment = environment;

    return basicInfo;
}

function createContainerQuota(cpuRequest,cpuLimit,memRequest,memLimit){
    var quota = new Object();
    quota.cpuRequest = cpuRequest;
    quota.cpuLimit = cpuLimit;
    quota.memRequest = memRequest;
    quota.memLimit = memLimit;

    return quota
}

function createContainerHealthyCheck(probeType,httpProtocol,httpPort,httpPath,httpHeaderKey,httpHeaderValue,tcpPort,command,initialDelaySeconds,
                                     periodSeconds,timeoutSeconds,failureThreshold,successThreshold){
    var probe = new Object();
    probe.probeType = probeType;
    probe.httpProtocol = httpProtocol;
    probe.httpPort = httpPort;
    probe.httpPath = httpPath;
    probe.httpHeaderKey = httpHeaderKey;
    probe.httpHeaderValue = httpHeaderValue;
    probe.tcpPort = tcpPort;
    probe.command = command;
    probe.initialDelaySeconds = initialDelaySeconds;
    probe.periodSeconds = periodSeconds;
    probe.timeoutSeconds = timeoutSeconds;
    probe.failureThreshold = failureThreshold;
    probe.successThreshold = successThreshold;

    return probe;

}

function createContainerData(basicInfo,quota,startup,readiness,liveness){
    var data = new Object();
    data.basicInfo = basicInfo;
    data.quota = quota;
    data.startup = startup;
    data.readiness = readiness;
    data.liveness = liveness;

    return data;
}

function volumeTypeChanged(formID,module,dcID,clusterID,namespace,groupID,option){
    var volumeType = option.value;
    switch (volumeType){
        case "0":
            pvcSwitch("0");
            emptyDirSwitch("0");
            hostPathSwitch("0");
            cmSwitch("0");
            secretSwitch("0");
            break;
        case "1":
            pvcSwitch("1");
            emptyDirSwitch("0");
            hostPathSwitch("0");
            cmSwitch("0");
            secretSwitch("0");
            break;
        case "2":
            pvcSwitch("0");
            emptyDirSwitch("1");
            hostPathSwitch("0");
            cmSwitch("0");
            secretSwitch("0");
            break;
        case "3":
            pvcSwitch("0");
            emptyDirSwitch("0");
            hostPathSwitch("1");
            cmSwitch("0");
            secretSwitch("0");
            break;
        case "4":
            pvcSwitch("0");
            emptyDirSwitch("0");
            hostPathSwitch("0");
            cmSwitch("1");
            secretSwitch("0");
            break;
        case "5":
            pvcSwitch("0");
            emptyDirSwitch("0");
            hostPathSwitch("0");
            cmSwitch("0");
            secretSwitch("1");
            break;
        default:
            pvcSwitch("0");
            emptyDirSwitch("0");
            hostPathSwitch("0");
            cmSwitch("0");
            secretSwitch("0");
            break;
    }

    if(volumeType == "0"){
        containerSwitch("0");
    } else {
        containerSwitch("1");
    }

    return;
}

function pvcSwitch(status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }

    var pvcObj = document.getElementById("pvcSelectID");
    pvcObj.options.length = 0 ;
    if( status == "1"){
        var clusterIDObj = document.getElementById("clusterID");
        var nsSelectedObj = document.getElementById("nsSelectedID");
        var clusterID = clusterIDObj.value;
        var ns = nsSelectedObj.options[nsSelectedObj.options.selectedIndex].value;
        var actionUri = "/api/" + apiVersion + "/pvc/getNameList?clusterID=" + clusterID +"&namespace=" + ns;
        var respValue = addWorkloadAjax("addDeployment",actionUri,"GET", false);
        if(respValue.errorCode != 0 ){
            formDataShowTip(respValue.data,"RED",0);
            return;
        }
        var pvcList =  respValue.data;
        if(pvcList.length > 0) {
            for (i = 0; i < pvcList.length; i++) {
                newOption = new Option(pvcList[i], pvcList[i]);
                pvcObj.options.add(newOption);
            }
        } else {
            newOption = new Option("请先在"+ ns + "命名空间内配置持久存储", "0",true,true);
            pvcObj.options.add(newOption);
        }

    }

    var pvcLineObj = document.getElementById("linepvcSelectLineID");
    pvcLineObj.style.display = statusStr;

    return;
}

function  emptyDirSwitch(status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }

    var emptyDirObj = document.getElementById("lineemptyDirtLineID");
    emptyDirObj.style.display =  statusStr;

    return;
}

function hostPathSwitch(status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }

    var hostPathBlock = document.getElementById("containerhostPathLineID");
    var hostPathLine = document.getElementById("linehostPathLineID");
    hostPathBlock.value = "";
    hostPathLine.value = "";
    hostPathBlock.style.display = statusStr;
    hostPathLine.style.display = statusStr;

    return;
}

function  cmSwitch(status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }
    var cmLineObj = document.getElementById("linecmSelectLineID")
    var cmObj = document.getElementById("cmSelectID");
    cmObj.options.length = 0;
    if(status == "1"){
        var clusterIDObj = document.getElementById("clusterID");
        var nsSelectedObj = document.getElementById("nsSelectedID");
        var clusterID = clusterIDObj.value;
        var ns = nsSelectedObj.options[nsSelectedObj.options.selectedIndex].value;
        var actionUri = "/api/" + apiVersion + "/configmap/getNameList?clusterID=" + clusterID +"&namespace=" + ns;
        var respValue = addWorkloadAjax("addDeployment",actionUri,"GET", false);
        if(respValue.errorCode != 0){
            formDataShowTip(respValue.data,"RED",0);
            return;
        }
        var cmList =  respValue.data;
        if(cmList.length > 0 ) {
            for (i = 0; i < cmList.length; i++) {
                newOption = new Option(cmList[i], cmList[i]);
                cmObj.options.add(newOption);
            }
        }else {
            newOption = new Option("请先在" + ns + "命名空间内配置配置字典", "0");
            cmObj.options.add(newOption);
        }
    }
    cmLineObj.style.display = statusStr;

    return;
}

function secretSwitch(status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }
    var secretLineObj = document.getElementById("linesecretSelectLineID")
    var secretObj = document.getElementById("secretSelectID");
    secretObj.options.length = 0;
    if( status == "1"){
        var clusterIDObj = document.getElementById("clusterID");
        var nsSelectedObj = document.getElementById("nsSelectedID");
        var clusterID = clusterIDObj.value;
        var ns = nsSelectedObj.options[nsSelectedObj.options.selectedIndex].value;
        var actionUri = "/api/" + apiVersion + "/secret/getNameList?clusterID=" + clusterID +"&namespace=" + ns;
        var respValue = addWorkloadAjax("addDeployment",actionUri,"GET", false);
        if(respValue.errorCode != 0){
            formDataShowTip(respValue.data,"RED",0);
            return;
        }
        var secretList =  respValue.data;
        if(secretList.length > 0 ) {
            for (i = 0; i < secretList.length; i++) {
                newOption = new Option(secretList[i], secretList[i]);
                secretObj.options.add(newOption);
            }
        } else {
            newOption = new Option("请先在" + ns + "命名空间内配置密文", "0");
            secretObj.options.add(newOption);
        }
    }
    secretLineObj.style.display = statusStr;

    return;
}

function containerSwitch(status){
    statusStr = "none";
    if(status == "1"){
        statusStr = "block";
    }

    var volumeMountBlockObj = document.getElementById("containercontainerSelectLineID");
    var volumeMountLineObj = document.getElementById("linevolumeMountLineID");
    var containerSelectLineObj = document.getElementById("linecontainerSelectLineID");
    var containerSelectObj = document.getElementById("containerSelectID");
    if(status == "1" ){
        var newContainerListsObj = document.getElementsByName("containerData[]");
        containerSelectObj.options.length = 0;
        if(newContainerListsObj.length < 1){
            newOption = new Option("请首先配置容器基本信息", "0");
            containerSelectObj.options.add(newOption);
        } else {
            for (i = 0; i < newContainerListsObj.length; i++) {
                var containerData = newContainerListsObj[i].value;
                var decodeData = Base64.decode(containerData);
                var jsonObj = JSON.parse(decodeData);
                var containerName = jsonObj.basicInfo.containerName;
                newOption = new Option(containerName, containerName);
                containerSelectObj.options.add(newOption);
            }
        }
    } else {
        containerSelectObj.options.length = 0;
        var mountPathObj = document.getElementById("volumeMountPathID");
        var mountSubPathObj = document.getElementById("volumeMountSubPathID");
        mountPathObj.value = "";
        mountSubPathObj.value = "";
    }
    containerSelectLineObj.style.display = statusStr;
    volumeMountLineObj.style.display = statusStr;
    volumeMountBlockObj.style.display = statusStr;

    return;
}

function addWorkloadAjax(formID, actionUrl,actionType,isMultiPart) {
    var data;
    if(isMultiPart){
        data = new FormData($("#"+formID)[0]);
    } else {
        data = "";
    }
    var ajaxRet = {"errorCode": "-1", "data": "未知错误"};
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
            ajaxRet.data = msg;
        },
        success: function(result) {
            ajaxRet.errorCode = result.errorCode;
            ajaxRet.data = result.data;
        }
    });

    return ajaxRet;
}

function addWorkloadValidateNewName(formID,module,dcID,clusterID,namespace,uri,obj){
    var nsSelectedObj = document.getElementById("nsSelectedID");
    var ns = nsSelectedObj.options[nsSelectedObj.options.selectedIndex].value;
    var objValue = obj.value;
    var actionUri = "/api/" + apiVersion + "/" + module + "/validateNewName?clusterID=" + clusterID +"&namespace=" + ns + "&objValue=" + objValue;
    var respValue = addWorkloadAjax(addFormId,actionUri,"GET", false);
    if(respValue.errorCode != 0){
        formDataShowTip(respValue.data,"RED",0);
        obj.focus();
        return;
    }

    return;
}

function addNewStorageMountFormData(){
    var volumeTypeIDObj = getRadioCheckedObj("volumeType");
    if(volumeTypeIDObj == null){
        formDataShowTip("请选择数据卷类型","RED",0);
        volumeTypeIDObj.focus();
        return;
    }
    var volumeType = volumeTypeIDObj.value;
    if(volumeType == "0"){
        formDataShowTip("你选择了不挂载存储","RED",0);
        volumeTypeIDObj.focus();
        return;
    }

    var volumeNameObj = document.getElementById("volumeNameID");
    var volumeName = volumeNameObj.value;
    if(volumeName == ""){
        formDataShowTip("数据卷名称不能为空","RED",0);
        volumeNameObj.focus();
        return;
    }

    var data = new Object();
    data.basicInfo = createStorageBasicInfo(volumeType,volumeName);
    data.pvcData = createStoragePvcData(volumeType);
    data.hostPathData = createStorageHostPathData(volumeType);
    var cmName = "";
    if(volumeType == "4"){
        var cmSelectObj =  document.getElementById("cmSelectID");
        cmName = cmSelectObj.options[cmSelectObj.options.selectedIndex].value;
    }
    data.cmData = cmName;

    var secretName = "";
    if(volumeType == "5"){
        var secretSelectObj =  document.getElementById("secretSelectID");
        secretName = secretSelectObj.options[secretSelectObj.options.selectedIndex].value;
    }
    data.secretData = secretName;

    var containerData = new Object();
    containerData.name = "";
    containerData.mountPath = "";
    containerData.subPath = "";
    if(volumeType != "0"){
        var containerSelectObj = document.getElementById("containerSelectID");
        var containerName = containerSelectObj.options[containerSelectObj.options.selectedIndex].value;
        var mountPathObj = document.getElementById("volumeMountPathID");
        var mountPath = mountPathObj.value;
        var mountSubPathObj = document.getElementById("volumeMountSubPathID");
        var mountSubPath = mountSubPathObj.value;
        containerData.name = containerName;
        containerData.mountPath = mountPath;
        containerData.subPath = mountSubPath;
    }
    data.containerData = containerData;

    var containerName =  containerData.name;
    var jsonData = JSON.stringify(data);
    var encodeData = Base64.encode(jsonData);

    addStorageDataToList(volumeType,containerName,volumeName,encodeData);
    resetStorageFormInputItemValue();

    return;
}

function createStorageBasicInfo(volumeType,volumeName){
    var storageBasicInfo = new Object();
    storageBasicInfo.volumeName = volumeName;
    storageBasicInfo.volumeType = volumeType;

    return storageBasicInfo;
}

function createStoragePvcData(volumeType){
    var storagePvcData = new Object();
    if(volumeType == "1"){
        var pvcSelectObj = document.getElementById("pvcSelectID");
        var pvcName =  pvcSelectObj.options[pvcSelectObj.options.selectedIndex].value;
        storagePvcData.name = pvcName;
    } else {
        storagePvcData.name = "";
    }

    return storagePvcData;
}

function createStorageHostPathData(volumeType){
    var hostPathData = new Object();
    if(volumeType == "3"){
        var hostPathObj = document.getElementById("hostPathID");
        var hostPath = hostPathObj.value;
        var hostPathTypeObj = document.getElementById("hostPathTypeID");
        var hostPathType = hostPathTypeObj.options[hostPathTypeObj.options.selectedIndex].value;
        hostPathData.hostPath = hostPath;
        hostPathData.hostPathType = hostPathType;
    } else {
        hostPathData.hostPath = "";
        hostPathData.hostPathType = "";
    }

    return hostPathData;
}

function addStorageDataToList(volumeType,containerName,volumeName, encodeData){
    var volumeTypeStr = ""
    switch (volumeType){
        case "1":
            volumeTypeStr = "持久数据卷";
            break;
        case "2":
            volumeTypeStr = "临时目录";
            break;
        case "3":
            volumeTypeStr = "HostPath";
            break;
        case "4":
            volumeTypeStr = "配置字典";
            break;
        case "5":
            volumeTypeStr = "密文";
            break;
        default:
            volumeTypeStr = "不挂载数据卷";
            break;
    }

    var trHtml = "<tr><td>" + containerName + "</td><td>" + volumeName + "</td><td>";
    trHtml = trHtml + volumeTypeStr + "</td><td>"
    trHtml = trHtml + "<input type=\"hidden\" name=\"storageMountData[]\" value=\"" + encodeData + "\">";
    trHtml = trHtml + "<a href=\"#\" onclick=\"editStorageMount(this)\" >修改</a> &nbsp;&nbsp;<a href=\"#\" onclick=\"removeStorageMount(this)\" >删除</a> </td></tr>";
    var storageListObj = document.getElementById("newStorageMountListID");
    var newRow = storageListObj.insertRow(storageListObj.rows.length);
    newRow.innerHTML = trHtml;

    return;
}

function getRadioCheckedObj(name){
    var obj = document.getElementsByName(name);
    if(obj != null){
        for(i = 0; i < obj.length; i++ ){
            if(obj[i].checked){
                return obj[i];
            }
        }
    }

    return null;
}

function resetStorageFormInputItemValue(){
    var volumeNameObj = document.getElementById("volumeNameID");
    volumeNameObj.value = "";

    var volumeTypeObj = document.getElementsByName("volumeType");
    volumeTypeObj[0].checked = true;

    pvcSwitch("0");
    emptyDirSwitch("0");
    hostPathSwitch("0");
    cmSwitch("0");
    secretSwitch("0");
    containerSwitch("0");

    return;
}

function editStorageMount(obj){
    var objParent = obj.parentElement;
    if(objParent == null ){
        formDataShowTip("操作错误，你稍后再试或联系系统管理员","RED",0);
        return;
    }

    var storageMountDataObj = objParent.querySelector('input');
    var storageMountData = storageMountDataObj.value;
    var decodeData = Base64.decode(storageMountData);
    var jsonObj = JSON.parse(decodeData);
    setStorageMountData(jsonObj);
    var trObj =  objParent.parentNode;
    trObj.remove();

    return;
}

function removeStorageMount(obj){
    var parentObj =  obj.parentNode;
    if(parentObj == null ){
        formDataShowTip("操作错误，你稍后再试或联系系统管理员","RED",0);
        return;
    }
    var trObj =  parentObj.parentNode;
    trObj.remove();
    resetStorageFormInputItemValue();
    return;
}

function setStorageMountData(jsonObj){
    var volumeNameObj = document.getElementById("volumeNameID");
    volumeNameObj.value = jsonObj.basicInfo.volumeName;
    var volumeType =  jsonObj.basicInfo.volumeType;
    var volumeTypeObj = document.getElementsByName("volumeType");
    for(i = 0; i < volumeTypeObj.length; i++ ){
        if(volumeType == volumeTypeObj[i].value){
            volumeTypeObj[i].checked = true;
            break;
        }
    }

    switch (volumeType){
        case "1":
            pvcSwitch("1");
            emptyDirSwitch("0");
            hostPathSwitch("0");
            cmSwitch("0");
            secretSwitch("0");
            var pvcSelectObj = document.getElementById("pvcSelectID");
            for(i = 0; i < pvcSelectObj.options.length; i++){
                if(pvcSelectObj.options[i].value == jsonObj.pvcData.name){
                    pvcSelectObj.options[i].selected = true;
                    break;
                }
            }
            break;
        case "2":
            pvcSwitch("0");
            emptyDirSwitch("1");
            hostPathSwitch("0");
            cmSwitch("0");
            secretSwitch("0");
            break;
        case "3":
            pvcSwitch("0");
            emptyDirSwitch("0");
            hostPathSwitch("1");
            cmSwitch("0");
            secretSwitch("0");
            var hostPath = jsonObj.hostPathData.hostPath;
            var hostPathType = jsonObj.hostPathData.hostPathType;
            var hostPathObj = document.getElementById("hostPathID");
            var hostPathTypeObj = document.getElementById("hostPathTypeID");
            hostPathObj.value = hostPath;
            for(i = 0; i < hostPathTypeObj.options.length; i++){
                if(hostPathTypeObj.options[i].value == hostPathType){
                    hostPathTypeObj.options[i].selected = true;
                    break;
                }
            }
            break;
        case "4":
            pvcSwitch("0");
            emptyDirSwitch("0");
            hostPathSwitch("0");
            cmSwitch("1");
            secretSwitch("0");
            var cmName = jsonObj.cmData;
            var cmSelectObj =  document.getElementById("cmSelectID");
            for(i = 0; i < cmSelectObj.options.length; i++){
                if(cmSelectObj.options[i].value == cmName){
                    cmSelectObj.options[i].selected = true;
                    break;
                }
            }
            break;
        case "5":
            pvcSwitch("0");
            emptyDirSwitch("0");
            hostPathSwitch("0");
            cmSwitch("0");
            secretSwitch("1");
            var secretName = jsonObj.secretData;
            var secretSelectObj =  document.getElementById("secretSelectID");
            for(i = 0; i < secretSelectObj.options.length; i++){
                if(secretSelectObj.options[i].value == secretName ){
                    secretSelectObj.options[i].selected = true;
                    break;
                }
            }
            break;
    }

    containerSwitch("1");

    var containerName =  jsonObj.containerData.name;
    var mountPath = jsonObj.containerData.mountPath;
    var subPath = jsonObj.containerData.subPath;
    var containerSelectObj = document.getElementById("containerSelectID");
    for(i = 0; i < containerSelectObj.options.length; i++){
        if(containerSelectObj.options[i].value == containerName) {
            containerSelectObj.options[i].selected = true;
            break;
        }
    }

    var mountPathObj = document.getElementById("volumeMountPathID");
    mountPathObj.value = mountPath;
    var mountSubPathObj = document.getElementById("volumeMountSubPathID");
    mountSubPathObj.value = subPath;

    return;
}

function addWorkload(formID, module, dcID, clusterID, namespace, actionType){
    var addTypeObj = document.getElementById("addType");
    var addTypeValue = addTypeObj.value;
    if(addTypeValue == "0"){
        var editor = ace.edit("addWorkloadEditor");
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

    var actionUri = "/api/" + apiVersion + "/" + module + "/add?clusterID=" + clusterID;
    var respValue = addWorkloadAjax("addDeployment",actionUri,"POST", true);
    $('#container').load(lastUrl);
    if(respValue.errorCode != 0 ){
        formDataShowTip(respValue.data,"RED",0);
        return;
    }

    return formDataShowTip("内容已经添加成功","GREEN",0);
}

function addObjClickFileButton(objectID){
    var obj = document.getElementById(objectID);
    obj.click();
}

function workloadAddLabel(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var labelContainerObj = document.getElementById("containernewNsLabel");
    var htmlContent ="<span><span>标签</span><input type=\"text\" name=\"labelKey[]\" size=\"30\" title=\"标签\"></span>";
    htmlContent = htmlContent + "<span> = </span>";
    htmlContent = htmlContent + "<span><span>值</span><input type=\"text\" name=\"labelValue[]\" size=\"30\" title=\"值\"></span>";
    htmlContent = htmlContent + "<span>";
    htmlContent = htmlContent + "<a href=\"#\" onclick=\"formDataWordsInputClick(&quot;addDeployment&quot;,&quot;deployment&quot;,&quot;newNsLabel&quot;,&quot;delLabel&quot;,&quot;#&quot;,&quot;workloadDelLabel&quot;,this)\">";
    htmlContent = htmlContent + "<span class=\"awesomeFont\"><b class=\"fa-trash\"></b></span></a></span>";

    var newLabelLineObj = document.createElement("div");
    newLabelLineObj.className = "WorkloadForInputline";
    newLabelLineObj.innerHTML = htmlContent;
    labelContainerObj.appendChild(newLabelLineObj);

    return;
}

function workloadDelLabel(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var objParent = obj.parentNode;
    var parentObj = objParent.parentNode;
    parentObj.remove();

    return;
}

function workloadAddAnnotation(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var annotaionContainerObj =  document.getElementById("containernewAnnotationLabel");
    var htmlContent ="<span><span>注解</span><input type=\"text\" name=\"annotationKey[]\" size=\"30\" title=\"注解\"></span>";
    htmlContent = htmlContent + "<span> = </span>";
    htmlContent = htmlContent + "<span><span>值</span><input type=\"text\" name=\"annotationValue[]\" size=\"30\" title=\"值\"></span>";
    htmlContent = htmlContent + "<span>";
    htmlContent = htmlContent + "<a href=\"#\" onclick=\"formDataWordsInputClick(&quot;addDeployment&quot;,&quot;deployment&quot;,&quot;newAnnotationLabel&quot;,&quot;annotationLabel&quot;,&quot;#&quot;,&quot;workloadDelAnnotaion&quot;,this)\">";
    htmlContent = htmlContent + "<span class=\"awesomeFont\"><b class=\"fa-trash\"></b></span></a></span>";

    var newAnnotationObj = document.createElement("div");
    newAnnotationObj.className = "WorkloadForInputline";
    newAnnotationObj.innerHTML = htmlContent;
    annotaionContainerObj.appendChild(newAnnotationObj);

    return;
}

function workloadDelAnnotaion(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var objParent = obj.parentNode;
    var parentObj = objParent.parentNode;
    parentObj.remove();

    return;
}

function workloadAddSelector(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var selectorContainerObj =  document.getElementById("containerselectorLabel");
    var htmlContent ="<span><span>标签选择器</span><input type=\"text\" name=\"selectorKey[]\" size=\"30\" title=\"标签选择器\"></span>";
    htmlContent = htmlContent + "<span> = </span>";
    htmlContent = htmlContent + "<span><span>值</span><input type=\"text\" name=\"selectorValue[]\" size=\"30\" title=\"值\"></span>";
    htmlContent = htmlContent + "<span>";
    htmlContent = htmlContent + "<a href=\"#\" onclick=\"formDataWordsInputClick(&quot;addDeployment&quot;,&quot;deployment&quot;,&quot;selectorLabel&quot;,&quot;selectorLabel&quot;,&quot;#&quot;,&quot;workloadDelSelector&quot;,this)\">";
    htmlContent = htmlContent + "<span class=\"awesomeFont\"><b class=\"fa-trash\"></b></span></a></span>";

    var newSelectorObj = document.createElement("div");
    newSelectorObj.className = "WorkloadForInputline";
    newSelectorObj.innerHTML = htmlContent;
    selectorContainerObj.appendChild(newSelectorObj);

    return;
}

function workloadDelSelector(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var objParent = obj.parentNode;
    var parentObj = objParent.parentNode;
    parentObj.remove();

    return;
}


function workloadDelHttpHeader(formID,module,dcID,clusterID,namespace,lineID,itemID,uri,obj){
    var objParent = obj.parentNode;
    var parentObj = objParent.parentNode;
    parentObj.remove();

    return;
}