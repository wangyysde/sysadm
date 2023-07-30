function isPort(port) {
    if (/^[1-9]\d*|0$/.test(port) && port * 1 >= 0 && port * 1 <= 65535){
        return true
    }
    return false;
}

// 将经过base64编码的utf-8字符串解码为普通字符串，支持中文
function b64DecodeUnicode(str) {
    return decodeURIComponent(Array.prototype.map.call(atob(str), function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
}

// 将utf-8字符串编码为base64字符串，支持中文
function b64EncodeUnicode(str) {
    return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, function(match, p1) {
        return String.fromCharCode('0x' + p1);
    }));
}