<template>
  <v-tabs
    :backgroundColor="$vuetify.theme.themes[theme].background"
    centered
    height="32px"
    sliderColor="#c33409"
    color="#c33409"
    app
  >
    <v-tab to="/">HOME</v-tab>
    <v-tab to="/leaderboard">LEADERBOARDS</v-tab>
    <v-tab to="/games">GAMES</v-tab>
    <v-tab to="/players">PLAYERS</v-tab>
    <v-tab to="/info">INFO</v-tab>
    <v-tab to="/download">DOWNLOAD</v-tab>
  </v-tabs>
</template>

<script>
export default {
  props: {
    backgroundColor: String
  },
  computed: {
    theme() {
      return this.$vuetify.theme.dark ? "dark" : "light";
    }
  }
};
</script>

<style></style>
