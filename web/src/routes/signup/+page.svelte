<script lang="ts">
	import PasswordField from '$lib/componenets/PasswordField.svelte';
	import { enhance } from '$app/forms';
	import type { SubmitFunction } from '@sveltejs/kit';
	import type { PageProps } from './$types';
	import { page } from '$app/state';
	import LoadingSpinner from '$lib/componenets/LoadingSpinner.svelte';
	import { logError, logSuccess } from '$lib/helpers';

	let { data }: PageProps = $props();

	let name = $state('');
	let password = $state('');
	let passwordMeetsCriteria = $state(false);
	let isLoading = $state(false);
	let err = $state('');
	let redirectTo = $state(page.url.searchParams.get('redirectTo') ?? '/');

	const updatePasswordMeetsCriteria = (newState: boolean) => {
		passwordMeetsCriteria = newState;
	};

	const handleUserCreation: SubmitFunction = () => {
		isLoading = true;
		return async ({ update, result }) => {
			try {
				if (result.type === 'failure') {
					throw Error(result.data?.message || 'Error signing in, try again');
				} else {
					logSuccess('Account created successfully');
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
</script>

<section id="main">
	{#if data.userId === -1}
		<h1>Shareit</h1>
		<h2>Sign-up</h2>
		<form method="POST" action="?/createNewUser" use:enhance={handleUserCreation}>
			<label for="name">Name</label>
			<input type="text" id="name" name="name" bind:value={name} required />

			<PasswordField {password} {updatePasswordMeetsCriteria} />

			<p id="redirect">Already have an account? <a href="/login">Login</a></p>
			{#if isLoading}
				<div id="loadingBox">
					<LoadingSpinner />
				</div>
			{:else}
				<input type="submit" value="Sign-up" disabled={!passwordMeetsCriteria} />
			{/if}

			<input type="hidden" id="redirectTo" name="redirectTo" bind:value={redirectTo} />
		</form>
	{:else}
		<p id="already-have-acc">Already signed in as {data.username}</p>
	{/if}
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

	input[type='text'] {
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

	#redirect {
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
</style>
