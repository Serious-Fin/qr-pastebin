import { deleteShare, getSharesForUser } from '$lib/share';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	const userId = locals.user?.id ?? -1;
	const shares = await getSharesForUser(locals.sessionId ?? '');
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
				return {
					errMsg: err.message
				};
			}
		}
	}
};
