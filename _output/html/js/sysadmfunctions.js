function isPort(port) {
    if (/^[1-9]\d*|0$/.test(port) && port * 1 >= 0 && port * 1 <= 65535){
        return true
    }
    return false;
}