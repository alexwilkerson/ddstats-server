<template>
  <div>
    <v-container v-if="!loading" :style="{padding: '0'}">
      <v-row>
        <v-col cols="12" :style="{padding: '0 12px'}">
          <News :data="data" />
        </v-col>
      </v-row>
      <v-row v-if="!$root.mobile">
        <v-col cols="12">
          <PlayerCard :data="data" />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <DailyServerStats :data="data" />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <GlobalServerStats :data="data" />
        </v-col>
      </v-row>
    </v-container>
    <v-progress-circular
      class="progress"
      v-else
      :size="100"
      :width="6"
      color="#c33409"
      indeterminate
    ></v-progress-circular>
  </div>
</template>

<script>
import News from "./News";
import GlobalServerStats from "./GlobalServerStats";
import PlayerCard from "./PlayerCard";
import DailyServerStats from "./DailyServerStats";

const moment = require("moment");

export default {
  data() {
    return {
      moment: moment
    };
  },
  props: ["data", "loading"],
  components: {
    News,
    GlobalServerStats,
    DailyServerStats,
    PlayerCard
  }
};
</script>

<style>
</style>