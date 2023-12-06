function detailsChangeCard(cardNo){
    var cardHead = document.getElementsByClassName("detailsCardheadline");
    var headSpans = cardHead[0].getElementsByTagName("span");
    var cardcontentdiv = document.getElementsByName("cardcontentdiv");

    for(i = 0; i < headSpans.length; i++ ){
        headSpans[i].className = "";
        cardcontentdiv[i].style.display = "none";
        if(i == cardNo){
            headSpans[i].className = "detailsActivecard";
            cardcontentdiv[i].style.display = "block";

        }
    }
    return;
}