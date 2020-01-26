<template>
  <v-app
    id="main"
    dark
    :style="{ background: $vuetify.theme.themes[theme].background }"
  >
    <v-app-bar
      app
      color="header"
      dark
      height="60px"
      elevation="1"
      :style="{ textAlign: 'center', zIndex: 500 }"
    >
      <v-spacer />
      <Logo
        width="40"
        @click="switchThemes()"
        :fill="$vuetify.theme.themes[theme].logo"
        class="shrink mr-2 dd-logo"
        transition="scale-transition"
      />
      <v-spacer />
    </v-app-bar>

    <v-content>
      <NavBar />
      <router-view></router-view>
    </v-content>

    <v-footer app color="footer" absolute>
      <v-spacer />
      <div class="footer-text">&#x263A; 2020 - VHS</div>
    </v-footer>
  </v-app>
</template>

<script>
import NavBar from "./components/NavBar";
import Logo from "./assets/logo.svg?inline";
export default {
  name: "DDSTATS",
  components: {
    NavBar,
    Logo
  },
  methods: {
    switchThemes: function() {
      this.$vuetify.theme.dark = !this.$vuetify.theme.dark;
      if (this.$vuetify.theme.dark) {
        window.$cookies.set("dark", true, "10y");
      } else {
        window.$cookies.remove("dark");
      }
    }
  },
  computed: {
    theme() {
      return this.$vuetify.theme.dark ? "dark" : "light";
    }
  },
  created() {
    if (window.$cookies.get("dark")) {
      this.$vuetify.theme.dark = true;
    }
  }
};
</script>

<style scoped>
.dd-logo {
  cursor: pointer;
}
</style>

<style>
@font-face {
  font-family: "alte_haas_grotesk_regular";
  src: url("./assets/alte_haas_grotesk_regular.ttf") format("truetype");
  font-weight: normal;
  font-style: normal;
}
@font-face {
  font-family: "alte_haas_grotesk_bold";
  src: url("./assets/alte_haas_grotesk_bold.ttf") format("truetype");
  font-weight: normal;
  font-style: normal;
}
body {
  /* background-color: #fffefc; */
  color: #34302e;
  font-family: "alte_haas_grotesk_regular", "Helvetica Neue", Helvetica, Arial;
}
h1 {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 22px;
  color: #34302e;
}
h2 {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 30px;
}
h3 {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 20px;
}
.footer-text {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 13px;
  color: var(--v-logo-base);
}
</style>
