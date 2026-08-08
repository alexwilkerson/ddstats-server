<template>
  <v-card>
    <v-toolbar flat color="th" dark dense>
      <v-toolbar-title>Top Games</v-toolbar-title>
    </v-toolbar>
    <v-simple-table v-if="!loading">
      <tbody>
        <tr v-for="game in data.games" :key="game.id" @click="selectGame(game)" class="pointer">
          <td class="dd-card-list-item-bold truncate-name">
            <v-icon
              v-if="$root.checkPlayerLive(game.player_id)"
              class="icon online-green"
              small
            >mdi-access-point</v-icon>
            <span class="truncate-name-text" :title="game.player_name">{{ game.player_name }}</span>
          </td>
          <td class="dd-card-list-item text-right shrink-column">{{ game.game_time.toFixed(4) }}s</td>
        </tr>
      </tbody>
    </v-simple-table>
  </v-card>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      data: {},
      loading: true
    };
  },
  methods: {
    selectGame: function(game) {
      this.$router.push("/games/" + game.id);
    },
    getGamesFromAPI() {
      this.loading = true;
      axios
        .get(
          process.env.VUE_APP_API_URL +
            `/api/v2/leaderboard?spawnset=v3&page_size=5&page_num=1`
        )
        .then(response => {
          this.data = response.data;
          this.loading = false;
        })
        .catch(error => window.console.log(error));
    }
  },
  mounted() {
    this.getGamesFromAPI();
  }
};
</script>

<style>
.since-text {
  color: #aeaeae;
  font-family: "alte_haas_grotesk", "Helvetica Neue", Helvetica, Arial;
  font-style: oblique;
  font-size: 12px !important;
}

.truncate-name {
  max-width: 0;
}

.shrink-column {
  width: 1%;
  white-space: nowrap;
}

.truncate-name-text {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}
</style>
