import type { PageServerLoad, Actions } from './$types';
import { getShareForEdit, type EditShareRequest, editShare } from '$lib/share';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals, params }) => {
	const share = await getShareForEdit(params.id, locals.sessionId ?? '');
	return {
		userId: locals.user?.id ?? -1,
		username: locals.user?.name ?? 'Anon',
		share
	};
};

export const actions: Actions = {
	editShare: async ({ request, locals }) => {
		const data = await request.formData();
		const shareId = data.get('shareId') as string;
		const sessionId = locals.sessionId ?? '';

		const editShareBody: EditShareRequest = {
			title: data.get('title') as string,
			content: data.get('content') as string,
			setPassword: data.get('setPassword') !== null,
			password: data.get('password') as string,
			expireIn: data.get('expireIn') as string,
			hideAuthor: data.get('shareId') !== null
		};

		try {
			await editShare(editShareBody, sessionId, shareId);
		} catch (err) {
			if (err instanceof Error) {
				return fail(400, { message: err.message });
			}
			return fail(500, { message: 'Unexpected server error' });
		}
	}
};
