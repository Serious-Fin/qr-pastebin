import { type User, tryCreateSessionForUser } from '$lib/user';
import { redirect } from '@sveltejs/kit';
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
		const user: User = {
			name: data.get('name') as string,
			password: data.get('password') as string
		};

		const redirectTo = (data.get('redirectTo') as string) ?? '/';
		let sessionId = '';
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
