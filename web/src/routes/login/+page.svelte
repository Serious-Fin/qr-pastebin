<script lang="ts">
	import { enhance } from '$app/forms';
	import type { SubmitFunction } from '@sveltejs/kit';

	let name = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let err = $state('');

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
				} else if (result.type === 'failure') {
					throw Error(result.data?.message || 'Unknown server error occurred');
				} else {
					throw Error('Could not query agent');
				}
			} catch (err) {
				if (err instanceof Error) {
					console.error('BOOM ERROR');
					return;
				}
			} finally {
				isLoading = false;
			}
		};
	};
</script>

<section id="main">
	<h1>Shareit</h1>
	<h2>Login</h2>
	<form method="POST" action="?/login" use:enhance={handleLoginErrors}>
		<label for="name">Name</label>
		<input type="text" id="name" name="name" bind:value={name} />
		<label for="password">Password</label>
		<input type="password" id="password" name="password" bind:value={password} />
		<p>Don't have an account? <a href="/signup">Sign-up</a></p>
		<input type="submit" value="Login" />
	</form>

	{#if err}
		<p id="err">{err}</p>
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
</style>
