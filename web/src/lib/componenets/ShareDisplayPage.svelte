<script lang="ts">
	import QRCode from 'qrcode';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import type { Share } from '$lib/share';

	let { share }: { share: Share } = $props();
	let svg: string = $state('');

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
</style>
