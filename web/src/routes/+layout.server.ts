import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	return {
		userId: locals.user?.id ?? -1,
		username: locals.user?.name ?? 'Anon'
	};
};
