<template>
  <v-tabs
    :backgroundColor="backgroundColor"
    centered
    height="32px"
    sliderColor="#c33409"
    color="#c33409"
    app
  >
    <v-tab>HOME</v-tab>
    <v-tab>GAMES</v-tab>
    <v-tab>PLAYERS</v-tab>
    <v-tab>INFO</v-tab>
    <v-tab>DOWNLOAD</v-tab>
  </v-tabs>
</template>

<script>
export default {
  props: {
    backgroundColor: String
  }
};
</script>

<style></style>
