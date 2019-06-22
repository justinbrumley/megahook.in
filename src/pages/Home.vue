<template>
  <div>
    <h1>Megahook</h1>
    <hr />
    <p>Megahook is a utility for forwarding webhooks to your local environment.</p>
    <hr />

    <h4>Downloads</h4>
    <a href="/static/latest/megahook.deb">Linux</a>
    <hr />
    <h4>Install Manually</h4>
    <p>
      <a href="https://golang.org/doc/install" target="_blank">Make sure that you have go installed first.</a>
      Then run the following:
    </p>
    <code>
      go get github.com/justinbrumley/megahook
      go install github.com/justinbrumley/megahook
    </code>
    <p>
      Or if you prefer using git:
    </p>
    <code>
      git clone github.com/justinbrumley/megahook
      cd megahook
      go install
    </code>
    <p>And finally, connect to the server and start receiving webhook traffic:</p>
    <code>
      megahook http://localhost:8080/my/favorite/webhook my-little-webhook
    </code>
    <p>
      You should be given a URL you can start using for your webhooks. If the name you chose
      is already taken, you will be given a randomly generated one.
    </p>
  </div>
</template>

<script>
  export default {};
</script>

<style lang="stylus">
  body
    background: #232840
    color: #9ac6c9
    font-feature-settings: "liga" 0
    line-height: 1.7
    font-family: IBM Plex Mono, monospace
    font-weight: 400
    max-width: 768px
    margin: auto

  h1
    text-align: center

  hr
    border-style: dashed
    border-color: #f649a7

  a
    color: #f649a7

  code
    padding: 20px 10px
    background: #1C2033
    border-radius: 4px
    margin-bottom: 10px
    display: inline-block
    font-size: 14px

  h1, h2, h3, h4, h5, h6
    text-shadow: 0 2px 10px #A2E0D3
</style>