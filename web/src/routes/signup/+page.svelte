<script lang="ts">
	import PasswordRequirements from '$lib/componenets/PasswordRequirements.svelte';

	let password = $state('');
	let passwordMeetsCriteria = $state(false);
	const lengthMet = $derived(password.length >= 8);
	const uppercaseMet = $derived(/[A-Z]/.test(password));
	const lowercaseMet = $derived(/[a-z]/.test(password));
	const numberMet = $derived(/[0-9]/.test(password));
	const symbolMet = $derived(/[!@#$%^&*()_+\-_<>,\.{}:;'"|]/.test(password));

	const checkState = () => {
		if (lengthMet && uppercaseMet && lowercaseMet && numberMet && symbolMet) {
			passwordMeetsCriteria = true;
		} else {
			passwordMeetsCriteria = false;
		}
	};
</script>

<section id="main">
	<h1>Shareit</h1>
	<h2>Sign-up</h2>
	<form>
		<label for="name">Name</label>
		<input type="text" id="name" name="name" />
		<label for="password">Password</label>
		<input
			type="password"
			id="password"
			name="password"
			bind:value={password}
			onchange={checkState}
		/>

		<PasswordRequirements {password} />

		<p id="redirect">Already have an account? <a href="/login">Login</a></p>
		<input type="submit" value="Sign-up" disabled={!passwordMeetsCriteria} />
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

	#redirect {
		font-size: 12pt;
		margin-top: 15px;
		align-self: baseline;
	}
</style>
