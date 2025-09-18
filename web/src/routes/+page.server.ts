import { createShare, type Share } from '$lib/share.js';
import { fail } from '@sveltejs/kit';

export const actions = {
	createShare: async ({ request }) => {
		const data = await request.formData();
		const params: Share = {
			title: data.get('title') as string,
			content: data.get('content') as string,
			expireIn: data.get('expireIn') as string,
			password: data.get('password') as string
		};
		try {
			const response = await createShare(params);
			return { response };
		} catch (err) {
			return fail(500, {
				message: err instanceof Error ? err.message : 'Unknown error'
			});
		}
	}
};
