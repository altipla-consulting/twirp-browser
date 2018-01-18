
// Code generated by protoc-gen-twirp_browser 1.0.0, DO NOT EDIT.
// Source: test/example.proto

class HaberdasherClient {
	constructor(server) {
		this.server = server;
	}
	
	MakeHat(req) {
		return this._doRequest('MakeHat', req);
	}
	
	MakeHat2(req) {
		return this._doRequest('MakeHat2', req);
	}
	
	_doRequest(method, req) {
		return fetch(`${this.server}/twirp/twitch.twirp.example.Haberdasher/${method}`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(req),
		});
	}
}
