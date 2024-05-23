<script>
  import { onMount } from 'svelte'
  import { Link, navigate } from 'svelte-routing' // export let urn
  // export let userrights
  import OpenSeadragon from '../lib/openseadragon/openseadragon'
  import ResizeBarLandscape from '../components/ResizeBarLanscape.svelte'
  import ResizeBarPortrait from '../components/ResizeBarPortrait.svelte'
  import LineEditor from '../components/LineEditor.svelte'
  import BlockEditor from '../components/BlockEditor.svelte'
  import MetadataEditor from '../components/MetadataEditor.svelte'
  import ImageEditor from '../components/ImageEditor.svelte'
  import { validateUrn } from '../lib/cts-urn'
  import { validateHttpUrl } from '../lib/url'
  import { isIIIFImage } from '../lib/iiif'
  import { getStaticOpts, getIIIFOpts, getInternalOpts } from '../lib/osd'

  export let id_col, urn
  let transcriptionFontSize = '1em'
  let tFontSize = 1

  let previewViewer = undefined,
    viewerOpts = undefined,
    previewVisible = false,
    previewFailed = false,
    selectedImageRef = undefined,
    didMount = false,
    selectedCatalogUrn, // = passage.catalog.urn,
    showMetadata = false,
    showImageEdit = false,
    showReadingLine = false,
    horizontalLine,
    previewContainer = undefined,
    previewHeight = 350,
    previewWidth = 500,
    err,
    isPortrait = false,
    newpassage,
    lineIndex,
    isLineEditor = false,
    isBlockEditor = false,
    transcriptionEditNav,
    passagetext,
    passage,
    user,
    userrights,
    passageId,
    currentTranscriptionLines = [],
    prevTranscriptionLines,
    newpassagetext,
    loading = true,
    savedata = false,
    protocol = 'static',
    images = [],
    schemes = [],
    imageURL,
    retryCount = 3,
    isFullyLoaded = false,
    isErrorLoading = false,
    hidePassageToolbar = false,
    hideTranscriptionToolbar = false,
    autoChangeOrientation = true,
    newUrn,
    sectionOpenSeaDragon,
    sectionTranscription,
    sectionOpenSeaDragonSize,
    sectionTranscriptionSize

  $: hidePassageToolbar =
    isBlockEditor || isLineEditor || showMetadata || showImageEdit

  $: hideTranscriptionToolbar = showMetadata || showImageEdit
  //console.info('from components colid, passage: ', id_col, passage)
  // $: if (!validateUrn(urn, { nid: 'cts' })) {
  //   //  urn = 'urn:cts:sktlit:skt0001.nyaya002.M3D:5.1.1'
  //   console.log('!!!', urn, validateUrn(urn, { nid: 'cts' }))
  //   urn = 'undefined'
  //   //  err = new Error('Passage not found')
  // }
  $: if (newUrn) {
    navigate(`/view/${id_col}/${newUrn}`)
  }

  onMount(() => {
    sectionOpenSeaDragon = document.getElementById('openseadragon_sec')
    sectionTranscription = document.getElementById('transcription_sec')

    //await loadData();
    didMount = true
    previewWidth = window.innerWidth / 2 - 6

    function updateFolioid() {
      const params = new URLSearchParams(window.location.search)
      const folioid = params.get('folioid')
      if (folioid) {
        // If a folioid exists in the query string, use it as the selected image ref
        selectedImageRef = folioid
      } else if (passage && passage.imageRefs && passage.imageRefs.length > 0) {
        // Otherwise, if there are image references, use the first one
        selectedImageRef = passage.imageRefs[0]
      } else {
        // If no image references exist, set selectedImageRef to null or some default value
        selectedImageRef = null
      }
    }

    // Update folioid when the component mounts
    updateFolioid()

    // And whenever the URL changes
    window.addEventListener('popstate', updateFolioid)

    return () => {
      // Clean up the event listener when the component is unmounted
      window.removeEventListener('popstate', updateFolioid)
    }
  })

  $: Promise.all([
    getPassage(id_col, urn),
    getUser(),
    getCollectionUserRights(id_col),
    getImages(id_col),
  ])
    .then(([p, u, ur, ci]) => {
      passage = p
      user = u
      userrights = ur
      passageId = p.passageid.split(':').pop()
      selectedCatalogUrn = p.catalog.urn
      if (p.transcriptionLines) {
        currentTranscriptionLines = [...p.transcriptionLines]
      }
      if (p.schemes) {
        schemes = [...p.schemes]
      }
      images = [...ci]
      loading = false
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
  async function getPassage(id_col, urn) {
    loading = true
    const res = await fetch(`/api/v1/passage/${id_col}/${urn}`)
    if (res.ok) {
      const d = await res.json()
      return d.data
    } else {
      throw new Error(res.body)
    }
  }

  async function getImages(id_col) {
    loading = true
    const res = await fetch(`/api/v1/collectionimages/${id_col}`)
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

  async function displayExternalMedia(imageUrl) {
    // console.log('displayExternalMedia', imageUrl)
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

  // FIXME: this is a pretty naive attempt to catch the passage ID
  $: if (!!passage) {
    passageId = passage.passageid.split(':').pop()
  }

  /* this should update the folio viewer a) once after mounting and b) when `passage` changes due to reactivity.
   * this is just a lazy trick to trigger the viewer update in coordination with svelte's reactivity */
  $: if (!!passage && didMount) {
    updateViewer(passage.imageRefs)
  }

  $: validSource = validateUrn(selectedImageRef) || validateHttpUrl(imageURL)

  $: if (validSource) {
    previewFailed = false
    // console.log('imageURL', imageURL)

    if (validateHttpUrl(imageURL)) {
      // console.log('displayExternalMedia', imageURL, validateHttpUrl(imageURL))
      displayExternalMedia(imageURL)
    } else if (validateUrn(selectedImageRef)) {
      //   console.log('validateUrn(selectedImageRef)', validateUrn(selectedImageRef), selectedImageRef)
      viewerOpts = getInternalOpts('preview', selectedImageRef)
      protocol = 'localDZ'
      viewerOpts.defaultZoomLevel = 1
    }
  }

  $: if (validSource && viewerOpts) {
    if (previewViewer) {
      previewVisible = false
      previewViewer.destroy()
    }
    createViewer(viewerOpts)
  }

  function handleImageUrlChange(event) {
    //console.log('handleImageUrlChange imageURLiii:')
    selectedImageRef = event.detail.imageURL
    getimageURLfromCollection()
    //console.log('handleImageUrlChange imageURL:', selectedImageRef)
    // Update the OpenSeadragon viewer with the new imageURL
  }

  function createViewer(opts) {
    opts.defaultZoomLevel = 1
    opts.preserveImageSizeOnResize = true
    opts.visibilityRatio = 1
    opts.constrainDuringPan = true

    const { tileSources, ...otherOpts } = opts
    document.getElementById('osd_buttons').innerHTML = ''
    previewViewer = OpenSeadragon(otherOpts)
    previewViewer.open(tileSources)
    createLine()
    previewViewer.addHandler('open-failed', (e) => {
      if (retryCount > 0) {
        retryCount--
        previewVisible = false
        previewViewer.destroy()
        previewFailed = true
        viewerOpts = getInternalOpts('preview', selectedImageRef)
        createViewer(opts)
        //updateViewer(passage.imageRefs)
      } else {
        isErrorLoading = true
        isFullyLoaded = false
        updateLoadingIndicator()
        previewVisible = false
        previewViewer.destroy()
        previewFailed = true
      }
    })

    previewViewer.addHandler('open', (e) => {
      isFullyLoaded = false
      isErrorLoading = false
      updateLoadingIndicator()
      var x = previewViewer.source.dimensions.x
      var y = previewViewer.source.dimensions.y
      if (y >= x && autoChangeOrientation) {
        isPortrait = true
      } else {
        isPortrait = false
      }
      previewVisible = true
    })

    previewViewer.world.addHandler('add-item', function (event) {
      var tiledImage = event.item
      tiledImage.addHandler('fully-loaded-change', function () {
        var newFullyLoaded = areAllFullyLoaded()
        if (newFullyLoaded !== isFullyLoaded) {
          isFullyLoaded = newFullyLoaded
          updateLoadingIndicator()
        }
      })
    })

    previewViewer.open(tileSources)
    handleToggleLine()
    handleToggleLine()
  }

  function areAllFullyLoaded() {
    var tiledImage
    var count = previewViewer.world.getItemCount()
    for (var i = 0; i < count; i++) {
      tiledImage = previewViewer.world.getItemAt(i)
      if (!tiledImage.getFullyLoaded()) {
        return false
      }
    }
    return true
  }

  function updateLoadingIndicator() {
    if (!document.querySelector('.loading')) {
      return
    }
    if (!isFullyLoaded && !isErrorLoading) {
      document.querySelector('.loading').style.display = 'flex'
      document.querySelector('.loading').textContent = 'Loading image...'
    } else if (isFullyLoaded || !isErrorLoading) {
      document.querySelector('.loading').style.display = 'none'
    } else if (isErrorLoading) {
      document.querySelector('.loading').style.display = 'flex'
      document.querySelector('.loading').textContent = 'Error loading image...'
    } else if (!isFullyLoaded) {
      document.querySelector('.loading').style.display = 'flex'
      document.querySelector('.loading').textContent = 'Loading image...'
    }
  }

  function getimageURLfromCollection() {
    // if selectedImageRef is in images.imagename array
    // get this images element from images array
    imageURL = undefined
    for (const element of images) {
      // console.log('Element:', element, selectedImageRef);
      if (element.imagename === selectedImageRef) {
        imageURL = element.location
        break
      }
      imageURL = selectedImageRef
    }
  }

  function updateViewer(refs) {
    if (previewViewer) {
      previewViewer.destroy()
    }
    //  console.log('viewerOpts', viewerOpts)
    if (Array.isArray(refs) && refs.length > 0) {
      if (!selectedImageRef || !refs.includes(selectedImageRef)) {
        selectedImageRef = refs[0]
      }
      // console.log('before getImageFromCollection', viewerOpts)
      getimageURLfromCollection()
    }
  }

  function handleSelect(event) {
    newUrn = selectedCatalogUrn + event.target.value
  }

  function handleWitnessSelection() {
    newUrn = selectedCatalogUrn
  }

  function handleFirstPassage() {
    newUrn = passage.firstPassage
  }
  function handlePreviousPassage() {
    newUrn = passage.previousPassage
  }
  function handleNextPassage() {
    if (passage.nextPassage !== '') {
      newUrn = passage.nextPassage
    }
  }
  function handleLastPassage() {
    if (passage.lastPassage !== '') {
      newUrn = passage.lastPassage
    }
  }

  function handleFolioSelection() {
    navigate(`/view/${id_col}/${newUrn}?folioid=${selectedImageRef}`)
  }

  function handleToggleMetadata() {
    showMetadata = !showMetadata
  }
  function handleHideMetadata() {
    showMetadata = false
  }

  function handleToggleLine() {
    showReadingLine = !showReadingLine
    if (showReadingLine) {
      document.getElementById('lineOpts').style.display = 'flex'

      horizontalLine.style.display = 'block'
      updateLinePosition()
      return
    }

    document.getElementById('lineOpts').style.display = 'none'
    horizontalLine.style.display = 'none'
  }

  function createLine() {
    horizontalLine = document.createElement('div')
    horizontalLine.id = 'hline'
    horizontalLine.style.position = 'absolute'
    horizontalLine.style.width = '100%'
    horizontalLine.style.height = '3px'
    horizontalLine.style.backgroundColor = 'red'
    horizontalLine.style.pointerEvents = 'none'
    horizontalLine.style.display = 'none'
    previewViewer.canvas.appendChild(horizontalLine)
  }

  function updateLinePosition() {
    const viewportOrigin = previewViewer.viewport.pixelFromPoint(
      previewViewer.viewport.getHomeBounds().getTopLeft()
    )
    const viewportY = previewViewer.viewport.getCenter(true).y
    var lineY = previewViewer.viewport.viewportToViewerElementCoordinates(
      new OpenSeadragon.Point(0, viewportY)
    ).y

    lineY +=
      (parseFloat(document.getElementById('lineOffsetRange').value) * lineY) /
      100

    //horizontalLine.style.top = `${lineY - viewportOrigin.y}px`;
    horizontalLine.style.top = `${lineY}px`

    // add colorPicker
    horizontalLine.style.backgroundColor =
      document.getElementById('myColorPicker').value
    localStorage.setItem(
      'linecolor',
      document.getElementById('myColorPicker').value
    )
  }

  function handleToggleImageEdit() {
    showImageEdit = !showImageEdit
    autoChangeOrientation = false
  }
  function handleHideImageEdit() {
    showImageEdit = false
    autoChangeOrientation = true
  }

  function handleResize(e) {
    console.log(e.detail)
    if (isPortrait) {
      previewWidth = e.detail.x - previewContainer.offsetLeft
    } else {
      previewHeight = e.detail.y - previewContainer.offsetTop
    }
  }

  function resizeUp() {
    sectionOpenSeaDragonSize = null
    sectionTranscriptionSize = null
  }

  function resizeDown() {
    var osd = window.getComputedStyle(sectionOpenSeaDragon)
    sectionOpenSeaDragonSize = {
      width: parseFloat(osd.width),
      height: parseFloat(osd.height),
    }
    var trans = window.getComputedStyle(sectionTranscription)
    sectionTranscriptionSize = {
      width: parseFloat(trans.width),
      height: parseFloat(trans.height),
    }
  }

  function resizeMove(e) {
    if (sectionTranscriptionSize && sectionOpenSeaDragonSize) {
      if (isPortrait) {
        sectionOpenSeaDragon.style.width =
          sectionOpenSeaDragonSize.width + e.detail.x + 'px'
        sectionTranscription.style.width =
          sectionTranscriptionSize.width - e.detail.x + 'px'
      } else {
        sectionOpenSeaDragon.style.height =
          sectionOpenSeaDragonSize.height + e.detail.y + 'px'
        sectionTranscription.style.height =
          sectionTranscriptionSize.height - e.detail.y + 'px'
      }
    }
  }

  function handleUpdateNewPassageText(event) {
    newpassagetext = event.detail
  }

  function handleBlockEditorClick() {
    isBlockEditor = !isBlockEditor
    passagetext = passage.transcriptionLines.join('\r\n')
  }

  function handleLineEditorClick() {
    isLineEditor = !isLineEditor

    newpassage = passage
    prevTranscriptionLines = [...passage.transcriptionLines]
    currentTranscriptionLines = [...passage.transcriptionLines]
    lineIndex = 0
    const spans = document.querySelectorAll(
      '.transcription  span[contenteditable]'
    )

    spans.forEach((span) => {
      span.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
          event.returnValue = false
        }
      })
      span.addEventListener('input', (event) => {
        const index = event.target.getAttribute('data-index')
        const newText = event.target.textContent
        currentTranscriptionLines[index] = newText
      })
      span.addEventListener('click', (event) => {
        const selectedLine = event.target.closest('[data-index]')
        if (selectedLine) {
          lineIndex = selectedLine.getAttribute('data-index')
          console.log('Selected line index:', lineIndex)
          transcriptionEditNav = true
        }
      })
    })
  }

  function handleTranscriptionEditorCancel() {
    if (isBlockEditor) {
      isBlockEditor = false
    }
    if (isLineEditor) {
      isLineEditor = false
      currentTranscriptionLines = [...prevTranscriptionLines]
    }
    //location.reload(true)
  }

  function handleTranscriptionEditorSave() {
    if (isBlockEditor) {
      SaveBlockEditorTrancription()
    }
    if (isLineEditor) {
      SaveLineEditorTrancription()
    }
    //location.reload(true)
  }

  async function SaveLineEditorTrancription() {
    // passage.transcriptionlines = currentTranscriptionLines.map(
    //   (obj) => obj.line
    // )
    // passagetext = passage.transcriptionLines.join('\r\n')
    newpassage.transcriptionlines = [...currentTranscriptionLines]
    await fetch('/api/v2/savepassage', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newpassage),
    }).then((res) => {
      window.location.href = `/view/${id_col}/${passage.passageid}`
    })
  }

  async function SaveBlockEditorTrancription() {
    if (newpassagetext == undefined) {
      console.log('newpassagetext empty')
      return
    }
    let data = {
      colid: passage.colid,
      passageid: passage.passageid,
      text: newpassagetext,
    }

    console.log('SaveBlockEditorTrancription data', data)

    document.body.style.cursor = 'wait'
    try {
      const response = await fetch('/api/v2/savepassagetext', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
      // If the fetch operation is successful, reload the page
      if (response.ok) {
        window.location.href = `/view/${id_col}/${passage.passageid}`
      } else {
        console.error('Failed to save data:', response)
      }
    } catch (error) {
      console.error('Error during fetch:', error)
    } finally {
      // Restore the cursor to its default style
      document.body.style.cursor = 'default'
    }
  }

  function setTranscriptionFontSize() {
    //var elem = document.getElementById("transcription")
    //elem.style.fontSize = transcriptionFontSize+"em"
    transcriptionFontSize = tFontSize + 'em'
  }

  function fontDown() {
    if (tFontSize > 0.5) {
      tFontSize -= 0.1
    }
    setTranscriptionFontSize()
  }

  function fontReset() {
    tFontSize = 1
    setTranscriptionFontSize()
  }

  function fontUp() {
    if (tFontSize < 1.5) {
      tFontSize += 0.1
    }
    setTranscriptionFontSize()
  }

</script>

<style>
  #app {
    height: 100%;
    display: flex;
    overflow: hidden;
  }

  :global(.ffull) {
    display: flex;
    flex-direction: column;
    height: 100vh; /* This ensures that the container is exactly the size of the viewport */
    overflow: scroll;
    align-items: stretch;
  }

  :global(.fbody) {
    margin: 0;
    padding: 0;
    height: 100%;
  }

  :global(.fnavbar) {
  }

  :global(.fmain-content) {
    flex-grow: 1; /* Takes up remaining space */
    overflow: scroll;
  }

  :global(.ffooter) {
    position: static;
  }

  .transcriptionBtn {
    border: none;
    background-color: lightblue;
    color: #0077be;
    font-weight: bold;
    padding: 0;
    width: 35px;
    height: 35px;
  }
  .transcriptionBtn:focus {
    outline: none;
  }

  body,
  .vertical-split.portrait {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
  .wait-cursor {
    cursor: wait;
  }
  :global(.toolbar) {
    display: flex;
    flex-direction: row;
    align-items: center;

    height: 40px;
    overflow: scroll;
    margin: 0;
    border-top: 1px solid var(--toolbar-border-color);
    border-bottom: 1px solid var(--toolbar-border-color);
    padding: 0;

    list-style: none;
    background: var(--toolbar-bg-color, white);
    /*border-top: 1px solid black;*/
    /*border-bottom: 1px solid black;*/
  }

  #osd_toolbar {
    height: 34px;
    width: 100%;
    background: var(--toolbar-bg-color, white);
    /*border-bottom: 1px solid black;*/
  }

  #osd_buttons {
    width: 100%;
  }

  :global(.toolbar, .toolbar select) {
    font: 400 14px/100% 'Inter', sans-serif;
    color: var(--toolbar-text-color, black);
  }

  :global(.toolbar.stacked-below) {
    border-top: 0;
  }

  :global(.toolbar li) {
    flex-shrink: 0;
    margin: 0;
    padding: 2px 6px;
  }

  /* join left */
  :global(.toolbar li.jl) {
    margin-left: 2px;
  }
  /* space left */
  :global(.toolbar li.pl) {
    margin-left: 6px;
  }
  /* border left */
  :global(.toolbar li.bl) {
    border-left: 1px solid var(--toolbar-border-color);
  }
  /* fill left */
  :global(.toolbar li.fl) {
    margin-left: auto;
  }

  :global(.toolbar li:last-child) {
    border-right: 0;
  }

  :global(.toolbar li label) {
    margin: 0;
    padding: 0;

    font-weight: 600;
    font-size: inherit;
  }

  :global(.toolbar li code) {
    font: 14px/100% 'IBM Plex Mono', monospace;
    font-weight: bold;
    color: inherit;
  }

  :global(.pane) {
    flex-basis: 10%;
    flex-shrink: 0;
    flex-grow: 1;
  }

  :global(.pane.static) {
    flex-grow: 0;
  }

  :global(.pane.grow) {
    display: flex;
    flex-direction: column;
  }

  :global(.pane.vertical-split.portrait) {
    border-left: 1px solid var(--toolbar-border-color);
    border-top: 0;
    flex-direction: column;
  }
  :global(.vertical-split) {
    display: flex;
    flex-direction: row;
    height: 100%;
    width: 100%;
  }

  :global(.vertical-split .pane) {
    border-left: 1px solid var(--toolbar-border-color);
  }
  :global(.vertical-split .pane:first-child) {
    border-left: 0;
  }

  :global(.preview) {
    display: flex;
    flex-direction: column;
    flex-grow: 0;
    flex-shrink: 0;

    box-sizing: border-box;
    padding: 4px;

    background: var(--pane-bg-color);
    width: var(--width, 100%);
  }

  :global(.openseadragon-container) {
    flex-grow: 1;
  }

  :global(.transcription) {
    box-sizing: border-box;
    padding: 16px;
  }

  :global(.metadata) {
    font: 14px/130% 'Inter', sans-serif;
    background: white;
  }

  :global(.metadata dl dt) {
    padding: 8px 16px 4px;
  }

  .metadata dl dd {
    padding: 0px 16px 8px;
  }

  :global(
      .metadata dl dt:nth-child(4n + 1),
      .metadata dl dt:nth-child(4n + 1) + dd
    ) {
    background: var(--pane-bg-color);
  }

  :global(.close-pane) {
    font-size: 20px;
    font-weight: 600;
  }

  :global(.close-pane:hover) {
    text-decoration: none;
  }
  :global(button) {
    background-color: #0077be;
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 5px;
    cursor: pointer;
  }

  :global(button.is-landscape) {
    /* Styles for landscape layout */
  }

  :global(button.is-portrait) {
    /* Styles for portrait layout */
  }

  .btn {
    min-width: 120px;
  }

  .btn_cancel {
    background: #ababab;
  }
  .btn_cancel:hover {
    background: #dadada;
    color: #ffffff;
  }

  .btn_cancel:active {
    background: #616161;
    color: #ffffff;
  }

  .shadow * {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08);
    pointer-events: none;
    cursor: not-allowed;
    opacity: 0.8;
  }

  .openseadragon .dijitToolbar {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
  }

  .flex {
    display: flex;
  }

  .frow {
    flex-direction: row;
  }

  .fcol {
    flex-direction: column;
  }

  .fgrow {
    flex-grow: 1;
  }
  :global(.scroll) {
    overflow: scroll;
  }

  :global(.pheight) {
    height: 100%;
  }

  .sectionh {
    height: 40vh;
  }

  .sectionw {
    width: 50vw;
    height: 80vh;
  }
  .transh {
    height: 85%;
  }

  .transw {
    width: 100%;
    height: 92%;
  }

  .openseadragon .dijitToolbar {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
  }

