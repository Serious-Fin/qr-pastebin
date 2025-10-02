import type { Handle } from '@sveltejs/kit';
import { tryGetSessionForUser } from '$lib/user';

export const handle: Handle = async ({ event, resolve }) => {
	const sessionId = event.cookies.get('session');
	if (sessionId) {
		event.locals.sessionId = sessionId;

		try {
			console.log(sessionId);
			const user = await tryGetSessionForUser(sessionId);
			console.log(user);
			event.locals.user = {
				id: user.id,
				name: user.name
			};
		} catch {}
	}
	return resolve(event);
};
