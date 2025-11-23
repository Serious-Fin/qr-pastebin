import type { RequestHandler } from './$types';
import { error, redirect } from '@sveltejs/kit';
import { v4 as uuidv4 } from 'uuid';
import { PUBLIC_GH_OAUTH_CLIENT_ID } from '$env/static/public';
import { GH_OAUTH_CLIENT_SECRET } from '$env/static/private';
import { checkIfGithubUserExists, createNewUser, tryCreateSessionForUser } from '$lib/user';

const ghStateCookieName = 'gh_state';

interface GithubUser {
	id: number;
	name: string;
}

export const GET: RequestHandler = async ({ url, cookies }) => {
	const redirectUri = `https://localhost:5173/api/github`;
	if (url.searchParams.has('code')) {
		const storedState = cookies.get(ghStateCookieName);
		const returnedState = url.searchParams.get('state');
		if (!storedState || storedState !== returnedState) {
			error(401, { message: 'Error logging in, try again: XSS attack detected' });
		}

		let sessionId: string;
		try {
			const code = url.searchParams.get('code') ?? '';
			const accessToken = await getAccessToken(code, redirectUri);
			const user = await getUserInfo(accessToken);

			const userExists = await checkIfGithubUserExists(user.id);
			console.log(`userExists: ${userExists}`);
			if (!userExists) {
				await createNewUser({
					id: user.id,
					name: user.name,
					password: '',
					isOauth: true
				});
				console.log(`user created`);
			}
			sessionId = await tryCreateSessionForUser({
				id: user.id,
				name: user.name,
				password: '',
				isOauth: true
			});
			console.log(`Session id: ${sessionId}`);
		} catch (err) {
			error(401, { message: `Error logging in, try again: ${JSON.stringify(err)}` });
		}

		cookies.set('session', sessionId, {
			httpOnly: true,
			sameSite: 'lax',
			path: '/',
			maxAge: 60 * 60 * 24 * 7
		});

		const redirectTo = url.searchParams.get('redirectTo');
		const decodedRedirectTo = redirectTo ? decodeURIComponent(redirectTo) : '/';
		throw redirect(302, `${decodedRedirectTo}`);
	} else {
		const scopes = 'read:user user:email';
		const state = uuidv4();
		const githubOAuthUrl =
			`https://github.com/login/oauth/authorize` +
			`?client_id=${encodeURIComponent(PUBLIC_GH_OAUTH_CLIENT_ID)}` +
			`&redirect_uri=${encodeURIComponent(redirectUri)}` +
			`&scopes=${encodeURIComponent(scopes)}` +
			`&state=${encodeURIComponent(state)}`;

		cookies.set(ghStateCookieName, state, {
			httpOnly: true,
			sameSite: 'lax',
			path: '/',
			maxAge: 60 * 10
		});
		redirect(307, githubOAuthUrl);
	}
};

async function getAccessToken(code: string, redirectUri: string): Promise<string> {
	try {
		const response = await fetch('https://github.com/login/oauth/access_token', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json'
			},
			body: JSON.stringify({
				client_id: PUBLIC_GH_OAUTH_CLIENT_ID,
				client_secret: GH_OAUTH_CLIENT_SECRET,
				code,
				redirect_uri: redirectUri
			})
		});
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`Error receiving github access token ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		const { access_token } = await response.json();
		return access_token;
	} catch (err) {
		throw new Error(`github access_token request failed: ${err}`);
	}
}

async function getUserInfo(accessToken: string): Promise<GithubUser> {
	try {
		const response = await fetch('https://api.github.com/user', {
			headers: {
				Authorization: `Bearer ${accessToken}`
			}
		});
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }));
			throw new Error(
				`response was non 2xx ${response.status} - ${errorBody.message || 'Unknown error'}`
			);
		}
		const body = await response.json();
		return {
			id: body.id,
			name: body.login
		};
	} catch (err) {
		throw new Error(`could not get github user details: ${err}`);
	}
}
