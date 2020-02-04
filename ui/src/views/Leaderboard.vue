<template>
  <div class="wrapper">
    <h1>{{ spawnset }} Leaderboard</h1>
    <div v-if="loaded">
      <v-select
        :options="spawnsets"
        class="style-chooser"
        :clearable="false"
        @input="go"
        placeholder="Select Leaderboard"
      ></v-select>
      <div v-if="!$root.mobile" class="dagger-legend">
        <v-icon fill="#c33409">$dagger</v-icon>
        >= {{data.devil_dagger_time}}s
        <v-icon fill="#ffcd00">$dagger</v-icon>
        >= {{data.gold_dagger_time}}s
        <v-icon fill="#acacac">$dagger</v-icon>
        >= {{data.silver_dagger_time}}s
        <v-icon fill="#ff8300">$dagger</v-icon>
        >= {{data.bronze_dagger_time}}s
        <v-icon fill="#000">$dagger</v-icon>
        &lt; {{data.bronze_dagger_time}}s
      </div>
      <LeaderboardTable :loadingTable="loadingTable" :data="data" />
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
      this.loaded = false;
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
  max-width: 860px;
  margin: auto;
}
h1 {
  text-align: center;
  margin-bottom: 15px;
  color: var(--v-primary-base);
}
.v-select {
  max-width: 400px;
  margin: 0 auto 20px auto;
}
.dagger-legend {
  text-align: center;
  margin: 0 auto;
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 13px;
  margin-bottom: 20px;
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
  font-family: "alte_haas_grotesk", "Helvetica Neue", Helvetica, Arial;
}

.style-chooser .vs__dropdown-option {
  margin-left: -23px;
  color: var(--v-deselected-base);
}

.style-chooser .vs__dropdown-option--highlight {
  color: var(--v-primary-base);
  background-color: var(--v-highlight-base);
}

.vs__actions {
  cursor: pointer;
}

.vs__clear,
.vs__open-indicator {
  fill: var(--v-primary-base);
}
</style>
