function redirectPage(actionUri) {
    var url = pageUrl + actionUri;
    $('#container').load(url);
    return;
}