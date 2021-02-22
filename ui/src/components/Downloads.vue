<template>
  <v-card :style="{ padding: '12px' }" class="grotesk">
    <v-row align="center" :style="{ margin: '12px' }">
      <v-col class="text-center">
        <v-btn
          color="logo"
          xLarge
          class="text-center"
          :href="'/static/' + data.releases[0].file_name"
          :style="{ marginBottom: '12px' }"
          dark
          >Download v{{ data.releases[0].version }}</v-btn
        >
        <br />
        <a
          class="grotesk"
          :style="{ fontSize: '15px' }"
          href="https://github.com/alexwilkerson/ddstats-go"
          >View code on GitHub</a
        >
      </v-col>
    </v-row>
    <div v-for="release in data.releases" :key="release.version">
      <span :style="{ color: '#c33409' }">
        <strong>v{{ release.version }}</strong></span
      ><span :style="{ fontSize: '15px' }">
        -- {{ moment(release.time_stamp).format("MMM D, YYYY") }}</span
      >
      <ul :style="{ fontSize: '15px' }">
        <li v-for="(note, i) in release.notes" :key="release.version + i">
          {{ note }}
        </li>
      </ul>
      <br />
    </div>
  </v-card>
</template>

<script>
const moment = require("moment");
export default {
  props: ["data"],
  data() {
    return {
      moment: moment
    };
  }
};
</script>

<style></style>
