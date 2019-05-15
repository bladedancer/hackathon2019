async function register() {
    const register = document.getElementById('register');
    const status = document.getElementById('status');
    const name = document.getElementById('name');

    status.innerHTML = `Registering ${name.value}`;
    register.classList.add('busy');
    
    const req = {
        method: 'GET',
        url: window.location + "/api/v1/totp/register",
        params: {
            name: name.value
        }
    };

    try {
        const response = await axios(req);
        console.log(response);
        if (response.status !== 201) {
            throw response;
        }

        register.classList.add('complete');
        register.classList.remove("busy");

        new QRCode(document.getElementById("qrcode"), {
            text: response.data.secret.otpauth_url,
            width: 128,
            height: 128,
            colorDark : "#000000",
            colorLight : "#ffffff",
            correctLevel : QRCode.CorrectLevel.H
        });

    } catch(err) {
        status.innerHTML = err;
    }
}

(function() {
    document.getElementById('submit').addEventListener('click', register);

})();

