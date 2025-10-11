<script lang="ts">
	import type { PageProps } from './$types';

	import { FetchShareStatus, type Share } from '$lib/share';
	import ShareNotFound from '$lib/componenets/ShareNotFound.svelte';
	import SharePasswordProtected from '$lib/componenets/SharePasswordProtected.svelte';
	import ShareDisplayPage from '$lib/componenets/ShareDisplayPage.svelte';

	let { data }: PageProps = $props();
	let share = $state(data.share);
	let status = $state(data.status);
	let isAdmin = data.isAdmin ?? false;

	const updateShare = (newShare: Share, newStatus: FetchShareStatus) => {
		share = newShare;
		status = newStatus;
	};
</script>

<section id="main">
	<h1>Shareit</h1>
	{#if status === FetchShareStatus.Accessible && share}
		<ShareDisplayPage {isAdmin} {share} />
	{:else if status === FetchShareStatus.NotFound}
		<ShareNotFound />
	{:else if status === FetchShareStatus.NeedPassword}
		<SharePasswordProtected {updateShare} />
	{/if}
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

	@media (min-width: 768px) {
		#main {
			padding: 30px 50px;
		}
	}
</style>
