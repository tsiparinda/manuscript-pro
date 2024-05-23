function redirectAfter(url, time) {
    window.setTimeout(function () {
        window.location.href = url
    }, time)
}