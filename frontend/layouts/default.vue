<template>
  <div class="app-container">
    <div class="sidebar">
      <a href="/"><img src="https://megahook.in/static/img/logo.png"></a>
      <div class="nav-items">
        <a class="nav-item" href="/">Home</a>
        <a class="nav-item" href="https://blog.megahook.in/testing-shopify-with-megahook/">Example</a>
      </div>
    </div>
    <nuxt />
  </div>
</template>

<style lang="stylus">
html
  font-family: 'Source Sans Pro', -apple-system, BlinkMacSystemFont, 'Segoe UI',
    Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  word-spacing: 1px;
  -ms-text-size-adjust: 100%;
  -webkit-text-size-adjust: 100%;
  -moz-osx-font-smoothing: grayscale;
  -webkit-font-smoothing: antialiased;
  box-sizing: border-box;

body
  background: #232840
  color: #9ac6c9
  font-feature-settings: "liga" 0
  line-height: 1.7
  font-family: IBM Plex Mono, monospace
  font-weight: 400
  margin: auto
  padding: 0

*,
*:before,
*:after
  box-sizing: border-box;
  margin: 0;

a
  color: #f649a7

img
  width: 300px
  display: block
  margin: -20px auto

.app-container
  display: flex

.sidebar
  text-align: center
  border-right: 1px dashed #f649a7
  height: 100vh
  position: sticky
  top: 0

.nav-items
  display: flex
  flex-direction: column

  .nav-item
    margin: 0.5rem 0
    font-size: 18px
    text-decoration: none
    font-weight: bold

@media screen and (max-width: 900px)
  body
    padding: 0 1rem

  .app-container
    flex-direction: column

  .sidebar
    height: auto
    position: relative
    border-right: 0

  .nav-items
    display: none
</style>
