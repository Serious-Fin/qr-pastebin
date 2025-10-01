import { tryCreateSessionForUser } from '$lib/user';
import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';

export const POST: RequestHandler = async ({ url, request, cookies }) => {
	const formData = await request.formData();
	const name = formData.get('name') as string;
	const password = formData.get('password') as string;

	let sessionId: string;
	try {
		sessionId = await tryCreateSessionForUser({ name, password });
	} catch (err) {
		if (err instanceof Error) {
			return json({ errMsg: err.message }, { status: 401 });
		} else {
			return json({ errMsg: 'Unknown error, try again later' }, { status: 500 });
		}
	}

	cookies.set('session', sessionId, {
		httpOnly: true,
		sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 * 24 * 7
	});

	const redirectTo = url.searchParams.get('redirectTo');
	const decodedRedirectTo = redirectTo ? decodeURIComponent(redirectTo) : '/';
	return json({ success: true, redirectTo: decodedRedirectTo });
};
