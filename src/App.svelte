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
    '/testpage': wrap({asyncComponent: () => import('./pages/testpage.svelte')}),
    '*': wrap({asyncComponent: () => import('./pages/not_found.svelte')}),
  }

</script>

<style global lang="postcss">
  @tailwind base;
  @tailwind components;
  @tailwind utilities;

  @import url('https://fonts.googleapis.com/css2?family=Montserrat&display=swap');
  * {
    font-family: 'Montserrat', sans-serif;
  }

  .link {
    color: #1DB954;
  }

  .link:hover {
    color: #199c47;
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
