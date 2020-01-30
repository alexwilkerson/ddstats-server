<template>
  <v-app
    id="main"
    dark
    :style="{ background: $vuetify.theme.themes[theme].background }"
  >
    <v-app-bar
      app
      color="header"
      height="60px"
      elevation="1"
      :style="{ textAlign: 'center', zIndex: 500 }"
    >
      <v-app-bar-nav-icon
        v-if="$vuetify.breakpoint.xsOnly"
        @click.stop="drawer = !drawer"
      ></v-app-bar-nav-icon>
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
      <NavBar v-if="!$vuetify.breakpoint.xsOnly" />
      <v-navigation-drawer
        v-model="drawer"
        fixed
        temporary
        app
        :style="{ top: '60px' }"
      >
        <v-list nav dense>
          <v-list-item-group v-model="group">
            <v-list-item to="/">
              <v-list-item-title>HOME</v-list-item-title>
            </v-list-item>
            <v-list-item to="/leaderboard">
              <v-list-item-title>LEADERBOARDS</v-list-item-title>
            </v-list-item>
            <v-list-item to="/games">
              <v-list-item-title>GAMES</v-list-item-title>
            </v-list-item>
            <v-list-item to="/players">
              <v-list-item-title>PLAYERS</v-list-item-title>
            </v-list-item>
            <v-list-item to="/info">
              <v-list-item-title>INFO</v-list-item-title>
            </v-list-item>
            <v-list-item to="/download">
              <v-list-item-title>DOWNLOAD</v-list-item-title>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-navigation-drawer>
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
  data: () => ({
    drawer: false,
    group: null
  }),
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
  },
  watch: {
    group() {
      this.drawer = false;
    }
  }
};
</script>

<style scoped>
.dd-logo {
  position: absolute;
  top: 50%;
  left: 50%;
  margin-top: -18px;
  margin-left: -18px;
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
  -webkit-font-smoothing: antialiased !important;
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
