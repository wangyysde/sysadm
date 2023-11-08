function addObjRadioCustizeAction(actionUrl,objID, relatedIsDisplay,subObjID,radioOption){
    var radioValue = radioOption.value;
    var objSelectID = "spanobjectName";
    var objSelect = document.getElementById(objSelectID);
    var cronID = "spancrontab";
    var objCron = document.getElementById(cronID);

    if(radioValue == "4"){
        objSelect.style.display = "none";
        objCron.style.display = "block";
    } else {
        objSelect.style.display = "block";
        objCron.style.display = "none";
    }

    return;
}