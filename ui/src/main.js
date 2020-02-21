import Vue from "vue";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import router from "./router";
import VueCookies from "vue-cookies";
import VueNativeSock from "vue-native-websocket";
import EventBus from "./event-bus";

Vue.use(VueNativeSock, process.env.VUE_APP_WS_URL, {
  reconnection: true,
  reconnectionAttempts: 5,
  reconnectionDelay: 3000
});

Vue.config.productionTip = false;

import "@/assets/global.css";

new Vue({
  VueCookies,
  vuetify,
  router,
  data() {
    return {
      mobile: window.innerWidth <= 700,
      players: [],
      state: {
        game_time: 0,
        gems: 0,
        homing_daggers: 0,
        enemies_alive: 0,
        enemies_killed: 0,
        level_two_time: 0,
        level_three_time: 0,
        level_four_time: 0,
        death_type: -1,
        is_replay: false
      },
      status: "Dead",
      watchers: 0,
      ws_queue: []
    };
  },
  methods: {
    checkPlayerLive(id) {
      for (let i = 0; i < this.players.length; i++) {
        if (this.players[i].player_id == id) return true;
      }
      return false;
    },
    daggerColor(game_time) {
      if (game_time >= 500) {
        return "#c33409";
      } else if (game_time >= 250) {
        return "#ffcd00";
      } else if (game_time >= 120) {
        return "#acacac";
      } else if (game_time >= 60) {
        return "#ff8300";
      } else {
        return "#000";
      }
    }
  },
  created() {
    addEventListener("resize", () => {
      this.mobile = innerWidth <= 700;
    });
    this.$options.sockets.onopen = function() {
      while (this.ws_queue.length > 0) {
        this.$socket.send(this.ws_queue.pop());
      }
    };
    this.$options.sockets.onmessage = function(msg) {
      let data = JSON.parse(msg.data);
      let body = JSON.parse(data.body);
      switch (data.func) {
        case "player_list":
          this.$root.players = body.players;
          break;
        case "player_logged_in":
          this.$root.players = [...this.$root.players, body];
          break;
        case "player_logged_off":
          this.$root.players = this.$root.players.filter(
            player => player.id == body.player_id
          );
          break;
        case "submit":
          if (body !== undefined) {
            this.state = body;
          }
          if (body.status === undefined) {
            this.status = "";
          } else {
            if (body.death_type != -1) {
              this.status = "Dead";
            } else {
              if (body.status !== "Dead" && body.status !== "Alive") {
                this.state.level_two_time = 0;
                this.state.level_three_time = 0;
                this.state.level_four_time = 0;
              }
              this.$root.status = body.status;
            }
          }
          break;
        case "status":
          this.status = body;
          if (body != "Dead") {
            this.state = {
              game_time: 0,
              gems: 0,
              homing_daggers: 0,
              enemies_alive: 0,
              enemies_killed: 0,
              level_two_time: 0,
              level_three_time: 0,
              level_four_time: 0,
              death_type: -1,
              is_replay: false
            };
          }
          break;
        case "user_count":
          this.watchers = body.count;
          break;
        case "game_submitted":
          EventBus.$emit("game_submitted", body);
          break;
      }
    };
  },
  render: h => h(App)
}).$mount("#app");
