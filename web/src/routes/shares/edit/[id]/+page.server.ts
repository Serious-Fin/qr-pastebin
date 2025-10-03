import type { PageServerLoad } from './$types';
import { getShareForEdit } from '$lib/share';

export const load: PageServerLoad = async ({ locals, params }) => {
	const share = await getShareForEdit(params.id, locals.sessionId ?? '');
	return {
		userId: locals.user?.id ?? -1,
		username: locals.user?.name ?? 'Anon',
		share
	};
};
