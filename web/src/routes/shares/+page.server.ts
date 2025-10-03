import { deleteShare, getSharesForUser, type Share } from '$lib/share';
import { fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	const userId = locals.user?.id ?? -1;
	let shares: Share[];
	try {
		shares = await getSharesForUser(locals.sessionId ?? '');
	} catch {
		shares = [];
	}

	return {
		userId,
		username: locals.user?.name ?? 'Anon',
		shares
	};
};

export const actions: Actions = {
	deleteShare: async ({ request, locals }) => {
		const data = await request.formData();
		const shareId = data.get('shareId') as string;
		const sessionId = locals.sessionId ?? '';

		try {
			await deleteShare(shareId, sessionId);
		} catch (err) {
			if (err instanceof Error) {
				return fail(400, { message: err.message });
			}
			return fail(500, { message: 'Unexpected server error' });
		}
	}
};
