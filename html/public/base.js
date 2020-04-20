const postToServer = function (url, json_object, onSuccess, onFailure) {
    let request = new XMLHttpRequest();
    request.onreadystatechange = function () {
        if (request.readyState === XMLHttpRequest.DONE) {
            let response = request.responseText;
            if (request.status >= 200 && request.status < 400) {
                onSuccess(response);
            } else {
                onFailure(response);
            }
        }
    };
    request.open("POST", url);
    request.setRequestHeader("Content-Type","text/plain; charset=UTF-8");
    request.send(JSON.stringify(json_object));
};

const getFromServer = function (url, onSuccess, onFailure) {
    let request = new XMLHttpRequest();
    request.onreadystatechange = function () {
        if (request.readyState === XMLHttpRequest.DONE) {
            let response = request.responseText;
            if (request.status >= 200 && request.status < 400) {
                onSuccess(response);
            } else {
                onFailure(response);
            }
        }
    };
    request.open("GET", url);
    request.send();
};

const genericFailure = function (response) {
    alert(response);
};