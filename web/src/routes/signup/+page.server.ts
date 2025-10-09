import { createNewUser, type UserCredentials, tryCreateSessionForUser } from '$lib/user';
import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { UserAlreadyExistsError } from '$lib/user';

export const load: PageServerLoad = ({ locals }) => {
	return {
		userId: locals.user?.id ?? -1,
		username: locals.user?.name ?? 'Anon'
	};
};

export const actions = {
	createNewUser: async ({ request, cookies }) => {
		const data = await request.formData();
		const user: UserCredentials = {
			name: data.get('name') as string,
			password: data.get('password') as string
		};
		try {
			await createNewUser(user);
		} catch (err) {
			if (err instanceof Error) {
				return fail(err instanceof UserAlreadyExistsError ? 400 : 500, { message: err.message });
			}
			return fail(500, { message: 'Unknown server error' });
		}

		let sessionId = '';
		const redirectTo = (data.get('redirectTo') as string) ?? '/';
		try {
			sessionId = await tryCreateSessionForUser(user);
		} catch (err) {
			if (err instanceof Error) {
				return fail(400, { message: err.message });
			}
			return fail(500, { message: 'Unknown server error' });
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
