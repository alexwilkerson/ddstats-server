<template>
  <v-card v-if="motd !== '' && news !== {}">
    <v-toolbar flat color="th" dark dense>
      <v-toolbar-title>News</v-toolbar-title>
    </v-toolbar>
    <p class="dd-card-list-item">
      <strong>News from {{ moment(news.time_stamp).format("MMMM Do, YYYY") }}:</strong>
      <br />
      {{ news.body }}
    </p>
    <p class="dd-card-list-item">
      <strong>Current MOTD:</strong>
      <br />
      {{ motd }}
    </p>
    <div v-if="data.devil_daggers_list.length > 0" class="dd-card-list-item">
      <strong>New Devil Daggers:</strong>
      <br />
      <ul :style="{listStyle: 'none', margin: 0, padding: 0}">
        <li v-for="player in data.devil_daggers_list" :key="player.player_id">
          <v-icon fill="#c33409" small>$dagger</v-icon>
          <strong>{{ player.player_name }}</strong>
          - {{ player.game_time }}s
        </li>
      </ul>
    </div>
  </v-card>
</template>

<script>
import axios from "axios";

const moment = require("moment");

export default {
  props: ["data"],
  data() {
    return {
      motd: "",
      news: {},
      moment: moment
    };
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
        .get(process.env.VUE_APP_API_URL + `/api/v2/motd`)
        .then(response => {
          this.motd = response.data.motd;
        })
        .catch(error => window.console.log(error));
    },
    getNewsFromAPI() {
      axios
        .get(
          process.env.VUE_APP_API_URL + `/api/v2/news?page_size=1&page_num=1`
        )
        .then(response => {
          this.news = response.data.news[0];
        })
        .catch(error => window.console.log(error));
    }
  },
  mounted() {
    this.getMOTDFromAPI();
    this.getNewsFromAPI();
  }
};
</script>

<style>
</style>