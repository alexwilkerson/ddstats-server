<template>
  <div class="wrapper">
    <div v-if="!loading">
      <h1>
        <v-icon class="icon" fill="#c33409">$flourish_left</v-icon>
        {{ data.player_name }}
        <v-icon class="icon" fill="#c33409">$flourish_right</v-icon>
      </h1>
      <PlayerInfo :data="data" />
      <div>
        <h1 class="recorded-games-header">Recorded Games</h1>
        <RecentPlayerGamesTable @onPlayerNameLoad="onPlayerNameLoad" />
      </div>
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
import RecentPlayerGamesTable from "../components/RecentPlayerGamesTable";
import PlayerInfo from "../components/PlayerInfo";
import "vue-select/dist/vue-select.css";

export default {
  data() {
    return {
      data: {},
      loading: true,
      playerName: ""
    };
  },
  components: {
    PlayerInfo,
    RecentPlayerGamesTable
  },
  methods: {
    onPlayerNameLoad(name) {
      this.playerName = name;
    },
    getPlayerFromAPI() {
      this.loading = true;
      axios
        .get(
          process.env.VUE_APP_API_URL +
            `/api/v2/player/update?id=${this.$route.params.id}`
        )
        .then(response => {
          this.data = response.data;
          this.loading = false;
        })
        .catch(error => window.console.log(error));
    }
  },
  mounted() {
    this.getPlayerFromAPI();
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
.recorded-games-header {
  margin-top: 20px;
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
