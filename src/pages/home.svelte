<script>
  import Header from '../components/header.svelte'
  import {link} from 'svelte-spa-router'
  import { getContext, onMount } from 'svelte'

  const client: any = getContext('client')

  onMount(async () => {
    const response = await $client.get(`/api/test`)
    if (response.ok) {
      console.log(response.body)
    }
  })

</script>

<Header/>

<p>Home</p>
<a use:link href="/testpage">Test Page</a>
