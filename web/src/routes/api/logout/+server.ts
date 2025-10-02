import { redirect } from '@sveltejs/kit';
import { page } from '$app/state';

export async function GET({ cookies, url }) {
	cookies.delete('session', { path: '/' });
	const redirectTo = url.searchParams.get('redirectTo') ?? '/';
	throw redirect(303, redirectTo);
}
