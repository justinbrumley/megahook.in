<template>
  <div>
    <img src="/static/img/logo.png" />
    <hr />
    <div class="description">Easily forward webhook traffic to your local environment</div>
    <hr />
    <div class="content">
      <div>
        <a href="#docker" id="docker">
          <h2># Docker</h2>
        </a>
        <p>
          The docker repo should make running Megahook easier, especially if you don't want to install go.
        </p>
        <a href="https://hub.docker.com/r/justinbrumley/megahook" target="_blank">Docker Repo</a>
        <p>Example Usage:</p>
        <code>
          docker run -d \
          <br />
          &nbsp;&nbsp;-e WEBHOOK_URL=http://localhost:3000/test -e WEBHOOK_NAME=my_hook \
          <br />
          &nbsp;&nbsp;--network host \
          <br />
          &nbsp;&nbsp;justinbrumley/megahook:latest
        </code>
        <p>
          All web traffic to https://megahook.in/m/my_hook should be redirected to http://localhost:3000/test (assuming the name `my_hook` is not taken)
        </p>
        <p>
          Using `--network host` is optional. This just allows the docker container to reach your machine's localhost instead of being confined to the docker container.
        </p>
        <p>
          Once running, you can confirm that the URLs are correct by checking the logs:
        </p>
        <code>
          docker logs --tail 30 CONTAINER_NAME
        </code>
        <hr />
      </div>
      <div>
        <a href="#manual" id="manual">
          <h2># Install Manually</h2>
        </a>
        <p>
          <a href="https://golang.org/doc/install" target="_blank">Make sure that you have go installed first.</a>
          Then run the following:
        </p>
        <code>
          go get github.com/justinbrumley/megahook
          <br />
          go install github.com/justinbrumley/megahook
        </code>
        <p>
          Or if you prefer using git:
        </p>
        <code>
          git clone github.com/justinbrumley/megahook
          <br />
          cd megahook
          <br />
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
    </div>
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
    margin: auto
    padding: 0 20px

  .description
    font-weight: bold
    text-align: center

  .content
    display: flex

    hr
      display: none

    > div:first-child
      border-right: 1px dashed #f649a7

    > div
      width: 50%
      padding: 0 15px
      box-sizing: border-box

  img
    width: 300px
    display: block
    margin: -20px auto

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
    width: 100%;
    box-sizing: border-box;

  @media (max-width: 1350px)
    .content
      flex-direction: column

      hr
        display: block

      > div:first-child
        border: none

      > div
        padding: 0
        width: 100%
</style>
