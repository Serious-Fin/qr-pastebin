export interface User {
    name: string;
    password: string
}

export class UserAlreadyExistsError extends Error {
    constructor() {
        super("Name already taken, try a different one");
        this.name = "UserAlreadyExistsError";
        Object.setPrototypeOf(this, UserAlreadyExistsError.prototype);
    }
}

export async function createNewUser(user: User) {
    try {
		const response = await fetch(`http://localhost:8080/user`, {
			body: JSON.stringify(user),
			headers: {
				'Content-Type': 'application/json'
			},
			method: 'POST'
		});
        if (response.status === 409) {
            throw new UserAlreadyExistsError()
        }
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error creating share ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
	} catch (err) {
        if (err instanceof UserAlreadyExistsError) {
            throw err
        }
		if (err instanceof Error) {
			throw Error(`Could not call create share endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while creating share: ${JSON.stringify(err)}`);
	}
}