<template>
  <v-container class="home-container">
    <v-row v-if="$root.mobile">
      <v-col cols="12">
        <LivePlayers />
      </v-col>
    </v-row>
    <h1 class="text-center">
      {{
      loading
      ? "Loading..."
      : "Stats from " + moment(data.time_stamp).format("MMMM Do, YYYY")
      }}
    </h1>
    <v-row>
      <v-col cols="12" sm="9">
        <Main :data="data" :loading="loading" />
      </v-col>
      <v-col cols="12" sm="3">
        <LivePlayers v-if="!$root.mobile" />
        <TopGames :style="{ marginTop: '24px' }" />
        <RecentGames :style="{ marginTop: '24px' }" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import axios from "axios";
import Main from "../components/Home/Main";
import LivePlayers from "../components/Home/LivePlayers";
import RecentGames from "../components/Home/RecentGames";
import TopGames from "../components/Home/TopGames";

const moment = require("moment");

export default {
  data() {
    return {
      data: {},
      loading: true,
      moment: moment
    };
  },
  components: {
    Main,
    LivePlayers,
    RecentGames,
    TopGames
  },
  methods: {
    getDailyFromAPI() {
      this.loading = true;
      axios
        .get(process.env.VUE_APP_API_URL + `/api/v2/daily`)
        .then(response => {
          this.data = response.data;
          this.loading = false;
        })
        .catch(error => window.console.log(error));
    }
  },
  mounted() {
    this.getDailyFromAPI();
  }
};
</script>

<style>
.home-container {
  padding-top: 40px;
  max-width: 950px;
}
</style>
