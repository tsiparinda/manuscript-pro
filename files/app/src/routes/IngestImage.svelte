<script>
  // import { stringify as stringifyQuery } from 'qs'
  import OpenSeadragon from '../lib/openseadragon/openseadragon'
  import { onMount } from 'svelte'
  import { navigate } from 'svelte-routing'
  import FormLine from '../components/FormLine.svelte'
  import Message from '../components/Message.svelte'
  import { validateUrn, validateNoColon } from '../lib/cts-urn'
  import TextInput from '../components/TextInput.svelte'
  import { validateHttpUrl } from '../lib/url'
  import { isIIIFImage } from '../lib/iiif'
  import { getStaticOpts, getIIIFOpts, getInternalOpts } from '../lib/osd'
  import Modal from '../components/Modal.svelte'
  import {
    listDZIFiles,
    s3_bucket_url_js,
    listDZIcollection,
  } from '../lib/s3helper.js'

  export let id_col

  let collection = ''
  let imageName = ''
  let imageUrl = ''
  let external = true
  let protocol = 'static'

  let statusMessage = null,
    timeoutHandle = null
  let collectionRef, imageNameRef
  let collections = []
  let nameExists = false
  let previewViewer = undefined,
    viewerOpts = undefined,
    previewVisible = false,
    previewErrored = false

  let colinfo,
    colimages = [],
    err = 'Loading information about the collection...'

  let showModal = false,
    modalHeader,
    modalText

  let namespace = '',
    namespaces = [],
    collectionID = '',
    collectionIDs = [],
    imageID = ''

  $: imageName = 'urn:cite2:' + namespace + ':' + collectionID + ':' + imageID
  $: collection = 'urn:cite2:' + namespace + ':' + collectionID + ':'
  $: validNames =
    validateUrn(collection, { noPassage: true }) && validateUrn(imageName)
  $: validSource = validateUrn(imageUrl) || validateHttpUrl(imageUrl)
  $: complete = validNames && validSource

  $: if (statusMessage !== null) {
    clearTimeout(timeoutHandle)
    timeoutHandle = setTimeout(() => (statusMessage = null), 10000)
  }
  $: errorMessage =
    statusMessage && statusMessage.toLowerCase().includes('error')
  $: external = !validateUrn(imageUrl)
  $: if (validNames) {
    // fetch(`/api/v2/imageinfo/${id_col}/${collection}/${imageName}`).then(
    //   async (res) => {
    //     const imageInfo = await res.json()
    //     nameExists = !!imageInfo.data.imagename
    //   }
    // )
    nameExists = !!collections.imagename
  } else if (nameExists) {
    nameExists = false
  }

  $: if (validSource) {
    previewErrored = false

    if (validateHttpUrl(imageUrl)) {
      displayExternalMedia(imageUrl)
    } else if (validateUrn(imageUrl)) {
      viewerOpts = getInternalOpts('preview', imageUrl)
      protocol = 'localDZ'
    }
  }

  $: if (!!collections) {
    colimages = collections.map((entry) => entry.colname)

    namespaces = collections.map((entry) => {
      let urnParts = entry.colname.split(':')
      if (urnParts.length >= 3) {
        return urnParts[2]
      } else {
        return null // or any other default value, or you can exclude it
      }
    })

    collectionIDs = collections.map((entry) => {
      let urnParts = entry.colname.split(':')
      if (urnParts.length >= 4) {
        return urnParts[3]
      } else {
        return null // or any other default value, or you can exclude it
      }
    })
  }

  onMount(async () => {
    const query = new URLSearchParams(location.search)
    if (query.has('collection')) {
      if (validateUrn(query.get('collection'), { noPassage: true })) {
        collection = query.get('collection')
        imageNameRef.focus()
      } else {
        query.delete('collection')
        navigate(`/ingest?${query.toString()}`, { replace: true })
      }
    } else {
      // collectionRef.focus()
    }
    await fetchCollectionImages()
    await fetchCollectionInfo()
    //await getDZIcollection()
  })

  async function displayExternalMedia(imageUrl) {
    try {
      const [isManifest, imageManifest] = await isIIIFImage(imageUrl)
      if (isManifest) {
        viewerOpts = getIIIFOpts('preview', imageManifest)
        protocol = 'iiif'
      } else {
        viewerOpts = getStaticOpts('preview', imageUrl)
        protocol = 'static'
      }
    } catch (err) {
      if (!err.message.includes('NetworkError')) {
        console.error(err.message)
      }

      viewerOpts = getStaticOpts('preview', imageUrl)
      protocol = 'static'
    }
  }

  async function getDZIcollection() {
    const bucketName = 'brucheion'
    const prefix = 'nbh/J1img/positive/'
    // const prefix = '';
    const dziCollection = await listDZIcollection(bucketName, prefix)
    console.log('List of DZI files with folder paths and names:', dziCollection)
  }

  function createViewer(opts) {
    const { tileSources, ...otherOpts } = opts
    previewViewer = OpenSeadragon(otherOpts)

    previewViewer.addHandler('open-failed', () => {
      previewVisible = false
      previewViewer.destroy()
      previewErrored = true
    })

    previewViewer.addHandler('open', () => {
      previewVisible = true
    })

    previewViewer.open(tileSources)
  }

  /* We'll need to trick the Svelte reactivity here, since destroying a prior viewer before creating a new one will result
   * in a circular dependency within the $-statement. Hence, above we just create the viewer options and handle viewer
   * lifecycles in the below $-statement.
   */
  $: if (validSource && viewerOpts) {
    if (previewViewer) {
      previewVisible = false
      previewViewer.destroy()
    }

    createViewer(viewerOpts)
  }

  // ready
  async function fetchCollectionImages() {
    const res = await fetch(`/api/v1/collectionimages/${id_col}`)
    if (res.ok) {
      collections = (await res.json()).data
    } else {
      throw new Error(res.body)
    }
  }

  async function fetchCollectionInfo() {
    const res = await fetch(`/api/v1/collection/${id_col}`)
    if (res.ok) {
      colinfo = (await res.json()).data
      //  console.log('colinfo', colinfo)
    } else {
      throw new Error(res.body)
    }
  }

  async function handleSubmit(event) {
    event.preventDefault()
    if (!complete) {
      return
    }

    let query = {
      colid: parseInt(id_col),
      imagename: imageName,
      colname: collection,
      protocol: protocol,
      license: 'CC-BY-4.0',
      external: external,
      location: imageUrl,
    }

    const res = await fetch('/api/v2/addimagetocite', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(query),
    })

    if (res.ok) {
      modalHeader = 'Success'
      modalText = 'The image has been saved to the image collection'
      showModal = true
      // showMetadata = false
    } else {
      // Handle error
      console.error(`Ingestion failed: HTTP ${res.status} ${await res.text()}`)
      modalHeader = 'Error'
      modalText = 'An error occurred. Please try later.'
      showModal = true
      return
    }
  }

