<script lang="ts">
	import { enhance } from '$app/forms';
	import type { SubmitFunction } from '@sveltejs/kit';
	import type { PageProps } from './$types';
	import { page } from '$app/state';
	import LoadingSpinner from '$lib/componenets/LoadingSpinner.svelte';
	import { logError, logSuccess } from '$lib/helpers';

	let { data }: PageProps = $props();

	let name = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let redirectTo = $state(page.url.searchParams.get('redirectTo') ?? '/');

	const handleLoginErrors: SubmitFunction = () => {
		isLoading = true;
		return async ({ update, result }) => {
			try {
				if (result.type === 'failure') {
					throw Error(result.data?.message || 'Error logging in, try again');
				} else {
					logSuccess('Logged in successfully');
					await update();
				}
			} catch (err) {
				if (err instanceof Error) {
					if (result?.status && result.status >= 500) {
						logError('Server error', err);
					} else {
						logError(err.message, err);
					}
				}
			} finally {
				isLoading = false;
			}
		};
	};

	const loginViaGithub = () => {
		window.location.href = '/api/github';
	};
</script>

<section id="main">
	{#if data.userId === -1}
		<h1>Shareit</h1>
		<h2>Login</h2>
		<form method="POST" action="?/login" use:enhance={handleLoginErrors}>
			<label for="name">Name</label>
			<input type="text" id="name" name="name" bind:value={name} required />
			<label for="password">Password</label>
			<input type="password" id="password" name="password" bind:value={password} required />
			<p>Don't have an account? <a href="/signup">Sign-up</a></p>

			{#if isLoading}
				<div id="loadingBox">
					<LoadingSpinner />
				</div>
			{:else}
				<input type="submit" value="Login" />
			{/if}

			<input type="hidden" id="redirectTo" name="redirectTo" bind:value={redirectTo} />
		</form>

		<button class="github_login" onclick={loginViaGithub}>
			<img src="/github-mark.svg" alt="github logo" />Sign in with GitHub
		</button>
	{:else}
		<p id="already-have-acc">Already logged in as {data.username}</p>{/if}
</section>

<style>
	#main {
		background-color: var(--lightest);
		color: var(--black);

		display: flex;
		flex-direction: column;
		align-items: center;

		width: 100vw;
		height: 100vh;
		padding: 25px 10px;
	}

	h1 {
		font-size: 22pt;
		margin-bottom: 25px;
		font-weight: 500;
	}

	h2 {
		font-size: 15pt;
		font-weight: 500;
	}

	form {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 80%;
		margin-top: 20px;
	}

	label {
		align-self: baseline;
		margin-top: 20px;
	}

	input[type='text'],
	input[type='password'] {
		padding: 5px 7px;
		width: 100%;
	}

	input[type='submit'] {
		padding: 10px 20px;
		background-color: var(--accent);
		border: none;
		box-shadow: 2px 2px 2px 2px rgba(0, 0, 0, 0.199);
		border-radius: 10px;
		margin-top: 20px;
		font-size: 12pt;
	}

	input[type='submit']:active {
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.199);
		transform: translateY(2px);
	}

	p {
		font-size: 12pt;
		margin-top: 15px;
		align-self: baseline;
	}

	#already-have-acc {
		margin-top: 40px;
	}

	#loadingBox {
		margin-top: 20px;
	}

	.github_login {
		background-color: white;
		padding: 10px 14px;
		margin-top: 40px;
		border-radius: 4px;
		border: 0.6px solid rgb(218, 220, 224);
		font-size: 16px;
		letter-spacing: 0.25px;
		font-family: 'Google Sans', arial, sans-serif;
		font-weight: 500;
		color: #3c4043f2;

		display: flex;
		align-items: center;
		gap: 10px;
	}

	.github_login img {
		width: 18px;
	}

	.github_login:active {
		background-color: rgb(236, 243, 254);
		border: 2px solid rgb(0, 99, 155);
	}

	@media (min-width: 768px) {
		#main {
			padding: 30px 50px;
		}

		form {
			max-width: 350px;
		}

		label {
			margin-bottom: 5px;
		}

		p {
			margin-top: 30px;
		}
	}
</style>
