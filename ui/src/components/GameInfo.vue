<template>
  <div class="wrapper">
    <h1>Game Stats {{ $route.params.id }}</h1>
    <v-simple-table class="stats-table" dense>
      <template v-slot:default>
        <tbody>
          <tr>
            <td class="text-left">Player</td>
            <td class="text-right">
              {{ gameInfo.player_name }} ({{ gameInfo.player_id }})
            </td>
          </tr>
          <tr>
            <td class="text-left">Death Type</td>
            <td class="text-right">{{ gameInfo.death_type }}</td>
          </tr>
          <tr>
            <td class="text-left">Spawnset</td>
            <td class="text-right">{{ gameInfo.spawnset }}</td>
          </tr>
          <tr>
            <td class="text-left">Recorded</td>
            <td class="text-right">
              {{ moment(gameInfo.time_stamp).fromNow() }}
            </td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
    <v-row no-gutters justify="center">
      <v-col cols="12" sm="2" align="center">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <div v-on="on" class="icon-info">
              <v-icon class="icon" fill="#c33409" small>$stopwatch</v-icon>
              <span>{{ gameInfo.game_time }}</span>
            </div>
          </template>
          <span>Game Time</span>
        </v-tooltip>
      </v-col>
      <v-col cols="12" sm="2" align="center">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <div v-on="on" class="icon-info">
              <v-icon class="icon" fill="#c33409" small>$gem</v-icon>
              <span>{{ gameInfo.gems }}</span>
            </div>
          </template>
          <span>Gems</span>
        </v-tooltip>
      </v-col>
      <v-col cols="12" sm="2" align="center">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <div v-on="on" class="icon-info">
              <v-icon class="icon" fill="#c33409">$dagger</v-icon>
              <span>{{ gameInfo.homing_daggers }}</span>
            </div>
          </template>
          <span>Homing Daggers</span>
        </v-tooltip>
      </v-col>
      <v-col cols="12" sm="2" align="center">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <div v-on="on" class="icon-info">
              <v-icon class="icon" fill="#c33409" small>$target</v-icon>
              <span>{{ gameInfo.accuracy }}%</span>
            </div>
          </template>
          <span>Accuracy</span>
        </v-tooltip>
      </v-col>
      <v-col cols="12" sm="2" align="center">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <div v-on="on" class="icon-info">
              <v-icon class="icon" fill="#c33409" small>$skull</v-icon>
              <span>{{ gameInfo.enemies_alive }}</span>
            </div>
          </template>
          <span>Enemies Alive</span>
        </v-tooltip>
      </v-col>
      <v-col cols="12" sm="2" align="center">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <div v-on="on" class="icon-info">
              <v-icon class="icon" fill="#c33409" small>$splat</v-icon>
              <span>{{ gameInfo.enemies_killed }}</span>
            </div>
          </template>
          <span>Enemies Killed</span>
        </v-tooltip>
      </v-col>
    </v-row>
  </div>
</template>

<script>
const moment = require("moment");
export default {
  data() {
    return {
      moment: moment
    };
  },
  props: {
    gameInfo: {
      type: Object,
      required: true
    }
  }
};
</script>

<style scoped>
.stats-table {
  max-width: 600px;
  margin: 0 auto 20px auto;
}
.stats-table tr {
  background: #f6f2ee;
}
.wrapper {
  max-width: 780px;
  margin: 0 auto 40px auto;
  background: #f6f2ee;
}
.icon {
  margin-top: -4px;
  margin-right: 6px;
}
.game-info {
  padding: 0;
}
.game-info li {
  font-size: 15px;
  display: inline;
  cursor: default;
}
.game-info div {
  display: inline;
}
.icon-info span {
  font-size: 15px;
  cursor: default;
}
h1 {
  text-align: center;
  padding-bottom: 20px;
}
tbody {
  background: #ebe7e4;
}
ul {
  list-style: none;
  font-size: 14px;
}
</style>
