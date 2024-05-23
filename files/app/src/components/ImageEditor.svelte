<!-- ImageReference.svelte -->
<script>
  import { onMount } from 'svelte'
  import { createEventDispatcher } from 'svelte'
  import { writable } from 'svelte/store'
  import Modal from './Modal.svelte'

  const dispatch = createEventDispatcher()

  export let passage, colID, savedata

  let ctsurn // Replace with the value from your original code
  // let imageRef = '' // Replace with the value from your original code
  // let imageJS = '' // Replace with the value from your original code
  let selectedCollection = '',
    selectedImage = writable(null),
    collections = [],
    s3collections = [],
    imagenames = [],
    locations = [],
    colnames = [],
    imageURL,
    imageRefs = [],
    activeTab = 'local'

  let showModal = false,
    modalHeader,
    modalText

  let isLoading = false

  function addImageToCollection() {
    if (selectedImage && !imageRefs.includes(selectedImage)) {
      imageRefs = [...imageRefs, selectedImage]
    }
  }

  function removeImageFromCollection(imageName) {
    imageRefs = imageRefs.filter((image) => image !== imageName)
  }

  onMount(async () => {
    imageRefs = [...passage.imageRefs]
    ctsurn = passage.passageid
    populateCollectionDropdown()
    populateS3CollectionDropdown()
  })

  $: handleSubmit(savedata)

  $: if (!!collections) {
    collections.forEach((entry) => {
      const option = document.createElement('option')
      option.value = entry.colname
      option.text = entry.colname
    })

    imagenames = collections.map((entry) => entry.imagename)
    locations = collections.map((entry) => entry.location)

    colnames = collections
      .map((entry) => entry.colname)
      .filter((value, index, self) => self.indexOf(value) === index)
  }

  $: filteredImagenames = collections
    .filter((entry) => entry.colname === selectedCollection)
    .map((entry) => entry.imagename)
    .filter((imagename) => !imageRefs.includes(imagename))

  $: filteredLocalImagenames = s3collections.filter(
    (entry) => !imageRefs.includes(entry)
  )
  //  .map((entry) => entry)
  //  .filter(( ) => !imageRefs.includes( entry))

  async function populateCollectionDropdown() {
    const collectionurl = `/api/v1/collectionimages/${colID}`
    const res = await fetch(collectionurl)
    if (res.ok) {
      collections = (await res.json()).data
    } else {
      throw new Error(res.body)
    }
  }

  async function populateS3CollectionDropdown() {
    const collectionurl = `/api/v1/localimages`
    const res = await fetch(collectionurl)
    if (res.ok) {
      s3collections = (await res.json()).data
    } else {
      throw new Error(res.body)
    }
  }

  function handleImageChange(event) {
    const selectedIndex = imagenames.indexOf(selectedImage)

    if (selectedIndex !== -1) {
      imageURL = locations[selectedIndex]
      // console.log('Selected imageURL:', imageURL)

      // Emit the custom event with the new imageURL
      dispatch('imageurlchange', { imageURL })
    } else {
      console.log('Image not found in the imagenames array')
    }
  }

  async function handleSubmit() {
    if (!savedata) return

    isLoading = true // Set isLoading to true before making the request

    const id = colID
    // const imageref = imageRefs.replaceAll('#', '+')
    const formData = {
      colid: parseInt(id),
      ctsurn: ctsurn,
      imageref: imageRefs,
    }

    const response = await fetch('/api/v2/savereference', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(formData),
    })

    isLoading = false // Set isLoading back to false after the request is complete

    if (response.ok) {
      // Handle successful POST action
      modalHeader = 'Success'
      modalText = 'Image references have been saved'
      showModal = true
    } else {
      // Handle error
      modalHeader = 'Error'
      modalText = 'Image references have not been saved'
      showModal = true
    }
    savedata = false
    window.location.href = `/view/${colID}/${passage.passageid}`
  }

  function selectImage(image) {
    //selectedImage.set(image); // Set the value of the selectedImage store to the clicked image
    // console.log('Image ', image, imageRefs)
    const selectedIndex = imageRefs.indexOf(image)

    if (selectedIndex !== -1) {
      imageURL = imageRefs[selectedIndex]
      // console.log('Selected imageURL:', imageURL)
    } else {
      console.log('Image not found in the imagenames array')
    }
    dispatch('imageurlchange', { imageURL })
  }

