const speakeasy = require("speakeasy");

function secret(req, cb) {
	const name = req.params.name || 'Axway Demo';
	const secret = speakeasy.generateSecret({
		name,
		length: 32,
		issuer: 'Axway Demo'
	});
	cb.next(null, secret);
}

function token(req, cb) {
	const secret = req.params.secret;
	const ttl = req.params.ttl || 30;

	const resp = {
        "token": speakeasy.totp({
            secret: secret,
            encoding: "base32"
        }),
        "remaining": (ttl - Math.floor((new Date()).getTime() / 1000.0 % ttl))
	};
	cb.next(null, resp);
}

function validate(req, cb) {
	console.log (req);
	const token = req.params.token;
	const secret = req.params.secret;
	const valid = speakeasy.totp.verify({
		secret: secret,
		encoding: "base32",
		token: token,
		window: 0
	});

	valid ? cb.valid() : cb.invalid();
}

exports = module.exports = {
	secret,
	token,
	validate
}