<section id="main">
	<h1>Quick Share</h1>
	<h2>New share</h2>
	<form action="POST" method="POST">
		<textarea id="content" name="content" required></textarea>
		<label>
			Title:
			<input type="text" class="property-input" name="title" />
		</label>

		<label>
			Expiration:
			<input type="date" class="property-input" name="expiration" />
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
