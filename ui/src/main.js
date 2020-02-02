import Vue from "vue";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import router from "./router";
import VueCookies from "vue-cookies";
import VueNativeSock from "vue-native-websocket";

Vue.use(VueNativeSock, "ws://localhost:5000/ws", {
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
      players: []
    };
  },
  methods: {
    checkPlayerLive(id) {
      for (let i = 0; i < this.players.length; i++) {
        if (this.players[i].player_name == id) return true;
      }
      return false;
    }
  },
  created() {
    addEventListener("resize", () => {
      this.mobile = innerWidth <= 700;
    });
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
      }
    };
  },
  render: h => h(App)
}).$mount("#app");
