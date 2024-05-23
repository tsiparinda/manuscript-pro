<script>
  import { afterUpdate } from 'svelte'
  import BlockEditor from '../components/BlockEditor_self.svelte'

  export let id_col

  let err, collection, user, userrights, headerHeight, header, pane
  let nonEmptyCollection = false

  $: Promise.all([
    getCollection(id_col),
    getUser(),
    getCollectionUserRights(id_col),
  ])
    .then(([c, u, ur]) => {
      collection = c
      user = u
      userrights = ur
    })
    .catch((e) => (err = e))

  async function getCollectionUserRights(id_col) {
    const res = await fetch(`/api/v1/collectionuserrights/${id_col}`)
    if (res.ok) {
      const d = await res.json()
      return d.data
    } else {
      throw new Error(res.body)
    }
  }
  async function getCollection(id_col) {
    const res = await fetch(`/api/v1/collection/${id_col}`)
    if (res.ok) {
      const d = await res.json()
      return d.data
    } else {
      throw new Error(res.body)
    }
  }

  async function getUser() {
    const res = await fetch(`/api/v1/user`)
    if (res.ok) {
      const d = await res.json()
      return d.data
    } else {
      throw new Error(res.body)
    }
  }

  $: if (collection && collection.buckets) {
    collection.buckets.sort((a, b) =>
      a.bucket.toLowerCase().localeCompare(b.bucket.toLowerCase())
    )
    collection.buckets.forEach((bucket) => {
      bucket.citations.sort((a, b) => {
        const aParts = a
          .replace(/[-+/_|].*/, '')
          .split('.')
          .map(Number)
        const bParts = b
          .replace(/[-+/_|].*/, '')
          .split('.')
          .map(Number)

        for (let i = 0; i < Math.min(aParts.length, bParts.length); i++) {
          if (aParts[i] < bParts[i]) {
            return -1
          } else if (aParts[i] > bParts[i]) {
            return 1
          }
        }

        // If all parts compared are equal, sort by length
        return aParts.length - bParts.length
      })
    })
  }

  $: if (collection && Array.isArray(collection.buckets)) {
    nonEmptyCollection = true
  }

  afterUpdate(async () => {
    // Get the height of the header element
    header = document.querySelector('.header')
    pane = document.querySelector('.pane')
  })

  $: headerHeight = header ? header.clientHeight : 0
  // Set the padding-top for the pane element
  $: if (pane) {
    pane.style.paddingTop = `${headerHeight + 10}px`
  }

  let visibleEditorIndex = -1

  function toggleEditor(index) {
    if (visibleEditorIndex === index) {
      visibleEditorIndex = -1 // Hide the BlockEditor if it's already visible
    } else {
      visibleEditorIndex = index // Show the BlockEditor
    }
  }

  let showAllEditors = false

  function toggleAllEditors() {
    showAllEditors = !showAllEditors
  }

</script>

<style>
  .pane {
    flex-basis: 50%;
    flex-shrink: 0;
    flex-grow: 1;
    position: relative;
    z-index: 10;
    display: block;
    padding-top: 0; /* Initialize padding to 0 */
  }

  .header {
    display: flex;
    position: fixed;
    background: #ffffff;
    width: 100%;
    z-index: 999;
  }

  .toggle-button {
    font-size: 1.2em; /* Adjust to the size you want */
    padding: 5px 10px; /* Adjust to the padding you want */
    /* Add more styles if needed, like colors, border, etc */
  }

</style>

{#if err}
  <h5>Error: You do not have permission to view this collection</h5>
  <h6>
    You can return to the
    <a href="/">Brucheion home page</a>
    to see all available collections
    <span>{err}</span>
  </h6>
{:else}
  {#if !!collection}
    <div class="header">
      <h2 style="margin-left: 10px">{collection.title}</h2>
      {#if nonEmptyCollection}
        <button on:click={toggleAllEditors}>Toggle All Editors</button>
      {/if}
    </div>
    <div class="pane">
      {#if nonEmptyCollection}
        {#each collection.buckets as { bucket, citations }}
          <div class="pane1">
            <a
              href={`/view/${id_col}/${bucket}`}
              target="_blank"
              rel="noreferrer">
              {bucket}
            </a>
            {#each citations as citation, i (citation)}
              <div class="pane1">
                <button
                  class="toggle-button"
                  on:click={() => toggleEditor(i)}>&gt;</button>
                <a
                  href={`/view/${id_col}/${bucket + citation}`}
                  target="_blank"
                  rel="noreferrer">
                  {bucket.split(':')[bucket.split(':').length - 2]}:{citation}
                </a>
              </div>
              {#if showAllEditors || visibleEditorIndex === i}
                <div class="pane1">
                  <BlockEditor {id_col} {bucket} {citation} />
                </div>
              {/if}
            {/each}
          </div>
        {/each}
      {:else}
        <p>No witness found.</p>
      {/if}
    </div>
  {/if}
{/if}
