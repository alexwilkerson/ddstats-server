<template>
  <div class="wrapper">
    <h1>Download</h1>
    <div v-if="!loading">
      <Downloads :data="data" />
    </div>
    <v-progress-circular
      class="progress"
      v-else
      :size="100"
      :width="6"
      color="#c33409"
      indeterminate
    ></v-progress-circular>
  </div>
</template>

<script>
import axios from "axios";
import Downloads from "../components/Downloads";

export default {
  data: () => ({
    loading: true,
    data: {}
  }),
  components: {
    Downloads
  },
  methods: {
    getReleasesFromAPI() {
      this.loading = true;
      axios
        .get(
          process.env.VUE_APP_API_URL +
            `/api/v2/releases?page_size=5&page_num=1`
        )
        .then(response => {
          this.data = response.data;
          this.loading = false;
        })
        .catch(error => window.console.log(error));
    }
  },
  mounted() {
    this.getReleasesFromAPI();
  }
};
</script>

<style scoped>
.wrapper {
  padding-left: 10px;
  padding-right: 10px;
  padding-top: 40px;
  padding-bottom: 40px;
  max-width: 860px;
  margin: auto;
}
h1 {
  text-align: center;
  margin-bottom: 15px;
  color: var(--v-primary-base);
}
</style>