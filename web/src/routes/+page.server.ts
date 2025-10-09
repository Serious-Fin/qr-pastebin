import { createShare, type ShareRequest } from '$lib/share.js';
import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	return {
		userId: locals.user?.id ?? -1,
		username: locals.user?.name ?? 'Anon'
	};
};

export const actions = {
	createShare: async ({ request }) => {
		const data = await request.formData();
		const title = data.get('title') ? (data.get('title') as string) : '';
		const content = data.get('content') ? (data.get('content') as string) : '';
		const setPassword = data.get('setPassword') !== null;
		const password = data.get('password') ? (data.get('password') as string) : '';
		const expireIn = data.get('expireIn') as string;
		const hideAuthor = data.get('hideAuthor') !== null;
		const authorId = parseInt((data.get('userId') as string) ?? '-1');

		if (setPassword && password == '') {
			return fail(400, {
				message: 'Can not set empty password. Uncheck "Set Password" for no password'
			});
		}

		const params: ShareRequest = {
			title,
			content,
			setPassword,
			password,
			expireIn,
			hideAuthor,
			authorId
		};
		let newShareId = '';
		try {
			newShareId = await createShare(params);
		} catch (err) {
			return fail(500, {
				message: err instanceof Error ? err.message : 'Unknown error'
			});
		}
		throw redirect(303, `/${newShareId}`);
	}
};
