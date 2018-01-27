// Code generated by protoc-gen-king_browser 1.0.0, DO NOT EDIT.
// Source: test/example/example.proto

class ContactMessagesServiceClient {
	constructor({server = '', authorization = '', hook = null} = {}) {
		this.server = server;
		this.authorization = authorization;
		this.hook = hook;
	}
	
	Foo(req) {
		return this._doRequest('Foo', req);
	}
	
	Bar(req) {
		return this._doRequest('Bar', req);
	}
	
	_doRequest(method, req) {
		let opts = {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(req),
		};
		if (this.authorization) {
			opts.headers.Authorization = this.authorization;
		}
		return fetch(`${this.server}/_/king.example.ContactMessagesService/${method}`, opts)
			.then(resp => {
				if (resp.status !== 200) {
					let err = new Error(response.statusText);
					err.response = response;
					throw err;
				}

				if (hook) {
					hook(response);
				}

				return response.json();
			});
	}
}
