<template>
  <div class="page">
    <a href="/"><img src="/static/img/logo.png" /></a>
    <hr />
    <div class="description">Easily forward webhook traffic to your local environment</div>
    <hr />
    <div v-if="hasMore" class="more-button">
      <a href="#" v-on:click="showMore">Show More Requests &uarr;</a>
    </div>
    <div class="card" v-for="record in records" v-bind:id="record.id">
      <div class="card__header">
        <a v-bind:href="'#' + record.id">{{record.timestamp}}</a>
      </div>
      <div class="card__content">
        <div>
          <h4>Request</h4>
          <code><pre>{{record.request}}</pre></code>
        </div>
        <div>
          <h4>Response</h4>
          <code><pre>{{record.response}}</pre></code>
        </div>
      </div>
      <div class="card__footer"></div>
    </div>
  </div>
</template>

<script>

import moment from 'moment';

function formatRecords(records) {
  return records.sort((a, b) => (
    a.timestamp < b.timestamp ? 1 : -1
  )).map((r) => {
    const record = {
      id: r.timestamp,
      timestamp: moment(r.timestamp * 1e3).format('HH:mm:ss MMM DD, YYYY'),
      request: JSON.stringify(r.request, null, 2),
      response: JSON.stringify(r.response, null, 2),
    };

    return record;
  });
}

export default {
  data() {
    return {
      hiddenRecords: [],
      records: [],
      hasMore: false,
    };
  },
  created() {
    this.fetchRecords();
  },
  methods: {
    showMore() {
      this.hasMore = false;
      this.records = this.hiddenRecords;
      window.scrollTo(0, 0);
    },
    fetchRecords() {
      // Fetch records via API
      const { name } = this.$route.params;

      fetch(`/api/hooks/${name}/history`).then((res) => res.json()).then((res) => {
        if (!this.records.length) {
          this.records = formatRecords(res.records);
        } else if (this.records.length !== res.records.length) {
          // Hold on until use requests to see more
          this.hasMore = true;
          this.hiddenRecords = formatRecords(res.records);
        }

        setTimeout(this.fetchRecords, 10e3);
      }).catch((err) => {
        console.log('Failed to fetch history:', err);
      });
    },
  }
};

</script>

<style scoped lang="scss">

.more-button {
  text-align: center;
  padding: 10px 0 20px;
  position: sticky;
  top: 0;
}

.more-button > a {
  text-shadow: 0 0 5px #232840;
}

.card {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #f649a7;
  margin-bottom: 1rem;
}

.card__header {
  padding: 20px;
  border-bottom: 1px solid #f649a7;
}

.card__content {
  padding: 20px;
}

code {
  padding: 8px;
}

pre {
  margin: 0;
  padding: 0;
  white-space: break-spaces;
  word-wrap: break-word;
}

</style>
