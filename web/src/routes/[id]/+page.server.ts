import type { PageServerLoad } from './$types';
import { getShare, isSharePasswordProtected, FetchShareStatus, type Share, type GetPasswordProtectedShareRequest, getPasswordProtectedShare } from '$lib/share';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params }) => {
	// Load "Enter password" view if share is password protected
	let hasPassword: boolean
	try {
		hasPassword = await isSharePasswordProtected(params.id)
	} catch {
		return {
			status: FetchShareStatus.NotFound
		}
	}
	if (hasPassword) {
		return {
			status: FetchShareStatus.NeedPassword
		}
	}

	// GET share if it's not password protected
	let share: Share
	try {
		share = await getShare(params.id)
	} catch {
		return {
			status: FetchShareStatus.NotFound
		}
	}
	return {
		share: share,
		status: FetchShareStatus.Accessible
	};
};

export const actions = {
	getPasswordProtectedShare: async ({ request }) => {
		const data = await request.formData();
		let id = data.get("id") as string
		id = id.substring(1)
		const params: GetPasswordProtectedShareRequest = {
			password: data.get('password') as string
		};
		let share: Share
		try {
			share = await getPasswordProtectedShare(id, params);
		} catch (err) {
			return fail(500, {
				message: err instanceof Error ? err.message : 'Unknown error'
			});
		}
		return {share}
	}
};
