<script lang='typescript'>
	import {onMount, setContext} from 'svelte'
  import { writable } from 'svelte/store'
  import Router from 'svelte-spa-router'
  import { wrap } from 'svelte-spa-router/wrap.js'
  import Client from './client.js'

  const client = writable(new Client())
  setContext('client', client)

  const routes = {
    '/': wrap({asyncComponent: () => import('./pages/home.svelte')}),
    '*': wrap({asyncComponent: () => import('./pages/not_found.svelte')}),
  }

</script>

<link rel="preconnect" href="https://fonts.gstatic.com">
<link href="https://fonts.googleapis.com/css2?family=Montserrat:wght@400;700&family=Pacifico&display=swap" rel="stylesheet">

<style global lang="postcss">
  @tailwind base;
  @tailwind components;
  @tailwind utilities;

  * {
    font-family: 'Montserrat', sans-serif;
  }

  .header-text {
    font-family: 'Pacifico', cursive;
  }

  .link {
    color: #1DB954;
  }

  .link:hover {
    color: #199c47;
  }

  .text-spotify-green {
    color: #1DB954;
  }

  .bg-spotify-green {
    background: #1DB954;
  }

  .bg-spotify-dark-green {
    background: #199c47;
  }
</style>

<div class="pseudo-body">
  <Router {routes} />
</div>
