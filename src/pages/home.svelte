<script>
  import Header from '../components/header.svelte'
  import {link} from 'svelte-spa-router'
  import { getContext, onMount } from 'svelte'

  const client: any = getContext('client')

  let artist: string
  let song: string
  let lyrics: string

  onMount(async () => {
    const response = await $client.get(`/api/lyrics/find`)
    if (response.ok) {
      console.log(response.body)
      artist = response.body.Artist
      song = response.body.Song
      lyrics = response.body.Lyrics
    }
  })
</script>

<Header/>

<div class="w-full flex justify-center p-5 bg-gray-100">
  <div class='px-16 py-10 w-full lg:w-1/2 bg-white rounded-md shadow-lg'>
    <p class='text-5xl font-bold mb-1'>{song}</p>
    <p class='text-2xl text-gray-600 mb-5'>{artist}</p>
    {@html lyrics}
    {#if lyrics}
      <hr class="my-5">
      <p>These lyrics were taken from <a href="https://genius.com" class="text-blue-500 hover:text-blue-700">genius.com</a></p>
    {/if}
  </div>
</div>
