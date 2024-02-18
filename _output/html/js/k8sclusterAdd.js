function k8sclusterAddRadioClick(actionUri,obj,option){
    var radioValue = "0";
    if(option){
        radioValue = option.value;
    }

    var apiServerLineObj = document.getElementById("apiServerLine");
    var caLineObj = document.getElementById("caLine");
    var certLineObj = document.getElementById("certLine");
    var keyLineObj = document.getElementById("keyLine");
    var tokenStringLineObj = document.getElementById("tokenStringLine");
    var tokenLineObj = document.getElementById("tokenLine");
    var kubeConfigLineObj = document.getElementById("kubeConfigLine");

    switch (radioValue){
        case "0":
            apiServerLineObj.style.display = "";
            caLineObj.style.display = "";
            certLineObj.style.display = "";
            keyLineObj.style.display = "";
            tokenStringLineObj.style.display = "none";
            tokenLineObj.style.display = "none";
            kubeConfigLineObj.style.display = "none";
            break;
        case "1":
            apiServerLineObj.style.display = "";
            caLineObj.style.display = "";
            certLineObj.style.display = "none";
            keyLineObj.style.display = "none";
            tokenStringLineObj.style.display = "";
            tokenLineObj.style.display = "";
            kubeConfigLineObj.style.display = "none";
            break;
        case "2":
            apiServerLineObj.style.display = "none";
            caLineObj.style.display = "none";
            certLineObj.style.display = "none";
            keyLineObj.style.display = "none";
            tokenStringLineObj.style.display = "none";
            tokenLineObj.style.display = "none";
            kubeConfigLineObj.style.display = "";
            break;
    }

    return;
}

function k8sclusterAddClickButton(){
   window.open("/download/createtokencommand.txt","_blank");
   return;
}
