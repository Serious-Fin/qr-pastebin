import { getSharesForUser } from '$lib/share';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	const userId = locals.user?.id ?? -1;
	const shares = await getSharesForUser(`${userId}`);
	return {
		userId,
		username: locals.user?.name ?? 'Anon',
		shares
	};
};