</script>

<style>
  .tab-content {
    min-height: 150px; /* Adjust the height as needed */
    overflow-y: auto; /* Enable vertical scrolling if content exceeds the height */
  }
  .tabs ul {
    display: flex;
    flex-direction: row;
    list-style: none;
    padding: 0;
    margin: 0;
  }

  /* Tab navigation styles */
  .tabs {
    border-bottom: 1px solid #ccc; /* Add a line to separate the navigation from the content */
    font-family: 'Arial', sans-serif; /* Change the font of the tab names */
    margin-bottom: 1rem;
  }

  .tabs li {
    margin-right: 1rem;
    cursor: pointer;
    padding: 0.5rem 1rem;
    border: 1px solid transparent;
    border-radius: 4px 4px 0 0;
  }

  /* Active tab styles */
  .tabs li.is-active {
    font-weight: bold; /* Make the font of the active tab name bold */
    border-color: #ccc;
    border-bottom-color: transparent;
    background-color: #f9f9f9;
  }

  li {
    margin-bottom: 5px; /* Adjust this value as needed */
  }
  .image-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .delete-button {
    margin-left: 20px; /* Adjust this value as needed */
  }

  .select-container {
    width: 100%;
  }
  select {
    width: 100%;
  }

</style>

<div
  class="container is-fluid"
  style="cursor: {isLoading ? 'progress' : 'default'}">
  <section>
    <div class="tile is-ancestor">
      <div class="tile is-parent">
        <div class="tile is-vertical">
          <div class="tile is-child">
            <!-- Add tab navigation -->
            <div class="tabs">
              <ul>
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <li
                  class:is-active={activeTab === 'local'}
                  on:click={() => (activeTab = 'local')}>
                  <!-- svelte-ignore a11y-missing-attribute -->
                  <a>External Images</a>
                </li>
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <li
                  class:is-active={activeTab === 's3'}
                  on:click={() => (activeTab = 's3')}>
                  <!-- svelte-ignore a11y-missing-attribute -->
                  <a>Internal Images</a>
                </li>
              </ul>
            </div>

            <div class="tab-content">
              <!-- Local Images Tab -->
              {#if activeTab === 'local'}
                <div class="form-group">
                  <label for="sel1">Select Collection:</label>
                  <select
                    bind:value={selectedCollection}
                    class="form-control"
                    id="image_colSelect">
                    <option disabled>Choose collection</option>
                    {#if !!collections}
                      {#each colnames as col}
                        <option value={col}>{col}</option>
                      {/each}
                    {:else}
                      <option disabled>No collections</option>
                    {/if}
                  </select>
                </div>

                <div class="form-group">
                  <label for="sel1">Select Image:</label>
                  <select
                    bind:value={selectedImage}
                    class="form-control"
                    id="image_urnSelect"
                    on:change={handleImageChange}>
                    <option disabled>Choose folio</option>
                    {#each filteredImagenames as imagename}
                      <option value={imagename}>{imagename}</option>
                    {/each}
                  </select>
                </div>
              {/if}

              <!-- S3 Bucket Images Tab -->
              {#if activeTab === 's3'}
                <div class="form-group">
                  <select bind:value={selectedImage} class="form-control">
                    <option value="">-- select an image --</option>
                    {#each filteredLocalImagenames as image}
                      <option>{image}</option>
                    {/each}
                  </select>
                </div>
              {/if}
              <button
                class="button is-small is-primary"
                on:click={addImageToCollection}>
                Add
              </button>
            </div>
            <ul>
              <!-- Display the list of collection images -->
              {#each imageRefs as image}
                <li>
                  <!-- svelte-ignore a11y-click-events-have-key-events -->
                  <div
                    class="image-item"
                    on:click={(event) => {
                      event.stopPropagation()
                      selectImage(image)
                    }}>
                    {image}
                    <button
                      class="btn btn-primary  delete-button"
                      on:click={(event) => {
                        event.stopPropagation() // Prevent the li click event from triggering
                        removeImageFromCollection(image)
                      }}>
                      DEL
                    </button>
                  </div>
                </li>
              {/each}
            </ul>
          </div>
        </div>
      </div>
    </div>
  </section>
</div>

{#if showModal}
  <Modal bind:showModal>
    <h3 slot="header">{modalHeader}</h3>

    <body>
      <h4>{modalText}</h4>
    </body>
  </Modal>
{/if}
