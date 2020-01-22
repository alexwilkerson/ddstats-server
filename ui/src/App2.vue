<template>
  <div id="app">
    <app-header />
    <main />
    <app-footer />
  </div>
</template>

<script>
import Header from "./components/Header.vue";
import Footer from "./components/Footer.vue";
export default {
  components: {
    "app-header": Header,
    "app-footer": Footer
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
  background-color: #fffefc;
  color: #34302e;
  font-family: "alte_haas_grotesk_regular", "Helvetica Neue", Helvetica, Arial;
  padding: 0;
  margin: 0;
}

h1 {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 45px;
}

h2 {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 30px;
}

h3 {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 20px;
}

#app {
  height: 100vh;
  display: flex;

  flex-direction: column;
}

main {
  flex: 1;
}
</style>
