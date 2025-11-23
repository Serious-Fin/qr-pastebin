import {
	type UserCredentials,
	UserUsingOauthError,
	WrongNameOrPassError,
	tryCreateSessionForUser
} from '$lib/user';
import { redirect, fail } from '@sveltejs/kit';
import type { Actions } from './$types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	return {
		userId: locals.user?.id ?? -1,
		username: locals.user?.name ?? 'Anon'
	};
};

export const actions: Actions = {
	login: async ({ request, cookies }) => {
		const data = await request.formData();
		const user: UserCredentials = {
			name: data.get('name') as string,
			password: data.get('password') as string
		};

		const redirectTo = (data.get('redirectTo') as string) ?? '/';
		let sessionId = '';
		try {
			sessionId = await tryCreateSessionForUser(user);
		} catch (err) {
			if (err instanceof Error) {
				return fail(
					err instanceof WrongNameOrPassError || err instanceof UserUsingOauthError ? 400 : 500,
					{ message: err.message }
				);
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