</script>

<style>
  .form-column {
    max-width: 724px;
  }

  .form {
    box-sizing: border-box;
    max-width: 700px;
    padding: 25px;
  }

  select {
    background-color: white;
    border-color: #dbdbdb;
    border-radius: 4px;
    color: #363636;
  }

  .preview-container {
    margin-top: 25px;
    padding: 25px;
    opacity: 0;

    transition: opacity 125ms ease-out;
  }

  .preview-container.visible {
    opacity: 1;
  }

  @media screen and (min-width: 1088px) {
    .preview-container {
      margin-top: 0;
    }
  }

  .preview {
    box-sizing: border-box;
    max-width: 701px;
    height: 601px;
    border: 2px solid rgba(230, 230, 230);
    border-radius: 4px;
    padding: 4px;

    background: rgba(246, 245, 245);
    box-shadow: 0px 0px 5px rgba(0, 0, 0, 0.15);
  }

</style>

<div class="container is-fluid">
  <section>
    <div class="columns is-desktop">
      <div class="column form-column">
        <form class="form" on:submit={handleSubmit}>
          {#if colinfo}
            <h3 style="width: 90%;">
              Define Identifier for "{colinfo.title}" collection
            </h3>
          {:else if err}
            <h5>{err}</h5>
          {/if}

          <FormLine
            id="URN_prefix"
            label="Image name in URN format"
            info="Like `urn:cite2:nbh:J1img.positive:J1_31r`. Fill all bottom filelds to take result!">
            <TextInput id="URN_prefix" bind:value={imageName} disabled="true" />
          </FormLine>

          <FormLine
            id="namespace"
            label="Namespace"
            info="Name of your image collection">
            <TextInput
              id="namespace"
              placeholder="Select or enter namespace"
              bind:value={namespace}
              validate={(value) => validateNoColon(value)}
              items={namespaces}
              invalidMessage="Please enter a valid namespace without COLON(:)." />
          </FormLine>

          <FormLine
            id="collection"
            label="Collection ID and/or Object Type"
            info="Like `nyaya.positive`">
            <TextInput
              id="collection"
              placeholder="Select or enter collection ID and/or object type"
              bind:value={collectionID}
              validate={(value) => validateNoColon(value)}
              items={collectionIDs}
              invalidMessage="Please enter a value without COLON(:)." />
          </FormLine>

          <FormLine
            id="image_id"
            label="Image Identifier"
            info="Image number in collection. Like `J1D_35`">
            <TextInput
              id="image_id"
              placeholder="Enter image identifier"
              bind:value={imageID}
              validate={(value) => validateNoColon(value)}
              invalidMessage="Please enter a valid image identifier without COLON(:)." />
          </FormLine>

          <FormLine
            id="source"
            label="Source"
            info="Link to external source like `https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png`">
            <TextInput
              id="source"
              placeholder="Resource URL"
              bind:value={imageUrl}
              validate={(value) => validateUrn(value) || validateHttpUrl(value)}
              invalidMessage="Please enter a valid CITE object URN or a HTTP(S)
              URL." />
            {#if previewErrored}
              <Message
                text="The media could not be loaded for preview due to errors.
                You can ingest it nonetheless." />
            {/if}
          </FormLine>

          <FormLine
            id="protocol"
            label="Type"
            info="Usually defined automatically when Source defined">
            <div class="select">
              <select id="protocol" bind:value={protocol}>
                <option value="static">Static</option>
                <!-- <option value="localDZ">Deep Zoom (local)</option> -->
                <option value="iiif">IIIF</option>
              </select>
            </div>
          </FormLine>

          <FormLine offset>
            <button
              type="button"
              class="btn btn-primary"
              disabled={!complete}
              on:click={handleSubmit}>
              Add Image
            </button>
            {#if statusMessage}
              <Message text={statusMessage} error={errorMessage} />
            {/if}
          </FormLine>
        </form>
      </div>
      <div class="column form-column">
        <div class="preview-container" class:visible={previewVisible}>
          <h3 class="title is-4">Preview</h3>
          <div id="preview" class="preview" />
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
