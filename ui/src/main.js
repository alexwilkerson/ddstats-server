import Vue from "vue";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import router from "./router";
import VueCookies from "vue-cookies";

Vue.config.productionTip = false;

import "@/assets/global.css";

new Vue({
  VueCookies,
  vuetify,
  router,
  data() {
    return {
      mobile: window.innerWidth <= 700
    };
  },
  created() {
    addEventListener("resize", () => {
      this.mobile = innerWidth <= 700;
    });
  },
  render: h => h(App)
}).$mount("#app");
