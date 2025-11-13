import { PUBLIC_API_ADDRESS } from '$env/static/public';

export interface ShareRequest {
	title: string;
	content: string;
	setPassword: boolean;
	password: string;
	expireIn: string;
	hideAuthor: boolean;
	authorId: number;
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
	authorName?: string;
	hideAuthor: boolean;
}

export interface GetPasswordProtectedShareRequest {
	password: string;
}

export enum FetchShareStatus {
	NotFound = 1,
	NeedPassword,
	Accessible
}

export class WrongPasswordError extends Error {
	constructor() {
		super('Wrong password, try again');
		this.name = 'WrongPasswordError';
		Object.setPrototypeOf(this, WrongPasswordError.prototype);
	}
}

export async function createShare(request: ShareRequest): Promise<string> {
	try {
		console.log(`Making request to ${PUBLIC_API_ADDRESS}/share`);
		console.log(`Trying health:`);
		const resp = await fetch(`${PUBLIC_API_ADDRESS}/health`);
		if (!resp.ok) {
			console.error('Health failed :(');
			console.error('Status' + resp.status);
		} else {
			console.log('health success!');
		}
		const response = await fetch(`${PUBLIC_API_ADDRESS}/share`, {
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
		const response = await fetch(`${PUBLIC_API_ADDRESS}/share/${id}`);
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

export async function getShareForEdit(id: string, sessionId: string): Promise<Share> {
	try {
		const response = await fetch(`${PUBLIC_API_ADDRESS}/share/${id}/edit`, {
			headers: {
				Authorization: `Bearer ${sessionId}`
			}
		});
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
		const response = await fetch(`${PUBLIC_API_ADDRESS}/share/${id}/protected`);
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error fetching password status ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		const parsedResponse: { isPasswordProtected: boolean } = await response.json();
		return parsedResponse.isPasswordProtected;
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call is password protected endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(
			`Unknown error while calling password protection endpoint: ${JSON.stringify(err)}`
		);
	}
}

export async function getPasswordProtectedShare(
	id: string,
	body: GetPasswordProtectedShareRequest
): Promise<Share> {
	try {
		const response = await fetch(`${PUBLIC_API_ADDRESS}/share/${id}/protected`, {
			body: JSON.stringify(body),
			headers: {
				'Content-Type': 'application/json'
			},
			method: 'POST'
		});
		if (response.status === 401) {
			throw new WrongPasswordError();
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
			throw err;
		}
		if (err instanceof Error) {
			throw Error(`Could not call get share endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while getting share: ${JSON.stringify(err)}`);
	}
}

export async function getSharesForUser(sessionId: string): Promise<Share[]> {
	try {
		const response = await fetch(`${PUBLIC_API_ADDRESS}/shares`, {
			headers: {
				Authorization: `Bearer ${sessionId}`
			}
		});
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error getting shares ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		return await response.json();
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call get shares endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while getting shares: ${JSON.stringify(err)}`);
	}
}

export async function deleteShare(shareId: string, sessionId: string): Promise<void> {
	try {
		const response = await fetch(`${PUBLIC_API_ADDRESS}/share/${shareId}`, {
			headers: {
				Authorization: `Bearer ${sessionId}`
			},
			method: 'DELETE'
		});
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error deleting share ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call delete share endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while deleting share: ${JSON.stringify(err)}`);
	}
}

export async function editShare(
	request: ShareRequest,
	sessionId: string,
	shareId: string
): Promise<void> {
	try {
		const response = await fetch(`${PUBLIC_API_ADDRESS}/share/${shareId}/edit`, {
			body: JSON.stringify(request),
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${sessionId}`
			},
			method: 'PATCH'
		});
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error editing share ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call edit share endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while editing share: ${JSON.stringify(err)}`);
	}
}
