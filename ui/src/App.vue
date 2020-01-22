<template>
  <v-app
    id="main"
    :style="{ background: $vuetify.theme.themes[theme].background }"
  >
    <v-app-bar
      app
      color="header"
      dark
      height="60px"
      elevation="1"
      :style="{ textAlign: 'center' }"
    >
      <v-spacer />
      <!-- <div class="d-flex align-center"> -->
      <v-img
        alt="Vuetify Logo"
        class="shrink mr-2"
        contain
        dense
        src="./assets/ddstats_logo_v2_header.png"
        transition="scale-transition"
        width="40"
      />
      <!-- <h1>DDSTATS</h1> -->
      <v-spacer />

      <!-- <h1>DDSTATS</h1> -->
      <!-- </div> -->

      <!-- <v-spacer /> -->
    </v-app-bar>

    <v-content>
      <NavBar :backgroundColor="$vuetify.theme.themes[theme].background" />
      <HelloWorld />
    </v-content>

    <v-footer app color="footer">
      <v-spacer />
      <div class="footer-text">&#x263A; 2020 - VHS</div>
    </v-footer>
  </v-app>
</template>

<script>
import HelloWorld from "./components/HelloWorld";
import NavBar from "./components/NavBar";

export default {
  name: "App",

  components: {
    HelloWorld,
    NavBar
  },

  data: () => ({
    //
  }),

  computed: {
    theme() {
      return this.$vuetify.theme.dark ? "dark" : "light";
    }
  }
};
</script>

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
}
</style>
