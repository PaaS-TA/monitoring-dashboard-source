var Utils = {
    requestHttp : async function (url, method, header, params) {
        let defaultHeader = {
            'Accept' : 'application/json',
            'Content-Type' : 'application/x-www-form-urlencoded'
        };
        Object.assign(defaultHeader, header);

        let body = params;
        if (defaultHeader["Content-Type"] == 'application/x-www-form-urlencoded') {
            body = this._makeUrlSearchParams(params);
        } else {
            body = JSON.stringify(params);
        }

        if (method == 'GET')
            delete defaultHeader['Content-Type'];

        const response = await fetch(url, {
            'method' : method,
            'mode' : 'cors',
            'headers' : defaultHeader,
            'body' : body
        });

        if (!response.ok) {
            console.error("error");
        }
        return response.json();
    },

    _makeUrlSearchParams : function (jsonData) {
        if (jsonData == null)
            return jsonData;
        else
            return Object.keys(jsonData).map(key => encodeURIComponent(key) + '=' + encodeURIComponent(jsonData[key])).join('&');
    }


}