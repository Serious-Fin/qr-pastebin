import type { PageServerLoad } from './$types';
import { getShare } from '$lib/share';

export const load: PageServerLoad = async ({ params }) => {
	return {
		share: getShare(params.id)
	};
};
