import type { Handle } from '@sveltejs/kit';
import { tryGetSessionForUser } from '$lib/user';

export const handle: Handle = async ({ event, resolve }) => {
	const sessionId = event.cookies.get('session');
	if (sessionId) {
		event.locals.sessionId = sessionId;

		try {
			const user = await tryGetSessionForUser(sessionId);
			event.locals.user = {
				id: user.id,
				name: user.name,
				role: user.role
			};
		} catch (err) {
			if (err instanceof Error) {
				console.error(JSON.stringify(err));
			}
		}
	}
	return resolve(event);
};
