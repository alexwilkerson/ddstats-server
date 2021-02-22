<template>
  <v-card class="stats-table">
    <v-simple-table dense>
      <template v-slot:default>
        <tbody>
          <tr>
            <td class="text-left">
              <v-icon
                class="icon"
                :style="{ marginLeft: '4px', marginRight: '10px' }"
                color="#c33409"
                small
              >mdi-trophy</v-icon>Rank
            </td>
            <td class="text-right">{{ data.rank.toLocaleString() }}</td>
          </tr>
          <tr>
            <td class="text-left">
              <v-icon
                class="icon"
                :style="{ marginLeft: '4px', marginRight: '10px' }"
                color="#c33409"
                small
              >mdi-account</v-icon>Player ID
            </td>
            <td class="text-right">{{ data.player_id }}</td>
          </tr>
          <tr
            class="pointer"
            v-if="data.high_score_game_id"
            @click="$router.push('/games/' + data.high_score_game_id)"
          >
            <td class="text-left">
              <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>Player Best Time
            </td>
            <td
              class="text-right"
              :style="{textDecoration: 'underline'}"
            >{{ data.game_time.toFixed(4) }}s</td>
          </tr>
          <tr v-else>
            <td class="text-left">
              <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>Player Best Time
            </td>
            <td class="text-right">{{ data.game_time.toFixed(4) }}s</td>
          </tr>
          <tr>
            <td class="text-left">
              <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>Average Game Time
            </td>
            <td class="text-right">{{ data.overall_average_game_time.toFixed(4) }}s</td>
          </tr>
          <tr @mouseover="gameTimeHover = true" @mouseleave="gameTimeHover = false">
            <td class="text-left">
              <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>Overall Game Time
            </td>
            <td class="text-right">
              {{
              gameTimeHover
              ? `${data.overall_game_time.toLocaleString()}s`
              : moment.duration(data.overall_game_time, "seconds").humanize()
              }}
            </td>
          </tr>
          <tr>
            <td class="text-left">
              <v-icon
                class="icon"
                :style="{ marginLeft: '4px', marginRight: '10px' }"
                color="#c33409"
                small
              >mdi-grave-stone</v-icon>Overall Deaths
            </td>
            <td class="text-right">{{ data.overall_deaths.toLocaleString() }}</td>
          </tr>
          <tr>
            <td class="text-left">
              <v-icon class="icon" fill="#c33409" small>$gem</v-icon>Overall Gems
            </td>
            <td class="text-right">{{ data.overall_gems.toLocaleString() }}</td>
          </tr>
          <tr>
            <td class="text-left">
              <v-icon class="icon" fill="#c33409" small>$splat</v-icon>Overall Enemies Killed
            </td>
            <td class="text-right">{{ data.overall_enemies_killed.toLocaleString() }}</td>
          </tr>
          <tr @mouseover="accuracyHover = true" @mouseleave="accuracyHover = false">
            <td class="text-left">
              <v-icon class="icon" fill="#c33409" small>$target</v-icon>Overall Accuracy
            </td>
            <td class="text-right" v-if="!accuracyHover">{{ data.accuracy }}%</td>
            <td class="text-right" v-else>
              {{ data.daggers_hit.toLocaleString() }} /
              {{ data.daggers_fired.toLocaleString() }}
            </td>
          </tr>
          <tr v-if="data.last_active" @mouseover="dateHover = true" @mouseleave="dateHover = false">
            <td class="text-left">
              <v-icon
                class="icon"
                :style="{ marginLeft: '4px', marginRight: '10px' }"
                color="#c33409"
                small
              >mdi-desktop-classic</v-icon>Last Active
            </td>
            <td class="text-right">
              {{
              dateHover
              ? moment
              .utc(data.last_active)
              .local()
              .format("lll")
              : moment(data.last_active).fromNow()
              }}
            </td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </v-card>
</template>

<script>
const moment = require("moment");
export default {
  props: ["data"],
  data() {
    return {
      moment: moment,
      gameTimeHover: false,
      accuracyHover: false
    };
  }
};
</script>

<style scoped>
.stats-table {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  border-radius: 2px;
  max-width: 650px;
  margin: 0 auto 20px auto;
  padding: 12px;
}
.icon {
  margin-top: -4px;
  margin-right: 6px;
}
ul {
  list-style: none;
  font-size: 14px;
}
.player-best-header {
  text-align: center;
  margin-bottom: 20px;
}
</style>