</style>

<div class="flex fcol pheight scroll" class:wait-cursor={loading}>
  {#if err}
    <h5>Error: You do not have permission to view this collection</h5>
    <h6>
      You can return to the
      <a href="/">Brucheion home page</a>
      to see all available collections
      <span>{err}</span>
    </h6>
  {:else}
    <nav class="flex" class:shadow={hidePassageToolbar}>
      <ul class="toolbar" style="width: 100vw">
        <li><label for="passage.passageid">Witness:</label></li>
        <li class="pl">
          <div class="select">
            <select
              bind:value={selectedCatalogUrn}
              on:change={handleWitnessSelection}>
              {#if passage && passage.textRefs}
                {#each passage.textRefs as ref}
                  <option value={ref}>{ref}</option>
                {/each}
              {:else}
                <option disabled>No witnesses</option>
              {/if}
            </select>
          </div>
        </li>

        <li class="pl">
          <a href="#top" on:click|preventDefault={handleFirstPassage}>&lt</a>
        </li>
        <li class="pl">
          <a href="#top" on:click|preventDefault={handlePreviousPassage}>←</a>
        </li>
        <!-- <li class="jl"><code>{passageId}</code></li> -->
        <select on:change={handleSelect} bind:value={passageId}>
          {#if passage && passage.schemes}
            {#each schemes as schema (schema)}
              <option value={schema}>{schema}</option>
            {/each}
          {:else}
            <option disabled>No schemes</option>
          {/if}
        </select>
        <li class="pl">
          <a href="#top" on:click|preventDefault={handleNextPassage}>→</a>
        </li>
        <li class="pl">
          <a href="#top" on:click|preventDefault={handleLastPassage}>&gt</a>
        </li>

        <li><label for="passage.passageid">Folio:</label></li>
        <li>
          <div class="select">
            <select
              bind:value={selectedImageRef}
              on:change={handleFolioSelection}>
              {#if passage && passage.imageRefs}
                {#each passage.imageRefs as ref}
                  <option value={ref}>{ref}</option>
                {/each}
              {:else}
                <option disabled>No image references</option>
              {/if}
            </select>
          </div>
        </li>
        <li class="pl">
          <a href="#top" on:click|preventDefault={handleToggleMetadata}>
            {#if showMetadata}Hide{:else}Show{/if}
            Metadata
          </a>
        </li>
        {#if userrights && userrights.canEditMetadata}
          <li class="pl">
            <a href="#top" on:click|preventDefault={handleToggleImageEdit}>
              {#if showImageEdit}Hide{:else}Show{/if}
              Image References
            </a>
          </li>
        {/if}
      </ul>
    </nav>

    <div
      class={`${isPortrait ? 'flex frow' : 'flex fcol'}`}
      style="width: 100%">
      <section
        id="openseadragon_sec"
        class={`flex  fcol ${isPortrait ? ' sectionw' : ' sectionh'}`}>
        <div id="osd_toolbar" class="flex frow">
          <div id="osd_buttons" class="flex frow" />
          <div class="fgrow" />
          <div
            id="lineOpts"
            class="flex frow"
            style="padding-top: 3px; padding-right: 7px; display: none">
            <input
              type="range"
              min="-100"
              max="100"
              value="0"
              id="lineOffsetRange"
              style="width: 200px"
              on:input={updateLinePosition} />
            <div style="padding-right: 5px" />
            <input
              type="color"
              id="myColorPicker"
              on:input={updateLinePosition}
              value={localStorage.getItem('linecolor') === null ? '#FF0000' : localStorage.getItem('linecolor')} />
          </div>
          <a
            href="#top"
            on:click|preventDefault={handleToggleLine}
            style="text-align: right; padding-right:3px; padding-top: 3px; white-space: nowrap">
            {#if showReadingLine}Disable{:else}Enable{/if}
            Reading Line
          </a>
        </div>
        <div
          class="loading"
          style="position:absolute; left: 10px; top: 40px; font-size: 16px;
          z-index: 10;">
          Loading image...
        </div>
        <div
          class="preview openseadragon fgrow"
          bind:this={previewContainer}
          id="preview" />
      </section>

      {#if !isPortrait}
        <ResizeBarLandscape
          on:up={resizeUp}
          on:down={resizeDown}
          on:move={resizeMove} />
      {:else}
        <ResizeBarPortrait
          on:up={resizeUp}
          on:down={resizeDown}
          on:move={resizeMove} />
      {/if}

      <section
        id="transcription_sec"
        class={`flex scroll ${isPortrait ? ' sectionw' : ' sectionh'}`}>
        <div class={`flex vertical-split${isPortrait ? ' portrait' : ''}`}>
          <div class="pane pheight">
            <ul class="toolbar" class:shadow={hideTranscriptionToolbar}>
              <li>
                <label
                  for="passage.passageid"
                  style="font-size: large">Transcription</label>
              </li>
              <li>
                <button
                  on:click={fontDown}
                  type="button"
                  class="transcriptionBtn"
                  style="font-size: large">A&#711</button>
                <button
                  on:click={fontReset}
                  type="button"
                  class="transcriptionBtn"
                  style="font-size: large">A</button>
                <button
                  on:click={fontUp}
                  type="button"
                  class="transcriptionBtn"
                  style="font-size: large">A&#710</button>
              </li>
              <li />
              {#if userrights && userrights.canEditTranscription}
                {#if !(isBlockEditor || isLineEditor)}
                  <li>
                    <button
                      type="button"
                      class=""
                      on:click={handleLineEditorClick}
                      disabled={isBlockEditor && !isLineEditor}>
                      Line Editor
                    </button>
                  </li>

                  <li>
                    <button
                      type="button"
                      class=""
                      on:click={handleBlockEditorClick}
                      disabled={!isBlockEditor && isLineEditor}>
                      Block Editor
                    </button>
                  </li>
                {/if}
                {#if isBlockEditor || isLineEditor}
                  <li>
                    <button
                      type="button"
                      class="btn_save"
                      on:click={handleTranscriptionEditorSave}>
                      Save
                    </button>
                  </li>
                  <li>
                    <button
                      type="button"
                      class="btn_cancel"
                      on:click={handleTranscriptionEditorCancel}>
                      Cancel
                    </button>
                  </li>
                {/if}

                <!-- <li>
                <a
                  href={`/tools/edittranscription/${id_col}/${passage.passageid}`}>
                  Legacy Editor
                </a>
              </li> -->
              {/if}
            </ul>
            <div class="scroll transh ${isPortrait ? ' transw' : ''}">
              {#if isBlockEditor}
                <BlockEditor
                  {passagetext}
                  {newpassagetext}
                  {transcriptionFontSize}
                  on:update={handleUpdateNewPassageText} />
              {:else}
                <LineEditor
                  {lineIndex}
                  {transcriptionEditNav}
                  {currentTranscriptionLines}
                  {isLineEditor}
                  {transcriptionFontSize} />
              {/if}
            </div>
          </div>
          {#if showMetadata}
            <div class="pane">
              <ul class="toolbar" style="display: flex;">
                <li><label for="metadata">Metadata</label></li>
                <span style="flex-grow: 1;" />
                {#if userrights && userrights.canEditMetadata}
                  <button on:click={() => (savedata = true)}>save</button>
                {/if}
                <li class="fl">
                  <a
                    href="#top"
                    class="close-pane"
                    on:click|preventDefault={handleHideMetadata}>
                    ×
                  </a>
                </li>
              </ul>
              <div class="metadata">
                <MetadataEditor
                  {passage}
                  {savedata}
                  canEditMetadata={userrights.canEditMetadata} />
              </div>
            </div>
          {/if}
          {#if showImageEdit}
            <div class="pane">
              <ul class="toolbar" style="display: flex;">
                <li><label for="imageeditor">Image Editor</label></li>
                <span style="flex-grow: 1;" />
                {#if userrights && userrights.canEditMetadata}
                  <button on:click={() => (savedata = true)}>save</button>
                {/if}
                <li class="fl">
                  <a
                    href="#top"
                    class="close-pane"
                    on:click|preventDefault={handleHideImageEdit}>
                    ×
                  </a>
                </li>
              </ul>
              <div class="imageeditor">
                <ImageEditor
                  {passage}
                  colID={id_col}
                  {savedata}
                  on:imageurlchange={handleImageUrlChange} />
              </div>
            </div>
          {/if}
        </div>
      </section>
    </div>
  {/if}
</div>
