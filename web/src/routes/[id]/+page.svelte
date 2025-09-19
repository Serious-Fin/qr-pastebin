<script lang="ts">
    import type { PageProps } from './$types'
    let { data }: PageProps = $props()
    const share = data.share
</script>

<section id="main">
	<h1>Quick Share</h1>
    {#if (share.title !== undefined && share.title !== "")}
        <h2>{share.title}</h2>
    {/if}
    <textarea id="content" name="content" readonly>{share.content}</textarea>
    {#if (share.expireAt !== undefined)}
        <h2>{share.expireAt}</h2>
    {/if}
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

    #content {
		width: 80vw;
		max-width: 90vw;
		min-height: 100px;
		margin: 15px 0;
		padding: 7px 8px;
		border-radius: 10px;
	}
</style>