<script>
  import { onMount } from 'svelte'

  let bucket = ''
  let citation = ''

  async function addBucketAndCitation() {
    const response = await fetch('/api/v1/collection', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ bucket, citation }),
    })
    if (!response.ok) {
      console.error('Failed to add bucket and citation')
    }
  }

</script>

<h1>Add new bucket and citation</h1>
<form on:submit|preventDefault={addBucketAndCitation}>
  <div>
    <label for="bucket">Bucket:</label>
    <input id="bucket" bind:value={bucket} required />
  </div>
  <div>
    <label for="citation">Citation:</label>
    <input id="citation" bind:value={citation} required />
  </div>
  <button type="submit">Add</button>
</form>
