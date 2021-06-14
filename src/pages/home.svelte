<script>
  import Header from '../components/header.svelte'
  import {link} from 'svelte-spa-router'
  import { getContext, onMount } from 'svelte'

  const client: any = getContext('client')

  let artist: string
  let song: string
  let lyrics: string

  let user: any
  let loading: boolean = true

  onMount(async () => {
    const userResponse = await $client.get(`/api/user/get`)
    if (userResponse.ok) {
      user = userResponse.body
    }

    if (user) {
      const lyricResponse = await $client.get(`/api/lyrics/find`)
      if (lyricResponse.ok) {
        artist = lyricResponse.body.Artist
        song = lyricResponse.body.Song
        lyrics = lyricResponse.body.Lyrics
      }
    }

    loading = false
  })

  const logIn = async () => {
    const response = await $client.get(`/api/login`)
    if (response.ok) {
      location.reload() // TODO: temporary, should be fixed after resolving auth opening in other tab
    }
  }
</script>

<Header/>

<div class="w-full flex justify-center p-5 bg-gray-100">
  <div class='px-16 py-10 w-full lg:w-1/2 bg-white rounded-md shadow-lg'>
    {#if loading}
      <p class='text-xl font-bold mb-1'>Finding Lyrics...</p>
    {:else}
      {#if user}
        <p class='text-5xl font-bold mb-1'>{song}</p>
        <p class='text-2xl text-gray-600 mb-5'>{artist}</p>
        {@html lyrics}
        {#if lyrics}
          <hr class="my-5">
          <p>These lyrics were taken from <a href="https://genius.com" class="text-blue-500 hover:text-blue-700">genius.com</a></p>
        {/if}
      {:else}
        <p class='text-3xl font-bold mb-1'>Lyrics</p>
        <p><span class="cursor-pointer link" on:click={logIn}>Log in with Spotify</span> to access</p>
      {/if}
    {/if}
  </div>
</div>
