import {
	createNewUser,
	UserAlreadyExistsError,
	type User,
	tryCreateSessionForUser
} from '$lib/user';
import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	return {
		userId: locals.user?.id ?? -1,
		username: locals.user?.name ?? 'Anon'
	};
};

export const actions = {
	createNewUser: async ({ request, cookies }) => {
		const data = await request.formData();
		const user: User = {
			name: data.get('name') as string,
			password: data.get('password') as string
		};
		try {
			await createNewUser(user);
		} catch (err) {
			if (err instanceof UserAlreadyExistsError) {
				return {
					errMsg: err.message,
					user: {
						name: user.name,
						password: user.password
					}
				};
			}
			return fail(500, {
				message: err instanceof Error ? err.message : 'Unknown error'
			});
		}

		let sessionId = '';
		const redirectTo = (data.get('redirectTo') as string) ?? '/';
		try {
			sessionId = await tryCreateSessionForUser(user);
		} catch (err) {
			if (err instanceof Error) {
				return {
					errMsg: err.message,
					user: {
						name: user.name,
						password: user.password
					}
				};
			}
		}

		cookies.set('session', sessionId, {
			httpOnly: true,
			sameSite: 'lax',
			path: '/',
			maxAge: 60 * 60 * 24 * 7
		});
		throw redirect(303, redirectTo);
	}
};
