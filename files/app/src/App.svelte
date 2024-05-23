<script>
  import { Router, Route } from 'svelte-routing'
  import IngestCex from './routes/IngestCex.svelte'
  import AddCollection from './routes/AddCollection.svelte'
  import IngestImage from './routes/IngestImage.svelte'
  import CollectionOverview from './routes/CollectionOverview.svelte'
  import CollectionList from './routes/Collections.svelte'
  import NotFound from './routes/NotFound.svelte'
  import ShareCollection from './routes/ShareCollection.svelte'
  import CollectionBulkEdit from './routes/CollectionBulkEdit.svelte'

  // export let url = location.pathname.toLowerCase()
  export let url = location.pathname
  // console.log('from dispatcher', url)

</script>

<svelte:head>
  <style>
    :root {
      --toolbar-bg-color: rgb(240, 240, 240);
      --toolbar-text-color: rgb(50, 50, 50);
      --toolbar-border-color: rgb(200, 200, 200);
      --pane-bg-color: rgba(248, 248, 248);
    }

  </style>
</svelte:head>

<Router {url}>
  <Route path="/collections/:id_auth/:id_page" component={CollectionList} />
  <Route path="/view/:id_col/:urn" component={CollectionOverview} />
  <Route path="/tools/bulkedit/:id_col" component={CollectionBulkEdit} />
  <Route path="/tools/sharecollection/:id_col" component={ShareCollection} />
  <Route path="/tools/ingestcex" component={IngestCex} />
  <Route path="/tools/addcollection" component={AddCollection} />
  <Route path="/tools/ingestimage/:id_col" component={IngestImage} />
  <Route path="*" component={NotFound} />
</Router>
