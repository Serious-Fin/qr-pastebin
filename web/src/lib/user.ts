export interface UserCredentials {
	name: string;
	password: string;
}

export interface User {
	id: number;
	name: string;
	role: number;
}

export class UserAlreadyExistsError extends Error {
	constructor() {
		super('Name already taken, try a different one');
		this.name = 'UserAlreadyExistsError';
		Object.setPrototypeOf(this, UserAlreadyExistsError.prototype);
	}
}
export class WrongNameOrPassError extends Error {
	constructor() {
		super('Wrong name or password');
		this.name = 'WrongNameOrPassError';
		Object.setPrototypeOf(this, WrongNameOrPassError.prototype);
	}
}

export async function createNewUser(user: UserCredentials) {
	try {
		const response = await fetch(`http://localhost:8080/user`, {
			body: JSON.stringify(user),
			headers: {
				'Content-Type': 'application/json'
			},
			method: 'POST'
		});
		if (response.status === 409) {
			throw new UserAlreadyExistsError();
		}
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error creating user ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
	} catch (err) {
		if (err instanceof UserAlreadyExistsError) {
			throw err;
		}
		if (err instanceof Error) {
			throw Error(`Could not call create user endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while creating user: ${JSON.stringify(err)}`);
	}
}

export async function tryCreateSessionForUser(user: UserCredentials): Promise<string> {
	try {
		const response = await fetch(`http://localhost:8080/user/session`, {
			body: JSON.stringify(user),
			headers: {
				'Content-Type': 'application/json'
			},
			method: 'POST'
		});
		if (response.status === 401) {
			throw new WrongNameOrPassError();
		}
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error creating session ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		const parsedResponse: { sessionId: string } = await response.json();
		return parsedResponse.sessionId;
	} catch (err) {
		if (err instanceof WrongNameOrPassError) {
			throw err;
		}
		if (err instanceof Error) {
			throw Error(`Could not call create session endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while creating session: ${JSON.stringify(err)}`);
	}
}

export async function tryGetSessionForUser(sessionId: string): Promise<User> {
	try {
		const response = await fetch(`http://localhost:8080/user/session/${sessionId}`);
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error getting session for user ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		return await response.json();
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call get session endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while getting session: ${JSON.stringify(err)}`);
	}
}
