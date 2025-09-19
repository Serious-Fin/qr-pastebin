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

export async function getShare(id: string): Promise<string> {
	try {
		const response = await fetch(`http://localhost:8080/share/${id}`);
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error getting share ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		const parsedResponse: Share = await response.json();
		return parsedResponse.id;
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call get share endpoint: ${JSON.stringify(err.message)}`);
		}
		throw new Error(`Unknown error while getting share: ${JSON.stringify(err)}`);
	}
}
