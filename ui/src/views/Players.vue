<template>
  <div class="wrapper">
    <div v-if="!loading">
      <h1>
        <v-icon class="icon" fill="#c33409">$flourish_left</v-icon>
        {{ data.player_name }}
        <v-icon class="icon" fill="#c33409">$flourish_right</v-icon>
      </h1>
      <PlayerInfo :data="data" />
      <LivePlayer v-if="$root.checkPlayerLive($route.params.id)" />
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
import LivePlayer from "../components/LivePlayer";
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
    LivePlayer,
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
  beforeDestroy() {
    this.$socket.send('{"func": "leave_room" }');
    this.$root.state = {};
  },
  mounted() {
    // this.$root.state = {
    //   game_time: 0,
    //   gems: 0,
    //   homing_daggers: 0,
    //   enemies_alive: 0,
    //   enemies_killed: 0,
    //   level_two_time: 0,
    //   level_three_time: 0,
    //   level_four_time: 0,
    //   death_type: -1,
    //   is_replay: false
    // };
    // this.$root.watchers = 0;
    this.getPlayerFromAPI();
    this.$socket.send(
      '{"func": "join_room", "body": "' + this.$route.params.id + '"}'
    );
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
