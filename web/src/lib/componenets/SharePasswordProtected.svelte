<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import { FetchShareStatus, type Share } from '$lib/share';
	import type { SubmitFunction } from '@sveltejs/kit';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import { logSuccess, logError } from '$lib/helpers';

	let { updateShare }: { updateShare: (share: Share, status: FetchShareStatus) => void } = $props();

	let isLoading = $state(false);

	const showShare: SubmitFunction = () => {
		isLoading = true;
		return async ({ update, result }) => {
			try {
				if (result.type === 'success' && result.data) {
					const share = result.data.share;
					const status = FetchShareStatus.Accessible;
					updateShare(share, status);
					await update();
					logSuccess('Password correct, share unlocked');
				} else if (result.type === 'failure') {
					throw Error(result.data?.message || 'Error unlocking share, try again');
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

<h2>Share is password protected</h2>

<form method="POST" action="?/getPasswordProtectedShare" use:enhance={showShare}>
	<label for="password"
		>Enter password
		<input type="hidden" name="id" value={page.url.pathname} />
		<input type="password" name="password" id="password" required />
	</label>
	{#if isLoading}
		<div id="loadingBox">
			<LoadingSpinner />
		</div>
	{:else}
		<input type="submit" value="Access" />
	{/if}
</form>

<style>
	h2 {
		font-size: 13pt;
		color: rgb(29, 29, 29);
		margin-bottom: 20px;
	}

	form {
		display: flex;
		flex-direction: column;
		width: 100%;
		align-items: center;
	}

	input {
		padding: 3px 5px;
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

	#loadingBox {
		margin-top: 20px;
	}

	@media (min-width: 768px) {
		h2 {
			margin-bottom: 70px;
		}

		input[type='password'] {
			width: 250px;
			margin-left: 15px;
		}

		input[type='submit'] {
			margin-top: 30px;
		}
	}

	@media (min-width: 1024px) {
		h1 {
			font-size: 22pt;
		}
	}
</style>
