<section id="main">
	<h1>Quick Share</h1>
	<h2>New share</h2>
	<form method="POST" action="?/createShare">
		<textarea id="content" name="content" required></textarea>
		<label>
			Title:
			<input type="text" class="property-input" name="title" />
		</label>

		<label>
			Expire in:
			<select class="property-input" name="expireIn">
				<option value="" selected>Never</option>
				<option value="10_minutes">10 minutes</option>
				<option value="1_hours">1 hour</option>
				<option value="1_days">1 day</option>
				<option value="1_weeks">1 week</option>
				<option value="2_weeks">2 weeks</option>
				<option value="1_months">1 month</option>
				<option value="6_months">6 months</option>
				<option value="1_years">1 year</option>
			</select>
		</label>

		<label>
			Password:
			<input type="password" class="property-input" name="password" />
		</label>

		<input type="submit" value="Share it" />
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
		font-size: 16pt;
		font-weight: 500;
	}

	form {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	#content {
		width: 80vw;
		max-width: 90vw;
		min-height: 100px;
		margin: 15px 0;
		padding: 7px 8px;
		border-radius: 10px;
	}
</style>
