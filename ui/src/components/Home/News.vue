<template>
  <v-card v-if="motd !== ''">
    <v-toolbar flat color="th" dark dense>
      <v-toolbar-title>News</v-toolbar-title>
    </v-toolbar>
    <p class="dd-card-list-item"><strong>MOTD:</strong><br />{{ motd }}</p>
    <p
      v-if="data.devil_daggers_list.length > 0"
      class="dd-card-list-item"
    >
    <strong>Game News:</strong><br />
    Congratulations to {{ newDevilDaggerers() }} on passing 500 seconds and getting their devil daggers!
    <ul :style="{listStyle: 'none', margin: 0, padding: 0}">
      <li v-for="player in data.devil_daggers_list" :key="player.player_id"><strong>{{ player.player_name }}</strong> - {{ player.game_time }}s</li>
    </ul>
    </p>

  </v-card>
</template>

<script>
import axios from "axios";

export default {
  props: ["data"],
  data() {
    return {
      motd: ""
    }
  },
  methods: {
    newDevilDaggerers() {
      let users = [];
      for (let i = 0; i < this.data.devil_daggers_list.length; i++) {
        users.push(this.data.devil_daggers_list[i].player_name);
      }
      return users.join(", ").replace(/,(?!.*,)/gim, ", and");
    },
    getMOTDFromAPI() {
      axios
        .get(
          process.env.VUE_APP_API_URL + `/api/v2/motd`
        )
        .then(response => {
          this.motd = response.data.motd;
        })
        .catch(error => window.console.log(error));
    }
  },
  mounted() {
    this.getMOTDFromAPI()
  }
};
</script>

<style>
</style>