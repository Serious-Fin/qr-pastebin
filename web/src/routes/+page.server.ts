import { createShare, type CreateShareRequest } from '$lib/share.js';
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
		const authorId = parseInt((data.get('userId') as string) ?? '-1');
		const hideAuthor = data.get('hideAuthor') !== null;
		const params: CreateShareRequest = {
			title: data.get('title') as string,
			content: data.get('content') as string,
			expireIn: data.get('expireIn') as string,
			password: data.get('password') as string,
			authorId,
			hideAuthor
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
