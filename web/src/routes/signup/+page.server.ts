import { createNewUser, UserAlreadyExistsError, type User } from '$lib/user';
import { fail } from '@sveltejs/kit';

export const actions = {
	createNewUser: async ({ request }) => {
		const data = await request.formData();
		const params: User = {
			name: data.get('name') as string,
			password: data.get('password') as string
		};
		try {
			await createNewUser(params);
		} catch (err) {
			if (err instanceof UserAlreadyExistsError) {
				return {
					errMsg: err.message,
					user: {
						name: params.name,
						password: params.password
					}
				};
			}
			return fail(500, {
				message: err instanceof Error ? err.message : 'Unknown error'
			});
		}

		// Try sign user in
		const form = new FormData();
		form.append('name', params.name);
		form.append('password', params.password);

		// Call your POST endpoint
		try {
			await fetch('/api/login', {
				method: 'POST',
				body: form
			});
		} catch (err) {
			if (err instanceof Error) {
				return {
					errMsg: err.message,
					user: {
						name: params.name,
						password: params.password
					}
				};
			}
		}
	}
};
