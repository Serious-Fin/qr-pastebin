export interface CreateShareRequest {
	content: string;
	title?: string;
	password?: string;
	expireIn?: string;
}

interface CreateShareResponse {
	id: string;
}

export interface Share {
	id: string;
	content: string;
	title?: string;
	expiresIn?: string;
	isPasswordProtected?: boolean;
	author?: string
}

export interface GetPasswordProtectedShareRequest {
	password: string
}

export enum FetchShareStatus {
	NotFound = 1,
	NeedPassword,
	Accessible
}

export class WrongPasswordError extends Error {
    constructor() {
        super("Wrong password, try again");
        this.name = "WrongPasswordError";
        Object.setPrototypeOf(this, WrongPasswordError.prototype);
    }
}

export async function createShare(request: CreateShareRequest): Promise<string> {
	try {
		const response = await fetch(`http://localhost:8080/share`, {
			body: JSON.stringify(request),
			headers: {
				'Content-Type': 'application/json'
			},
			method: 'POST'
		});
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error creating share ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		const parsedResponse: CreateShareResponse = await response.json();
		return parsedResponse.id;
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call create share endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while creating share: ${JSON.stringify(err)}`);
	}
}

export async function getShare(id: string): Promise<Share> {
	try {
		const response = await fetch(`http://localhost:8080/share/${id}`);
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error getting share ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		return await response.json();
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call get share endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while getting share: ${JSON.stringify(err)}`);
	}
}

export async function isSharePasswordProtected(id: string): Promise<boolean> {
	try {
		const response = await fetch(`http://localhost:8080/share/${id}/protected`);
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error fetching password status ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		const parsedResponse: {isPasswordProtected: boolean} = await response.json()
		return parsedResponse.isPasswordProtected;
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call is password protected endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while calling password protection endpoint: ${JSON.stringify(err)}`);
	}
}

export async function getPasswordProtectedShare(id: string, body: GetPasswordProtectedShareRequest): Promise<Share> {
	try {
		const response = await fetch(`http://localhost:8080/share/${id}/protected`, {
			body: JSON.stringify(body),
			headers: {
				'Content-Type': 'application/json'
			},
			method: 'POST'
		});
		if (response.status === 401) {
			throw new WrongPasswordError()
		}
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error getting share ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		return await response.json();
	} catch (err) {
		if (err instanceof WrongPasswordError) {
			throw err
		}
		if (err instanceof Error) {
			throw Error(`Could not call get share endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while getting share: ${JSON.stringify(err)}`);
	}
}
