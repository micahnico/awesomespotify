// [snowpack] add styles to the page (skip if no document exists)
if (typeof document !== 'undefined') {
  const code = "@-webkit-keyframes svelte-1kjgs54-rotation{from{transform:rotate(0deg)}to{transform:rotate(359deg)}}@keyframes svelte-1kjgs54-rotation{from{transform:rotate(0deg)}to{transform:rotate(359deg)}}.rotate.svelte-1kjgs54{-webkit-animation:svelte-1kjgs54-rotation 5s infinite linear;animation:svelte-1kjgs54-rotation 5s infinite linear}.custom-img-size.svelte-1kjgs54{height:200px;width:200px}.currently-playing-info.svelte-1kjgs54{display:flex;flex-direction:column}@media(min-width: 640px){.custom-img-size.svelte-1kjgs54{height:250px;width:250px}}@media(min-width: 1280px){.currently-playing-info.svelte-1kjgs54{flex-direction:row}}@media(min-width: 1536px){.custom-img-size.svelte-1kjgs54{height:300px;width:300px}}";

  const styleEl = document.createElement("style");
  const codeEl = document.createTextNode(code);
  styleEl.type = 'text/css';
  styleEl.appendChild(codeEl);
  document.head.appendChild(styleEl);
}