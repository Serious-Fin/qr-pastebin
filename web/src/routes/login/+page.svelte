<script lang="ts">
	import { enhance } from '$app/forms';
	import type { SubmitFunction } from '@sveltejs/kit';
	import type { PageProps } from './$types';
	import { page } from '$app/state';

	let { data }: PageProps = $props();

	let name = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let err = $state('');
	let redirectTo = $state(page.url.searchParams.get('redirectTo') ?? '/');

	const handleLoginErrors: SubmitFunction = () => {
		isLoading = true;
		return async ({ update, result }) => {
			try {
				await update();
				if (result.type === 'success') {
					if (result.data?.errMsg) {
						err = result.data.errMsg;
						name = result.data.user.name;
						password = result.data.user.password;
					}
				} else if (result.type === 'redirect') {
					return;
				} else if (result.type === 'failure') {
					throw Error(result.data?.message || 'Unknown server error occurred');
				} else {
					throw Error('Could not query agent');
				}
			} catch (err) {
				if (err instanceof Error) {
					console.error(`Error logging in: ${JSON.stringify(err)}`);
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
		<h2>Login</h2>
		<form method="POST" action="?/login" use:enhance={handleLoginErrors}>
			<label for="name">Name</label>
			<input type="text" id="name" name="name" bind:value={name} />
			<label for="password">Password</label>
			<input type="password" id="password" name="password" bind:value={password} />
			<p>Don't have an account? <a href="/signup">Sign-up</a></p>
			<input type="submit" value="Login" />
			{#if err}
				<p id="err">{err}</p>
			{/if}

			<input type="hidden" id="redirectTo" name="redirectTo" bind:value={redirectTo} />
		</form>
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

	#err {
		color: red;
		margin-top: 30px;
		align-self: center;
	}
</style>
