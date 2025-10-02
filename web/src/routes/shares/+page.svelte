<script lang="ts">
	import type { PageProps } from './$types';
	import { truncateString } from '$lib/helpers';

	let { data }: PageProps = $props();
</script>

<section id="main">
	<h1>Shareit</h1>
	<h2>You shares</h2>
	{#if data.shares.length === 0}
		<p>No shares, try creating one</p>
	{/if}
	{#each data.shares as share}
		<div id="share-box">
			{#if share.title}
				<p class="title">{share.title}</p>
			{/if}
			<p class="content">{truncateString(share.content, 100)}</p>

			<div class="buttons">
				<button class="button">Edit</button>
				<button class="button delete">Delete</button>
			</div>
		</div>
	{/each}
</section>

<style>
	#share-box {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 100%;
		background-color: var(--light);
		margin-bottom: 40px;
		box-shadow: 2px 2px 2px 2px rgba(0, 0, 0, 0.153);
		border-radius: 5px;
		padding: 20px 10px;
	}

	.title {
		font-weight: 600;
		align-self: center;
		margin-bottom: 15px;
	}

	.content {
		background-color: var(--lightest);
		padding: 10px;
		border-radius: 5px;
		margin-bottom: 15px;
		width: 100%;
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

	.buttons {
		display: flex;
		gap: 30px;
	}

	.delete {
		background-color: var(--red);
	}

	.button:active {
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.2);
		transform: translateY(2px);
	}

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
		margin-bottom: 20px;
		font-weight: 500;
	}

	h2 {
		font-size: 15pt;
		font-weight: 500;
		margin-bottom: 20px;
	}
</style>
