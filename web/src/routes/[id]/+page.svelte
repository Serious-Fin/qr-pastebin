<script lang="ts">
	import type { PageProps } from './$types';
	import QRCode from 'qrcode';
	import { onMount } from 'svelte';
	import { page } from '$app/state';

	let { data }: PageProps = $props();
	const share = data.share;

	let svg: string = $state('');

	onMount(async () => {
		svg = await QRCode.toString(page.url.href, {
			type: 'svg',
			errorCorrectionLevel: 'L',
			margin: 2,
			width: 300
		});
	});
</script>

<section id="main">
	<h1>Shareit</h1>
	{#if share.title !== undefined && share.title !== ''}
		<h2>{share.title}</h2>
	{/if}
	<textarea id="content" name="content" readonly>{share.content}</textarea>
	<div id="share-info">
		{#if share.author}
			<p>Created by: {share.author}</p>
		{/if}

		{#if share.expiresIn}
			<p>Expires in {share.expiresIn}</p>
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
</section>

<style>
	#main {
		background-color: var(--lightest);
		color: var(--black);

		display: flex;
		flex-direction: column;
		align-items: center;
		box-sizing: border-box;

		width: 100vw;
		height: 100vh;
		padding: 25px 10px;
	}

	h1 {
		font-size: 22pt;
		margin-bottom: 25px;
		font-weight: 500;
	}

	p {
		color: rgb(56, 56, 56);
		font-size: 11pt;
	}

	#content {
		width: 100%;
		max-width: 100%; /* so text box couldn't expand out of view */
		min-height: 150px;
		height: fit-content;
		margin: 15px 0;
		padding: 7px 8px;
		border-radius: 5px;
	}

	#share-info {
		width: 100%;
		margin-bottom: 20px;
		box-sizing: border-box;
		padding-left: 10px;

		display: flex;
		flex-direction: column;
	}

	.qr {
		width: 100%;
		height: fit-content;
		display: flex;
		justify-content: center;
	}
</style>
