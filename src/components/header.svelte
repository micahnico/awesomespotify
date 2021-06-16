<script>
  import { getContext, onMount } from 'svelte'

  const client: any = getContext('client')

  let user: any
  let loading: boolean = true

  onMount(async () => {
    const response = await $client.get(`/api/user/get`)
    if (response.ok) {
      user = response.body
      loading = false
    }
  })

  const logIn = async () => {
    const response = await $client.get(`/api/login`)
    if (response.ok) {
      location.reload() // TODO: temporary, should be fixed after resolving auth opening in other tab
    }
  }

  const logOut = async () => {
    const response = await $client.get(`/api/logout`)
    if (response.ok) {
      location.reload()
    }
  }
</script>

<div class="w-full flex justify-center pt-10 pb-5 px-6 bg-gray-100">
  <div class='w-full flex justify-between lg:w-3/4 xl:w-3/5'>
    <p class="text-spotify-green text-3xl font-bold">Awesome Spotify</p>
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
