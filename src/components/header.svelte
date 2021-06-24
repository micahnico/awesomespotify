<script>
  import { getContext, onMount } from 'svelte'

  const client = getContext('client')

  let user
  let loading = true

  onMount(async () => {
    const response = await $client.get(`/api/user/get`)
    if (response.ok) {
      user = response.body
      loading = false
    }
  })

  const logIn = async () => {
    const response = await $client.post(`/api/login`)
    if (response.ok) {
      location.href = response.body.url
    }
  }

  const logOut = async () => {
    const response = await $client.post(`/api/logout`)
    if (response.ok) {
      location.reload()
    }
  }
</script>

<div class="w-full flex justify-center pt-10 pb-5 px-10" style="background-color: #111111;">
  <div class='w-full flex justify-between w-full'>
    <p class="text-spotify-green text-2xl md:text-4xl font-bold header-text">Awesome Spotify</p>
    <div>
      {#if !loading}
        {#if user}
          <button class="bg-spotify-green text-white px-3 py-2 rounded-md text-sm font-medium" on:click={logOut}>Log Out</button>
        {:else}
          <button class="bg-spotify-green text-white px-3 py-2 rounded-md text-sm font-medium" on:click={logIn}>Log In</button>
        {/if}
      {/if}
    </div>
  </div>
</div>
