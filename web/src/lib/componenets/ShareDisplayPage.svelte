<script lang="ts">
	import QRCode from 'qrcode';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import type { Share } from '$lib/share';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import { enhance } from '$app/forms';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { logError, logSuccess } from '$lib/helpers';

	let { share, isAdmin }: { share: Share; isAdmin: boolean } = $props();
	let svg: string = $state('');
	let isLoading = $state(false);

	const handleDeleteSubmissionError: SubmitFunction = () => {
		return async ({ update, result }) => {
			isLoading = true;
			try {
				await update();
				if (result.type === 'success') {
					logSuccess('Share deleted');
					window.location.href = `/${share.id}`;
				} else if (result.type === 'failure') {
					throw Error(result.data?.message || 'Unknown server error occurred');
				}
			} catch (err) {
				if (err instanceof Error) {
					logError('Error deleting share, try again later', err);
				}
			} finally {
				isLoading = false;
			}
		};
	};

	onMount(async () => {
		svg = await QRCode.toString(page.url.href, {
			type: 'svg',
			errorCorrectionLevel: 'L',
			margin: 2,
			width: 200
		});
	});
</script>

{#if share.title !== undefined && share.title !== ''}
	<h2>{share.title}</h2>
{/if}
<textarea id="content" name="content" readonly>{share.content}</textarea>
<div id="share-info">
	{#if share.authorName}
		<p>Created by: {share.authorName}</p>
	{/if}

	{#if share.expiresIn}
		<p>{share.expiresIn}</p>
	{:else}
		<p>Does not expire</p>
	{/if}

	{#if share.isPasswordProtected}
		<p>Is password protected</p>
	{:else}
		<p>Not password protected</p>
	{/if}
</div>

<div class="qr">{@html svg}</div>

{#if isAdmin}
	<form method="POST" action="?/deleteShare" use:enhance={handleDeleteSubmissionError}>
		{#if isLoading}
			<LoadingSpinner />
		{:else}
			<input type="hidden" id="shareId" name="shareId" value={share.id} />
			<input type="submit" value="Delete" class="button delete" />
		{/if}
	</form>
{/if}

<style>
	p {
		color: rgb(65, 65, 65);
		font-size: 12pt;
	}

	#content {
		width: 100%;
		max-width: 100%; /* so text box couldn't expand out of view */
		min-height: 150px;
		height: fit-content;
		margin: 15px 0;
		padding: 7px 8px;
		border-radius: 5px;
		font-size: 12pt;
	}

	#share-info {
		width: 100%;
		margin-bottom: 40px;
		box-sizing: border-box;
		padding-left: 10px;

		display: flex;
		flex-direction: column;
	}

	.qr {
		width: 60%;
		height: fit-content;
		display: flex;
		justify-content: center;
	}

	.button {
		padding: 10px 20px;
		font-size: 12pt;
		background-color: var(--accent);
		border: none;
		border-radius: 10px;
		box-shadow: 2px 2px 2px 2px rgba(0, 0, 0, 0.2);
		width: fit-content;
	}

	.delete {
		background-color: var(--red);
		margin-top: 30px;
	}

	.button:active {
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.2);
		transform: translateY(2px);
	}
</style>
