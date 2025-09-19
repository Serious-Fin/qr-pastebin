import { createShare, type CreateShareRequest } from '$lib/share.js';
import { fail, redirect } from '@sveltejs/kit';

export const actions = {
	createShare: async ({ request }) => {
		const data = await request.formData();
		const params: CreateShareRequest = {
			title: data.get('title') as string,
			content: data.get('content') as string,
			expireIn: data.get('expireIn') as string,
			password: data.get('password') as string
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
