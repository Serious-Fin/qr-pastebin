import type { PageServerLoad } from './$types';
import {
	getShare,
	isSharePasswordProtected,
	FetchShareStatus,
	type Share,
	type GetPasswordProtectedShareRequest,
	getPasswordProtectedShare,
	WrongPasswordError,
	deleteShare
} from '$lib/share';
import { fail } from '@sveltejs/kit';
import { GOOGLE_API_KEY } from '$env/static/private';

export const load: PageServerLoad = async ({ params, locals }) => {
	// Check the role of viewing user
	const role = locals.user?.role;
	let isAdmin = false;
	if (role && role === 1) {
		isAdmin = true;
	}

	// Load "Enter password" view if share is password protected (except if user is admin)
	let hasPassword: boolean;
	try {
		hasPassword = await isSharePasswordProtected(params.id);
	} catch {
		return {
			status: FetchShareStatus.NotFound
		};
	}
	if (!isAdmin && hasPassword) {
		return {
			status: FetchShareStatus.NeedPassword
		};
	}

	// GET share if it's not password protected
	let share: Share;
	try {
		share = await getShare(params.id);
	} catch {
		return {
			status: FetchShareStatus.NotFound
		};
	}
	return {
		share: share,
		status: FetchShareStatus.Accessible,
		isAdmin
	};
};

export const actions = {
	getPasswordProtectedShare: async ({ request }) => {
		const data = await request.formData();
		let id = data.get('id') as string;
		id = id.substring(1);
		const params: GetPasswordProtectedShareRequest = {
			password: data.get('password') as string
		};
		let share: Share;
		try {
			share = await getPasswordProtectedShare(id, params);
		} catch (err) {
			if (err instanceof Error) {
				return fail(err instanceof WrongPasswordError ? 400 : 500, { message: err.message });
			}
			return fail(500, { message: 'Unknown server error' });
		}
		return { share };
	},
	deleteShare: async ({ request, locals }) => {
		const data = await request.formData();
		const shareId = data.get('shareId') as string;
		const sessionId = locals.sessionId ?? '';

		try {
			await deleteShare(shareId, sessionId);
		} catch (err) {
			if (err instanceof Error) {
				return fail(500, { message: err.message });
			}
			return fail(500, { message: 'Unexpected server error' });
		}
	},
	translate: async ({ fetch, request }) => {
		const url = `https://translation.googleapis.com/language/translate/v2?key=${GOOGLE_API_KEY}`;

		const data = await request.formData();
		const text = data.get('content') as string;
		const language = data.get('language') as string;

		try {
			const response = await fetch(url, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					q: text,
					target: language,
					format: 'text'
				})
			});
			if (!response.ok) {
				const errorMsg = await response.json().catch(() => response.statusText);
				return fail(500, { message: errorMsg });
			}
			const parsedResponse = await response.json();
			return parsedResponse['data']['translations'][0]['translatedText'];
		} catch (err) {
			return fail(500, { message: JSON.stringify(err) });
		}
	}
};
