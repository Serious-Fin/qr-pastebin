export interface Share {
	content: string;
	title?: string;
	password?: string;
	expireIn?: string;
}

interface CreateShareResponse {
	id: string;
}

export async function createShare(share: Share): Promise<string> {
	try {
		const response = await fetch(`http://localhost:8080/share`, {
			body: JSON.stringify({ share }),
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
