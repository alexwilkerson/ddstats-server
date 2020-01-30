<template>
  <div class="wrapper" v-if="loaded">
    <h1>{{ spawnset }} Leaderboard</h1>
    <v-select
      :options="spawnsets"
      class="style-chooser"
      :clearable="false"
      :value="spawnset"
      @input="go"
      placeholder="Select Spawnset"
    ></v-select>
    <LeaderboardTable :loadingTable="loadingTable" :data="data" />
  </div>
</template>

<script>
import LeaderboardTable from "../components/LeaderboardTable";
import vSelect from "vue-select";
import axios from "axios";
import "vue-select/dist/vue-select.css";

export default {
  data: () => ({
    loadingTable: false,
    loaded: false,
    spawnset: null,
    data: null,
    spawnsets: null
  }),
  methods: {
    go: function(spawnset) {
      this.$router.push("/leaderboard/" + spawnset);
      this.spawnset = spawnset;
      this.loadingTable = true;
      this.data = {
        items: [],
        page_size: 0,
        page_num: 1,
        total_game_count: 0
      };
      axios
        .get(
          process.env.VUE_APP_API_URL +
            "/api/v2/leaderboard?spawnset=" +
            spawnset
        )
        .then(response => {
          let spawnsets = [];
          for (let i = 0; i < response.data.spawnsets.length; i++) {
            if (
              (this.$route.params.name === undefined &&
                response.data.spawnsets[i] === "v3") ||
              response.data.spawnsets[i] === this.$route.params.name
            ) {
              continue;
            }
            spawnsets.push(response.data.spawnsets[i]);
          }
          this.spawnsets = spawnsets;
          this.data = response.data;
          this.loaded = true;
          this.loadingTable = false;
        })
        .catch(error => window.console.log(error));
    }
  },
  components: {
    "v-select": vSelect,
    LeaderboardTable
  },
  mounted() {
    let spawnset = this.$route.params.name;
    if (spawnset === undefined) {
      spawnset = "v3";
    }
    this.spawnset = spawnset;
    axios
      .get(
        process.env.VUE_APP_API_URL + "/api/v2/leaderboard?spawnset=" + spawnset
      )
      .then(response => {
        let spawnsets = [];
        for (let i = 0; i < response.data.spawnsets.length; i++) {
          if (
            (this.$route.params.name === undefined &&
              response.data.spawnsets[i] === "v3") ||
            response.data.spawnsets[i] === this.$route.params.name
          ) {
            continue;
          }
          spawnsets.push(response.data.spawnsets[i]);
        }
        this.spawnsets = spawnsets;
        this.data = response.data;
        this.loaded = true;
      })
      .catch(error => window.console.log(error));
  }
};
</script>

<style scoped>
.wrapper {
  padding-left: 10px;
  padding-right: 10px;
  padding-top: 40px;
  padding-bottom: 40px;
  max-width: 800px;
  margin: auto;
}
h1 {
  text-align: center;
  margin-bottom: 20px;
  color: var(--v-primary-base);
}
.v-select {
  max-width: 400px;
  margin: 0 auto 20px auto;
}
</style>

<style>
.style-chooser .vs__search::placeholder,
.style-chooser .vs__search,
.style-chooser .vs__dropdown-toggle,
.style-chooser .vs__dropdown-menu,
.style-chooser .vs__selected {
  background: var(--v-footer-base);
  color: var(--v-primary-base);
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
}

.style-chooser .vs__dropdown-option {
  margin-left: -23px;
  color: var(--v-deselected-base);
  font-family: "alte_haas_grotesk", "Helvetica Neue", Helvetica, Arial;
}

.style-chooser .vs__dropdown-option--highlight {
  color: var(--v-primary-base);
  background-color: var(--v-highlight-base);
}

.vs__clear,
.vs__open-indicator {
  fill: var(--v-primary-base);
}
</style>
