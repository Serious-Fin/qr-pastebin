<script lang="ts">
	let {
		password,
		updatePasswordMeetsCriteria
	}: { password: string; updatePasswordMeetsCriteria: (newState: boolean) => void } = $props();
	const lengthMet = $derived(password.length >= 8);
	const uppercaseMet = $derived(/[A-Z]/.test(password));
	const lowercaseMet = $derived(/[a-z]/.test(password));
	const numberMet = $derived(/[0-9]/.test(password));
	const symbolMet = $derived(/[!@#$%^&*()_+\-_<>,\.{}:;'"|]/.test(password));

	const checkState = () => {
		if (lengthMet && uppercaseMet && lowercaseMet && numberMet && symbolMet) {
			updatePasswordMeetsCriteria(true);
		} else {
			updatePasswordMeetsCriteria(false);
		}
	};
</script>

<label for="password">Password</label>
<input
	type="password"
	id="password"
	name="password"
	bind:value={password}
	oninput={checkState}
	required
/>

<article id="password-tips">
	<p id="length" class="password-tip" class:met={lengthMet} class:not-met={!lengthMet}>
		At least 8 characters long
	</p>
	<p id="uppercase" class="password-tip" class:met={uppercaseMet} class:not-met={!uppercaseMet}>
		Has an uppercase character
	</p>
	<p id="lowercase" class="password-tip" class:met={lowercaseMet} class:not-met={!lowercaseMet}>
		Has a lowercase character character
	</p>
	<p id="number" class="password-tip" class:met={numberMet} class:not-met={!numberMet}>
		Has a number
	</p>
	<p id="symbol" class="password-tip" class:met={symbolMet} class:not-met={!symbolMet}>
		Has a symbol (eg. &,$,!)
	</p>
</article>

<style>
	#password-tips {
		margin-top: 10px;
	}

	.password-tip {
		margin-top: 3px;
	}

	.not-met {
		color: red;
	}

	.met {
		color: green;
	}

	label {
		align-self: baseline;
		margin-top: 20px;
	}

	input[type='password'] {
		padding: 5px 7px;
		width: 100%;
	}
</style>
