<template>
  <div>
    <h1>Megahook</h1>
    <p>Megahook is a utility for forwarding webhooks to your local environment.</p>

    <h4>Install</h4>
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