<script lang="ts">
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
	let userId = $state(data.userId);
	let share = $state(data.share);
</script>

<section id="main">
	<h1>Shareit</h1>
	<h2>Edit share</h2>
	<form method="POST" action="?/createShare">
		<textarea id="content" name="content" required>{share.content}</textarea>

		<div id="grid">
			<label for="title">Title:</label>
			<input type="text" class="property-input" name="title" id="title" bind:value={share.title} />

			<label for="expireIn">Expire in:</label>
			<select class="property-input" name="expireIn" id="expireIn">
				<option value="no-change" selected>Don't change</option>
				<option value="">Never</option>
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

			<label for="password">New password:</label>
			<input
				type="password"
				class="property-input"
				name="password"
				id="password"
				placeholder="Leave blank for no change"
			/>

			<label for="hideAuthor">Hide author?</label>
			<input type="checkbox" id="hideAuthor" name="hideAuthor" bind:checked={share.hideAuthor} />
		</div>

		<input type="hidden" id="userId" name="userId" bind:value={userId} />

		<input type="submit" value="Save changes" />
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
		grid-template-columns: 120px auto;
		row-gap: 10px;
		min-width: 90%;
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
</style>
