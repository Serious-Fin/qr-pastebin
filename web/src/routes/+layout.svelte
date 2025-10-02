<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();

	const logout = () => {};
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div id="header-box">
	<div id="login-box">
		{#if data.userId !== -1}
			<button class="button" onclick={() => (window.location.href = '/api/logout')}>Logout</button>
			<p>{data.username}</p>
		{:else}
			<a href="/login">Login</a>
			<a href="/signup">Sign-up</a>
		{/if}
	</div>

	{#if data.userId !== -1}
		<button class="button">My shares</button>
	{/if}
</div>

{@render children?.()}

<style>
	#header-box {
		background-color: var(--light);
		display: flex;
		justify-content: space-between;
		padding: 10px 15px;
	}

	#login-box {
		display: flex;
		align-items: center;
		gap: 15px;
	}

	.button {
		padding: 10px 10px;
		font-size: 12pt;
		background-color: var(--accent);
		border: none;
		border-radius: 10px;
		box-shadow: 2px 2px 2px 2px rgba(0, 0, 0, 0.2);
	}

	.button:active {
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.199);
		transform: translateY(2px);
	}

	a {
		padding: 10px 3px;
	}
</style>
