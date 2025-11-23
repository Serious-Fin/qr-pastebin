<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { logError, logSuccess, fetchQuote } from '$lib/helpers';
	import LoadingSpinner from '$lib/componenets/LoadingSpinner.svelte';
	import { onMount } from 'svelte';

	let { data }: PageProps = $props();
	let userId = $state(data.userId);
	let setPassword = $state(false);
	let isLoading = $state(false);

	let placeholderQuote = $state('');
	onMount(async () => {
		placeholderQuote = await fetchQuote();
	});

	const handleShareCreation: SubmitFunction = () => {
		isLoading = true;
		return async ({ update, result }) => {
			try {
				if (result.type === 'failure') {
					throw Error(result.data?.message || 'Error signing in, try again');
				} else {
					logSuccess('Share created');
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
	<h1>Shareit</h1>
	<h2>New share</h2>
	<form method="POST" action="?/createShare" use:enhance={handleShareCreation}>
		<textarea id="content" name="content" required>{placeholderQuote}</textarea>

		<div id="grid">
			<label for="title">Title:</label>
			<input type="text" class="property-input" name="title" id="title" />

			<label for="expireIn">Expire in:</label>
			<select class="property-input" name="expireIn" id="expireIn">
				<option value="never" selected>Never</option>
				<option value="1_minutes">1 minute</option>
				<option value="10_minutes">10 minutes</option>
				<option value="1_hours">1 hour</option>
				<option value="1_days">1 day</option>
				<option value="1_weeks">1 week</option>
				<option value="2_weeks">2 weeks</option>
				<option value="1_months">1 month</option>
				<option value="6_months">6 months</option>
				<option value="1_years">1 year</option>
			</select>

			<label for="setPassword">Set password</label>
			<input type="checkbox" id="setPassword" name="setPassword" bind:checked={setPassword} />

			<label for="password">Password:</label>
			<input
				type="password"
				class="property-input"
				name="password"
				id="password"
				disabled={!setPassword}
			/>

			<label for="hideAuthor">Hide author</label>
			<input type="checkbox" id="hideAuthor" name="hideAuthor" />
		</div>

		<input type="hidden" id="userId" name="userId" bind:value={userId} />

		{#if isLoading}
			<div id="loadingBox">
				<LoadingSpinner />
			</div>
		{:else}
			<input type="submit" value="Create" />
		{/if}
	</form>
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

	#loadingBox {
		margin-top: 30px;
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
		width: 100%;
	}

	#content {
		width: 100%;
		max-width: 90vw;
		min-height: 150px;
		margin: 15px 0;
		padding: 7px 8px;
		border-radius: 10px;
	}

	#grid {
		display: grid;
		grid-template-columns: 100px auto;
		row-gap: 10px;
	}

	input,
	select {
		padding: 3px 5px;
	}

	input[type='checkbox'] {
		justify-self: baseline;
		width: 20px;
	}

	input[type='submit'] {
		padding: 10px 20px;
		background-color: var(--accent);
		border: none;
		box-shadow: 2px 2px 2px 2px rgba(0, 0, 0, 0.2);
		border-radius: 10px;
		margin-top: 20px;
		font-size: 12pt;
	}

	input[type='submit']:active {
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.199);
		transform: translateY(2px);
	}

	@media (min-width: 768px) {
		#main {
			padding: 30px 50px;
		}

		#content {
			margin-bottom: 30px;
		}

		#grid {
			row-gap: 15px;
			grid-template-columns: 170px auto;
		}

		input,
		select {
			padding: 5px 10px;
			font-size: 11pt;
		}

		input[type='checkbox'] {
			width: 25px;
		}
	}

	@media (min-width: 1024px) {
		#grid {
			row-gap: 20px;
			grid-template-columns: 170px auto;
		}

		input[type='submit'] {
			margin-top: 30px;
		}
	}
</style>
